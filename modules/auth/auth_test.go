package auth

import (
	"testing"

	"github.com/denghongcai/MessageHive/modules/message"
	"github.com/golang/protobuf/proto"
)

func TestAuth(t *testing.T) {
	SetAuthHandler("test", `{}`)
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
