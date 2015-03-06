package rpc

import (
	"net/http"

	"github.com/denghongcai/messagehive/onlinetable"
)

type AddArgs struct {
	Uid     string
	Uidlist []string
}

type GroupManager struct {
	onlinetable *onlinetable.Container
}

func (t *GroupManager) AddGroup(r *http.Request, args *AddArgs, reply *int) error {
	log.Info("gogog")
	return t.onlinetable.AddGroupEntity(args.Uid, args.Uidlist)
}
