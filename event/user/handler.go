package user

import (
	"github.com/denghongcai/MessageHive/message"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("main")

const (
	USER_ONLINE = iota
)

type Event struct {
	Uid  string
	Type int
}

type Config struct {
	eventchan chan *Event
	mainchan  chan *message.Container
}

func NewConfig(eventchan chan *Event, mainchan chan *message.Container) Config {
	return Config{
		eventchan: eventchan,
		mainchan:  mainchan,
	}
}

func Start(config Config) {
	for {
		e := <-config.eventchan
		switch e.Type {
		case USER_ONLINE:
			log.Debug("UID: %s online", e.Uid)
		}
	}
}
