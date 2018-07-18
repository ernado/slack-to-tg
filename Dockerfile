FROM golang:1.10.3

COPY . /go/src/github.com/ernado/slack-to-tg

RUN go install github.com/ernado/slack-to-tg

CMD ["/go/bin/slack-to-tg"]