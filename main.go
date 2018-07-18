package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nlopes/slack"
)

type Message struct {
	ID      string
	Channel string
	Content string
	Timeout time.Time
	Attempt int
}

type MessageManager struct {
	sync.Mutex
	chatID   int64
	bot      *tgbotapi.BotAPI
	wait     time.Duration
	messages map[string]Message
}

func (m *MessageManager) Add(e *slack.DesktopNotification) {
	m.Lock()
	defer m.Unlock()
	m.messages[e.Msg] = Message{
		ID:      e.Msg,
		Timeout: time.Now().Add(m.wait),
		Channel: e.Channel,
		Content: fmt.Sprintf("%s: %s - %s",
			e.Title, e.Subtitle, e.Content,
		),
	}
}

func (m *MessageManager) Delete(e *slack.IMMarkedEvent) {
	var toDelete []Message
	m.Lock()
	for _, message := range m.messages {
		if message.Channel != e.Channel {
			continue
		}
		toDelete = append(toDelete, message)
	}
	for _, s := range toDelete {
		fmt.Println("deleted", s.Content)
		delete(m.messages, s.ID)
	}
	m.Unlock()
}

func (m *MessageManager) Collect(t time.Time) {
	var toSend []Message
	m.Lock()
	for _, message := range m.messages {
		if message.Timeout.After(t) {
			continue
		}
		toSend = append(toSend, message)
	}
	for _, s := range toSend {
		delete(m.messages, s.ID)
	}
	m.Unlock()
	for _, s := range toSend {
		if s.Attempt > 3 {
			fmt.Println("attempts")
			continue
		}
		fmt.Println("sending:", s.Content)
		_, err := m.bot.Send(tgbotapi.NewMessage(m.chatID, s.Content))
		if err != nil {
			log.Println("failed to send:", err)
			m.Lock()
			s.Attempt++
			m.messages[s.ID] = s
			m.Unlock()
		}
	}
}

func main() {
	chatID, err := strconv.ParseInt(os.Getenv("TELEGRAM_TARGET"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}
	api := slack.New(os.Getenv("SLACK_TOKEN"))
	rtm := api.NewRTM()
	go rtm.ManageConnection()
	defer func() {
		log.Println("closing")
		rtm.Disconnect()
	}()
	m := MessageManager{
		messages: make(map[string]Message),
		wait:     time.Second * 10,
		bot:      bot,
		chatID:   chatID,
	}
	go func() {
		ticker := time.NewTicker(time.Second)
		for t := range ticker.C {
			m.Collect(t)
		}
	}()
	for e := range rtm.IncomingEvents {
		switch d := e.Data.(type) {
		case *slack.DesktopNotification:
			fmt.Println("got desktop notify")
			m.Add(d)
		case *slack.IMMarkedEvent:
			fmt.Println("got mark event")
			m.Delete(d)
		default:
			// fmt.Printf("type %T\n", e.Data)
		}
	}
}
