package onlinetable

import (
	"fmt"
	"testing"

	"github.com/denghongcai/MessageHive/message"
)

func TestNewAddGetDel(t *testing.T) {
	container := NewContainer()
	container.AddEntity("foo", make(chan *message.Container))
	uidlist := make([]string, 0)
	uidlist = append(uidlist, "foo")
	container.AddGroupEntity("bar", uidlist)
	entity, err := container.GetEntity("foo")
	if entity.Uid != "foo" {
		t.Fail()
	}
	container.DelEntity("foo")
	entity, err = container.GetEntity("foo")
	if err == nil {
		t.Fail()
	}
	entity, err = container.GetEntity("bar")
	if entity.Uid != "bar" {
		t.Fail()
	}
	fmt.Printf("Uid list: %v\n", entity.List)
}
