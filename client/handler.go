package client

import (
	"errors"
	"github.com/denghongcai/generalmessagegate/auth"
	"github.com/denghongcai/generalmessagegate/message"
	"github.com/denghongcai/generalmessagegate/onlinetable"
	"github.com/denghongcai/generalmessagegate/protocol"
	"github.com/op/go-logging"
	"net"
	"time"
)

var log = logging.MustGetLogger("main")

const (
	AUTH = iota
	CONNECTED
)

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

func (ins *Instance) stateMachine(pkt *protocol.Packet) error {
	switch ins.state {
	case AUTH:
		if err := auth.Authenticate(pkt.Msg); err != nil {
			return err
		}
		if err := ins.onlinetable.AddEntity(pkt.Msg.GetSID(), ins.inchan); err != nil {
			return errors.New("Entity add failed")
		}
		ins.Uid = pkt.Msg.GetSID()
		ins.state = CONNECTED
	case CONNECTED:
		if ins.Uid != pkt.Msg.GetSID() {
			return errors.New("Uid and SID mismatched")
		}
		ins.outchan <- pkt.Msg
	}
	return nil
}

func (ins *Instance) MainReadHandler() {
	defer ins.conn.Close()
	buffer := make([]byte, 0)
	for {
		ins.conn.SetReadDeadline(time.Now().Add(timeoutMinutes * time.Minute))
		tmp := make([]byte, 32)
		n, err := ins.conn.Read(tmp)
		if err != nil {
			log.Debug("Read: %s", err)
			break
		}
		buffer = append(buffer, tmp[:n]...)
		pkt := new(protocol.Packet)
		var s bool
		if s, err = pkt.UnPack(&buffer); s {
			if err = ins.stateMachine(pkt); err != nil {
				log.Debug("StateMachine: %s", err)
				break
			}
		} else if err != nil {
			log.Debug("Unpack: %s, received bytes: %d", err, n)
			break
		}
	}
	log.Info("Disconnected")
	ins.onlinetable.DelEntity(ins.Uid)
}

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
	log.Info("Disconnected")
	ins.onlinetable.DelEntity(ins.Uid)
}