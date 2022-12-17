package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"sync-bot/services"
	"sync-bot/types"
)

func main() {

	c, err := types.ParserConfig(os.Getenv(" SYNC_BOT_CONFIG"))
	if err != nil {
		log.Panic(err)
	}

	if c.Run.Debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	tg, err := services.NewTG(c)
	if err != nil {
		log.Panic(err)
	}
	tg.Run()

	<-make(chan interface{})
}
