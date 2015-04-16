// 客户端会话模块
package client

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/denghongcai/MessageHive/auth"
	"github.com/denghongcai/MessageHive/message"
	"github.com/denghongcai/MessageHive/onlinetable"
	"github.com/denghongcai/MessageHive/protocol"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("main")

// 状态机状态表
const (
	AUTH = iota
	CONNECTED
)

// TCP连接Keep-alive超时时间
const timeoutMinutes time.Duration = 1

type Instance struct {
	state          int
	lastActiveTime int64
	conn           net.Conn
	Uid            string
	outchan        chan *message.Container
	inchan         chan *message.Container
	onlinetable    *onlinetable.Container
}

// 新客户端会话创建
func NewInstance(conn net.Conn, outchan chan *message.Container, onlinetable *onlinetable.Container) *Instance {
	instance := new(Instance)
	instance.state = 0
	instance.conn = conn
	instance.outchan = outchan
	instance.onlinetable = onlinetable
	instance.inchan = make(chan *message.Container)
	return instance
}

func (ins *Instance) Handler() {
	go ins.MainReadHandler()
	go ins.MainWriteHandler()
}

// 客户端会话有限状态机
func (ins *Instance) stateMachine(pkt *protocol.Packet) error {
	switch ins.state {
	case AUTH:
		// 调用认证模块进行Token验证
		if err := auth.Authenticate(pkt.Msg); err != nil {
			return err
		}
		// 将用户加入在线表
		if err := ins.onlinetable.AddEntity(pkt.Msg.GetSID(), ins.inchan); err != nil {
			return errors.New(fmt.Sprintf("Entity add failed, uid: %s", ins.Uid))
		}
		ins.Uid = pkt.Msg.GetSID()
		ins.state = CONNECTED
		// 向路由模块的主消息通道压入用户上线消息
		ins.outchan <- pkt.Msg
	case CONNECTED:
		if ins.Uid != pkt.Msg.GetSID() {
			return errors.New(fmt.Sprintf("Uid and SID mismatched, uid: %s, sid: %s", ins.Uid, pkt.Msg.GetSID()))
		}
		ins.outchan <- pkt.Msg
	}
	return nil
}

// 客户端TCP连接读routine
func (ins *Instance) MainReadHandler() {
	defer ins.conn.Close()
	buffer := make([]byte, 0)
	for {
		ins.conn.SetReadDeadline(time.Now().Add(timeoutMinutes * time.Minute))
		tmp := make([]byte, 32)
		n, err := ins.conn.Read(tmp)
		if err != nil {
			log.Debug("Read: %s, uid: %s", err, ins.Uid)
			break
		}
		buffer = append(buffer, tmp[:n]...)
		pkt := new(protocol.Packet)
		var s bool
		// 解包Packet
		if s, err = pkt.UnPack(&buffer); s {
			if err = ins.stateMachine(pkt); err != nil {
				log.Debug("StateMachine: %s, uid: %s", err, ins.Uid)
				break
			}
		} else if err != nil {
			log.Debug("Unpack: %s, received bytes: %d, uid: %s", err, n, ins.Uid)
			break
		}
	}
	log.Info("Disconnected, uid: %s", ins.Uid)
	// 删除在线表中当前用户
	ins.onlinetable.DelEntity(ins.Uid)
}

// 客户端TCP连接写routine
func (ins *Instance) MainWriteHandler() {
	defer ins.conn.Close()
	for {
		msg := <-ins.inchan
		pkt := new(protocol.Packet)
		pkt.Msg = msg
		data := pkt.Pack()
		_, err := ins.conn.Write(data)
		nerr, ok := err.(*net.OpError)
		if ok && nerr.Temporary() {
			ins.inchan <- msg
			continue
		}
		if err != nil {
			ins.inchan <- msg
			break
		}
	}
	log.Info("Disconnected, uid: %s", ins.Uid)
	// 删除在线表中当前用户
	ins.onlinetable.DelEntity(ins.Uid)
}
