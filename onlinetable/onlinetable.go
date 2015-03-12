package onlinetable

import (
	"errors"
	"sync"
	"time"

	"github.com/denghongcai/MessageHive/message"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("main")

const (
	ENTITY_TYPE_USER = iota
	ENTITY_TYPE_GROUP
)

type Entity struct {
	Uid       string
	Type      int
	Pipe      chan *message.Container
	List      []string
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
	entity := &Entity{Uid: uid, Type: ENTITY_TYPE_USER, Pipe: pipe, LoginTime: time.Now().UTC()}
	ct.storage[uid] = entity
	ct.Unlock()
	log.Debug("Entity uid: %s added", uid)
	return nil
}

func (ct *Container) AddGroupEntity(uid string, uidlist []string) error {
	//TODO: Ensure every user in group is existed
	ct.Lock()
	entity := &Entity{Uid: uid, Type: ENTITY_TYPE_GROUP, List: uidlist, LoginTime: time.Now().UTC()}
	ct.storage[uid] = entity
	ct.Unlock()
	log.Debug("Group entity uid: %s added", uid)
	return nil
}

func (ct *Container) GetEntities() error {
	ct.Lock()
	ct.Unlock()
	return nil
}

func (ct *Container) DelEntity(uid string) error {
	ct.Lock()
	delete(ct.storage, uid)
	ct.Unlock()
	return nil
}
