// 在线表模块
package onlinetable

import (
	"errors"
	"sync"
	"time"

	"github.com/denghongcai/MessageHive/message"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("main")

// 实体类型定义
const (
	ENTITY_TYPE_USER = iota
	ENTITY_TYPE_GROUP
)

// 实体结构
type Entity struct {
	Uid       string
	Type      int
	Pipe      chan *message.Container
	List      []string
	LoginTime time.Time
}

// 在线表结构
type Container struct {
	sync.RWMutex                    // 同步锁
	storage      map[string]*Entity // 哈希表
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

// 通过UID获取实体
func (ct Container) GetEntity(uid string) (*Entity, error) {
	ct.RLock()
	if entity, ok := ct.storage[uid]; ok {
		ct.RUnlock()
		return entity, nil
	}
	ct.RUnlock()
	return new(Entity), errors.New("Entity not found")
}

// 向在线表中添加实体
func (ct *Container) AddEntity(uid string, pipe chan *message.Container) error {
	ct.Lock()
	delete(ct.storage, uid)
	entity := &Entity{Uid: uid, Type: ENTITY_TYPE_USER, Pipe: pipe, LoginTime: time.Now().UTC()}
	ct.storage[uid] = entity
	ct.Unlock()
	log.Debug("Entity uid: %s added", uid)
	return nil
}

// 向在线表中添加群组实体
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

// 通过UID删除实体
func (ct *Container) DelEntity(uid string) error {
	ct.Lock()
	if _, ok := ct.storage[uid]; ok {
		delete(ct.storage, uid)
		ct.Unlock()
		return nil
	} else {
		ct.Unlock()
		return errors.New("Entity delete failed")
	}
}
