# slack-to-tg
Send slack desktop notifications to telegram (if not read in 10s) with emoji support.

Inspired by [slack-to-telegram-bot](https://github.com/trestoa/slack-to-telegram-bot).
Uses [fork](https://github.com/ernado/slack) of [nlopes/slack](https://github.com/nlopes/slack) package.

## Usage

For configuration, set the following environment variables:
```bash
$ export SLACK_TOKEN=''     # Slack bot token
$ export TELEGRAM_TOKEN=''  # Telegram bot token
$ export TELEGRAM_TARGET='' # Target chat
```
For the target chat id, see that [stackoverflow question]( http://stackoverflow.com/questions/32423837/telegram-bot-how-to-get-a-group-chat-id-ruby-gem-telegram-bot).

You can build and use the docker image (or just use `ernado/slackbot`):
```bash
docker build -t <docker-image-url:docker-image-tag> .
docker push <docker-image-url:docker-image-tag>
docker run -d --name slack-to-telegram-bot --restart=always -e TELEGRAM_TOKEN='$TELEGRAM_TOKEN' -e TELEGRAM_TARGET='$TELEGRAM_TARGET' -e SLACK_TOKEN='$SLACK_TOKEN' <docker-image-url:docker-image-tag>
```
