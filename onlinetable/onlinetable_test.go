package onlinetable

import (
	"github.com/denghongcai/generalmessagegate/message"
	"testing"
)

func TestNewAddGetDel(t *testing.T) {
	container := NewContainer()
	container.AddEntity("foo", make(chan *message.Container))
	entity, err := container.GetEntity("foo")
	if entity.Uid != "foo" {
		t.Fail()
	}
	container.DelEntity("foo")
	entity, err = container.GetEntity("foo")
	if err == nil {
		t.Fail()
	}
}
