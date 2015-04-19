// 路由模块
package router

import (
	"encoding/json"
	"io"
	"strings"
	"time"

	event_user "github.com/denghongcai/MessageHive/event/user"
	"github.com/denghongcai/MessageHive/message"
	"github.com/denghongcai/MessageHive/onlinetable"
	"github.com/denghongcai/MessageHive/redis"
	"github.com/denghongcai/MessageHive/router/transient"
	"github.com/golang/protobuf/proto"
	"github.com/op/go-logging"
)

// 消息类型定义
const (
	MESSAGE_TYPE_IDENTITY uint = iota
	MESSAGE_TYPE_AUTHENTICATE
	MESSAGE_TYPE_HEARTBEAT
	MESSAGE_TYPE_RECEIPT
	MESSAGE_TYPE_TRANSIENT
	MESSAGE_TYPE_GROUP
	MESSAGE_TYPE_EVENT
	MESSAGE_TYPE_MAX
	MESSAEG_TYPE_ERROR = 31
)

// 群组消息类型定义
const (
	MESSAGE_GROUP_JOIN   = "join"
	MESSAGE_GROUP_INVITE = "invite"
	MESSAGE_GROUP_SEND   = "send"
	MESSAGE_GROUP_LEAVE  = "leave"
)

// 群组消息内容结构
type GroupBody struct {
	Action  string      `json:"action"`
	BodyRaw interface{} `json:"data"`

	List []string
	Data string
}

// 群组消息内容JSON解析
func GroupBodyDecode(r io.Reader) (x *GroupBody, err error) {
	x = new(GroupBody)
	if err = json.NewDecoder(r).Decode(x); err != nil {
		return
	}
	switch t := x.BodyRaw.(type) {
	case string:
		x.Data = t
	case []interface{}:
		bodyraw := x.BodyRaw.([]interface{})
		list := make([]string, 0)
		for i := range bodyraw {
			data, err := bodyraw[i].(string)
			if err == false {
				break
			}
			list = append(list, data)
		}
		x.List = list
	}
	return
}

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
	// Redis init
	redisPool := redis.NewPool(":6379")
	// Event handler init
	eventUserChan := make(chan *event_user.Event, 1000)
	eventUserConfig := event_user.NewConfig(eventUserChan, redisPool, config.mainchan, config.onlinetable)
	go func(config event_user.Config) {
		event_user.Start(eventUserConfig)
	}(eventUserConfig)
	// Transient handler init
	transientChan := make(chan *message.Container, 1000)
	transientConfig := transient.NewConfig(redisPool, transientChan)
	go func(config transient.Config) {
		transient.Handler(config)
	}(transientConfig)

	for {
		select {
		case msg := <-config.mainchan:
			sendflag := true
			sid := msg.GetSID()
			rid := msg.GetRID()
			mtype := msg.GetTYPE()
			sentity, err := config.onlinetable.GetEntity(sid)
			if err != nil {
				if hasBit(mtype, MESSAGE_TYPE_AUTHENTICATE) {
					e := &event_user.Event{
						Uid:  sid,
						Type: event_user.USER_OFFLINE,
					}
					eventUserChan <- e
				}
				log.Info(err.Error())
				break
			}
			response := new(message.Container)
			response.MID = proto.String(msg.GetMID())
			response.SID = proto.String("")
			response.RID = proto.String(sid)
			response.TYPE = proto.Uint32(0)
			response.STIME = proto.Int64(time.Now().Unix())
			response.BODY = proto.String("")
			for i := 0; i < int(MESSAGE_TYPE_MAX); i++ {
				// 消息分类处理
				if hasBit(mtype, uint(i)) {
					switch uint(i) {
					case MESSAGE_TYPE_AUTHENTICATE:
						e := &event_user.Event{
							Uid:  sid,
							Type: event_user.USER_ONLINE,
						}
						eventUserChan <- e
						sendflag = false
						break
					case MESSAGE_TYPE_GROUP:
						body := msg.GetBODY()
						groupbody, err := GroupBodyDecode(strings.NewReader(body))
						if err != nil {
							log.Info(err.Error())
						}
						switch groupbody.Action {
						case MESSAGE_GROUP_JOIN:
							_, err := config.onlinetable.GetEntity(rid)
							if err != nil {
								err = config.onlinetable.AddGroupEntity(rid, groupbody.List)
								if err != nil {
									log.Error(err.Error())
								}
							}
							sendflag = false
							// TODO
							break
						case MESSAGE_GROUP_SEND:
							// PASS
							break
						case MESSAGE_GROUP_INVITE:
							// TODO
							break
						case MESSAGE_GROUP_LEAVE:
							// TODO
							break
						}
						break
					}
				}
			}
			// 发送回应消息
			go func() {
				select {
				case sentity.Pipe <- response:
					log.Info("Response delivered to %s", sid)
				case <-time.After(time.Second):
					log.Error("Response failed to deliverd to %s", sid)
				}
			}()

			// Send to rid
			if sendflag {
				rentity, err := config.onlinetable.GetEntity(rid)
				if err != nil {
					log.Debug(err.Error())
					if hasBit(mtype, MESSAGE_TYPE_TRANSIENT) {
						// MESSAGE_TYPE_TRANSIENT
						// 向Transient队列压入消息
						transientChan <- msg
					}
					break
				}
				switch rentity.Type {
				case onlinetable.ENTITY_TYPE_GROUP:
					for _, v := range rentity.List {
						if v != sid {
							newmsg := *msg
							newmsg.RID = proto.String(v)
							config.mainchan <- &newmsg
						}
					}
				case onlinetable.ENTITY_TYPE_USER:
					go func() {
						select {
						case rentity.Pipe <- msg: // TODO: 这里可能造成死锁
							log.Info("Message delivered from %s to %s", sid, rid)
						case <-time.After(time.Second):
							config.mainchan <- msg
						}
					}()
				}
			}
		}
	}
	return nil
}

// Assume BigEndian

func hasBit(n uint32, pos uint) bool {
	val := n & (1 << pos)
	return (val > 0)
}

func setBit(n uint32, pos uint) uint32 {
	n |= (1 << pos)
	return n
}
