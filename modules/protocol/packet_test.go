package protocol

import (
	"testing"

	"github.com/denghongcai/MessageHive/modules/message"
	"github.com/golang/protobuf/proto"
)

func TestPackUnPack(t *testing.T) {
	pkt := new(Packet)
	pkt.Msg = &message.Container{
		SID:   proto.String("foo"),
		RID:   proto.String("pig"),
		TYPE:  proto.Uint32(0),
		STIME: proto.Int64(0),
		BODY:  proto.String("haha"),
	}
	msg := pkt.Pack()
	newpkt := new(Packet)
	newpkt.UnPack(&msg)
	if pkt.Msg.GetSID() != newpkt.Msg.GetSID() {
		t.Errorf("Data mismatch")
	}
}
