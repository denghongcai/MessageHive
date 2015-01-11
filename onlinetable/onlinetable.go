package onlinetable

import (
	"errors"
	"github.com/denghongcai/generalmessagegate/message"
	"github.com/op/go-logging"
	"sync"
	"time"
)

var log = logging.MustGetLogger("main")

type Entity struct {
	Uid       string
	Pipe      chan *message.Container
	LoginTime time.Time
}

type Container struct {
	sync.RWMutex
	storage map[string]*Entity
}

var instance *Container

var initctx sync.Once

func NewContainer() *Container {
	initctx.Do(func() {
		instance = new(Container)
		instance.storage = make(map[string]*Entity)
	})
	return instance
}

func (ct Container) GetEntity(uid string) (*Entity, error) {
	ct.RLock()
	if entity, ok := ct.storage[uid]; ok {
		ct.RUnlock()
		return entity, nil
	}
	ct.RUnlock()
	return new(Entity), errors.New("Entity not found")
}

func (ct *Container) AddEntity(uid string, pipe chan *message.Container) error {
	ct.Lock()
	delete(ct.storage, uid)
	entity := &Entity{Uid: uid, Pipe: pipe, LoginTime: time.Now().UTC()}
	ct.storage[uid] = entity
	ct.Unlock()
	log.Debug("Entity uid: %s added", uid)
	return nil
}

func (ct *Container) DelEntity(uid string) error {
	ct.Lock()
	delete(ct.storage, uid)
	ct.Unlock()
	return nil
}
