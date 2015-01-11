package main

import (
	"github.com/denghongcai/generalmessagegate/message"
	"github.com/denghongcai/generalmessagegate/protocol"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"time"
)

const (
	RESP = iota
	SYS
	EXCHANGE
	HEARTBEAT
	CREATGROUP
)

func CreatePacket(msgtype int, sid string, rid string, body string) []byte {
	pkt := new(protocol.Packet)
	pkt.Msg = new(message.Container)
	pkt.Msg.SID = proto.String(sid)
	pkt.Msg.RID = proto.String(rid)
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
