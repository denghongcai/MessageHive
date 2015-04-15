package user

import (
	"github.com/denghongcai/MessageHive/message"
	"github.com/denghongcai/MessageHive/onlinetable"
	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
	"time"

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
	eventchan   chan *Event
	pool        *redis.Pool
	mainchan    chan *message.Container
	onlinetable *onlinetable.Container
}

func NewConfig(eventchan chan *Event, pool *redis.Pool, mainchan chan *message.Container, onlinetable *onlinetable.Container) Config {
	return Config{
		eventchan:   eventchan,
		pool:        pool,
		mainchan:    mainchan,
		onlinetable: onlinetable,
	}
}

func Start(config Config) {
	for {
		e := <-config.eventchan
		switch e.Type {
		case USER_ONLINE:
			// Transient handle
			conn := config.pool.Get()
			for {
				data, err := redis.Bytes(conn.Do("RPOP", e.Uid))
				if err != nil {
					break
				}
				msg := new(message.Container)
				err = proto.Unmarshal(data, msg)
				if err != nil {
					log.Error(err.Error())
				}
				config.mainchan <- msg
			}
			conn.Close()
			if e.Uid != "00000001" {
				msg := &message.Container{}
				msg.SID = proto.String(e.Uid)
				msg.RID = proto.String("00000001") // 推送系统UID
				msg.TYPE = proto.Uint32(64)
				msg.STIME = proto.Int64(time.Now().UnixNano())
				msg.BODY = proto.String(`{"type": "online", "data": null}`)
				config.mainchan <- msg
			}
			log.Debug("UID: %s online", e.Uid)
		}
	}
}
