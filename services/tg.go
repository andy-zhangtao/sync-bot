package services

import (
	"context"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync-bot/services/commands"
	"sync-bot/share"
	"sync-bot/types"
)

type TG struct {
	bot      *tgbot.BotAPI
	interval int
	task     map[string]commands.Command
	conf     types.Config
}

func NewTG(c types.Config) (*TG, error) {
	bot, err := tgbot.NewBotAPI(c.Run.Token)
	if err != nil {
		return nil, err
	}

	if c.Run.Debug {
		bot.Debug = true
	}

	log.Debugf("%s Authorized Success", bot.Self.UserName)

	share.NewDockerTask()
	share.NewGithubHelper(c)
	return &TG{bot: bot, interval: c.Run.Interval, task: make(map[string]commands.Command), conf: c}, nil
}

func (t TG) Run() {
	go t.queryMessage()
}

func (t TG) queryMessage() error {
	u := tgbot.NewUpdate(0)
	u.Timeout = 30

	for update := range t.bot.GetUpdatesChan(u) {
		if update.Message != nil { // If we got a message
			if update.Message.ReplyToMessage != nil {
				log.Debugf("[%s] %s reply[%s]", update.Message.From.UserName, update.Message.Text, update.Message.ReplyToMessage.Text)
			} else {
				log.Debugf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			}

			cmd, err := t.parserCommand(update.Message.Text)
			if err != nil {
				log.Errorf("parserCommand error. [%s]", update.Message.Text)
				t.sendMessage(update.Message.Chat.ID, update.Message.MessageID, err.Error())
				continue
			}

			ctx := t.wrapContext(t.conf)
			cmd.SetContext(ctx)
			if update.Message.ReplyToMessage != nil {
				ctx = context.WithValue(ctx, "reply", update.Message.ReplyToMessage.Text)
				cmd.SetContext(ctx)
			}

			command := t.convertCommand(cmd)
			command.SetContext(context.WithValue(cmd.Context(), "source", update.Message.Chat.ID))
			command.SetSource(update.Message.Chat.ID)
			answer, _ := command.Answer()
			if answer != nil {
				t.sendMessage(update.Message.Chat.ID, update.Message.MessageID, *answer)
			}

			result, err := command.Run()
			if err != nil {
				t.sendMessage(update.Message.Chat.ID, update.Message.MessageID, err.Error())
				continue
			}

			t.sendMessage(update.Message.Chat.ID, update.Message.MessageID, result)
			//msg := tgbot.NewMessage(update.Message.Chat.ID, fmt.Sprintf("reply from bot, [%s]", update.Message.Text))
			//msg.ReplyToMessageID = update.Message.MessageID
			//
			//t.bot.Send(msg)
		}
	}

	return nil
}

func (t TG) parserCommand(msg string) (cmd commands.Command, err error) {
	msg = strings.TrimSpace(msg)
	cmds := strings.Split(msg, " ")

	var cmdDetail []string

	for _, c := range cmds {
		if _c := strings.TrimSpace(c); _c != "" {
			cmdDetail = append(cmdDetail, _c)
		}
	}

	return &commands.CommonCmd{
		Cmd:   cmdDetail[0],
		Value: strings.TrimSpace(strings.Trim(msg, cmdDetail[0])),
	}, nil

}

func (t TG) convertCommand(cmd commands.Command) commands.Command {
	switch cmd.Kind() {
	case types.DockerName:
		_c := &commands.DockerName{Name: cmd.Content()}
		_c.SetContext(cmd.Context())
		return _c
	case types.DockerBuild:
		_c := &commands.DockerBuild{Build: cmd.Content()}
		_c.SetContext(cmd.Context())
		return _c
	default:
		return cmd
	}
}

func (t TG) wrapContext(c types.Config) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "auth", c.GHelper.Auth)
	ctx = context.WithValue(ctx, "email", c.GHelper.Email)
	return ctx
}
