package main

import (
	"math/rand"
	"time"

	"github.com/denghongcai/messagehive/message"
	"github.com/denghongcai/messagehive/protocol"
	"github.com/golang/protobuf/proto"
)

const (
	RESP = iota
	AUTH
	EXCHANGE
	HEARTBEAT
	JOINGROUP
)

func setBit(n int, pos uint) int {
	n |= (1 << pos)
	return n
}

func CreatePacket(msgtype int, sid string, rid string, body string) []byte {
	pkt := new(protocol.Packet)
	pkt.Msg = new(message.Container)
	pkt.Msg.SID = proto.String(sid)
	pkt.Msg.RID = proto.String(rid)
	switch msgtype {
	case AUTH:
		msgtype = 0
		msgtype = setBit(msgtype, 0)
		msgtype = setBit(msgtype, 1)
	case JOINGROUP:
		msgtype = 0
		msgtype = setBit(msgtype, 0)
		msgtype = setBit(msgtype, 5)
	case EXCHANGE:
		msgtype = 0
		msgtype = setBit(msgtype, 0)
	}
	pkt.Msg.TYPE = proto.Uint32(uint32(msgtype))
	pkt.Msg.STIME = proto.Int64(time.Now().Unix())
	pkt.Msg.BODY = proto.String(body)
	msg := pkt.Pack()
	return msg
}

var letters = []rune("abcdefghipqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
