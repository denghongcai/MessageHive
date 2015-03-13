package transient

import (
	"github.com/denghongcai/MessageHive/message"
	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("main")

type Config struct {
	pool          *redis.Pool
	transientchan chan *message.Container
}

func NewConfig(pool *redis.Pool, transientchan chan *message.Container) Config {
	return Config{
		pool:          pool,
		transientchan: transientchan,
	}
}

func Handler(config Config) {
	for {
		msg := <-config.transientchan
		sid := msg.GetSID()
		conn := config.pool.Get()
		data, _ := proto.Marshal(msg)
		_, err := conn.Do("LPUSH", sid, data)
		if err != nil {
			log.Error(err.Error())
		}
		conn.Close()
	}
}
