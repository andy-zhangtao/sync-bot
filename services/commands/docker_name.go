package commands

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sync-bot/share"
	"sync-bot/types"
)

type DockerName struct {
	Name string
	ctx  context.Context
}

func (dn *DockerName) Kind() string {
	return "docker-name"
}

func (dn *DockerName) Content() string {
	return dn.Name
}

func (dn *DockerName) SetContext(ctx context.Context) {
	dn.ctx = ctx
}

func (dn *DockerName) Context() context.Context {
	return dn.ctx
}

func (dn *DockerName) SetReply(string) {}

func (dn *DockerName) SetSource(chatId int64) {
	if task, exist := share.DockTask()[chatId]; exist {
		task = append(task, types.DockerTask{
			Name:  dn.Name,
			Build: "",
		})
	} else {
		share.DockTask()[chatId] = []types.DockerTask{
			types.DockerTask{
				Name:  dn.Name,
				Build: "",
				Stage: types.Ready,
			},
		}
	}
}

func (dn *DockerName) Answer() (*string, bool) {
	result := "Tell me your build script"
	return &result, false
}

func (dn *DockerName) Run() (string, error) {
	return "", nil
}

func (dn *DockerName) Inspect(update tgbotapi.Update) {

}
