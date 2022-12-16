package services

import (
	"fmt"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"sync-bot/types"
	"time"
)

type TG struct {
	bot      *tgbot.BotAPI
	interval int
}

func NewTG(c types.Config) (*TG, error) {
	bot, err := tgbot.NewBotAPI(c.Run.Token)
	if err != nil {
		return nil, err
	}

	if c.Run.Debug {
		bot.Debug = true
	}

	log.Printf("%s Authorized Success", bot.Self.UserName)

	return &TG{bot: bot, interval: c.Run.Interval}, nil
}

func (t TG) Run() {
	ticker := time.NewTicker(time.Duration(t.interval) * time.Second)
	for {
		select {
		case <-ticker.C:
			go t.queryMessage()
		}
	}
}

func (t TG) queryMessage() error {
	u := tgbot.NewUpdate(0)
	u.Timeout = 30

	for update := range t.bot.GetUpdatesChan(u) {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbot.NewMessage(update.Message.Chat.ID, fmt.Sprintf("reply from bot, [%s]", update.Message.Text))
			msg.ReplyToMessageID = update.Message.MessageID

			t.bot.Send(msg)
		}
	}

	return nil
}
