package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"sync-bot/share"
	"sync-bot/utils"
)

type DockerBuild struct {
	Name  string
	Build string
}

func (dn *DockerBuild) Kind() string {
	return "docker-build"
}

func (dn *DockerBuild) Content() string {
	return dn.Build
}

func (dn *DockerBuild) SetReply(reply string) {
	dn.Name = utils.GrabCommand(reply)
}

func (dn *DockerBuild) SetSource(chatId int64) {
	for _, task := range share.DockTask()[chatId] {
		if task.Name != "" && task.Build == "" {
			task.Build = dn.Build
			dn.Name = task.Name
			break
		}
	}
}

func (dn *DockerBuild) Answer() (*string, bool) {
	result := fmt.Sprintf("[%s] has start", dn.Name)
	return &result, true
}

func (dn *DockerBuild) Run() error {

	log.Infof("name: %s  build: %s", dn.Name, dn.Build)
	return nil
}

func (dn *DockerBuild) Inspect(update tgbotapi.Update) {

}
