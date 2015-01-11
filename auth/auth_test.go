package auth

import (
	"github.com/denghongcai/generalmessagegate/message"
	"github.com/golang/protobuf/proto"
	"testing"
)

func TestAuth(t *testing.T) {
	msg := &message.Container{
		SID:   proto.String("hehe"),
		RID:   proto.String("haha"),
		TYPE:  proto.Uint32(0),
		STIME: proto.Int64(0),
		BODY:  proto.String(`{"token":"hehe"}`),
	}
	if err := Authenticate(msg); err != nil {
		t.Error(err)
	}
	emsg := &message.Container{
		SID:   proto.String("hehe"),
		RID:   proto.String("haha"),
		TYPE:  proto.Uint32(0),
		STIME: proto.Int64(0),
		BODY:  proto.String(`{"tokens":"hehe"}`),
	}
	if err := Authenticate(emsg); err == nil {
		t.Fail()
	}
}