package commands

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Command interface {
	Kind() string
	Content() string
	//SetReply(string)
	SetSource(int64)
	SetContext(ctx context.Context)
	Context() context.Context
	Answer() (*string, bool)
	Run() (string, error)
	Inspect(update tgbotapi.Update)
}

type CommonCmd struct {
	Cmd    string
	Value  string
	Reply  string
	Source int64
	ctx    context.Context
}

func (dn *CommonCmd) Kind() string {
	return dn.Cmd
}

func (dn *CommonCmd) Content() string {
	return dn.Value
}

func (dn *CommonCmd) SetContext(ctx context.Context) {
	dn.ctx = ctx
}

func (dn *CommonCmd) Context() context.Context {
	return dn.ctx
}

func (dn *CommonCmd) SetReply(reply string) {
	dn.Reply = reply
}

func (dn *CommonCmd) SetSource(source int64) {
	dn.Source = source
}

func (dn *CommonCmd) Answer() (*string, bool) {
	result := ""
	return &result, false
}

func (dn *CommonCmd) Run() (string, error) {
	return "", nil
}

func (dn *CommonCmd) Inspect(update tgbotapi.Update) {

}
