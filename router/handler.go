package router

import (
	"fmt"
	"github.com/denghongcai/generalmessagegate/message"
	"github.com/denghongcai/generalmessagegate/onlinetable"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("main")

type Config struct {
	mainchan    chan *message.Container
	onlinetable *onlinetable.Container
}

// Returns config of router
func NewConfig(mainchan chan *message.Container, onlinetable *onlinetable.Container) Config {
	return Config{
		mainchan:    mainchan,
		onlinetable: onlinetable,
	}
}

func Handler(config Config) error {
	for {
		select {
		case msg := <-config.mainchan:
			sid := msg.GetSID()
			rid := msg.GetRID()
			_, err := config.onlinetable.GetEntity(sid)
			if err != nil {
				log.Info(fmt.Sprintf("%s", err))
				break
			}
			rentity, err := config.onlinetable.GetEntity(rid)
			if err != nil {
				log.Info(fmt.Sprintf("%s", err))
				break
			}
			rentity.Pipe <- msg // TODO: this will cause dead lock
			log.Info("Message delivered from %s to %s", sid, rid)
		}
	}
	return nil
}
