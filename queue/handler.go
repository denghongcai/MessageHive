package queue

import (
	"github.com/denghongcai/MessageHive/message"
	"github.com/golang/protobuf/proto"
	"github.com/kr/beanstalk"
	"time"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("main")

type Config struct {
	server string
	in     chan *message.Container
	out    chan *message.Container
}

func NewConfig(server string, in chan *message.Container, out chan *message.Container) Config {
	return Config{
		server: server,
		in:     in,
		out:    out,
	}
}

func Start(config Config) {
	conn, err := beanstalk.Dial("tcp", config.server)
	go func(config Config) {
		id, data, err := conn.Reserve(time.Minute)
		if err != nil {
			log.Error(err.Error())
		} else {
			msg := new(message.Container)
			proto.Unmarshal(data, msg)
			select {
			case config.out <- msg:
				conn.Delete(id)
			case time.After(5 * time.Second):
				log.Error("%d message fail to proceed", id)
			}

		}
	}(config)
	for {
		msg := <-config.in
		data, _ := proto.Marshal(msg)
		conn.Put(data, 1, 0, 30*time.Second)
	}
}
