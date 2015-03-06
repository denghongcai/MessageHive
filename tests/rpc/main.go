package main

import (
	"log"
	"net"
	"net/rpc/jsonrpc"
)

type Args struct {
	Uid     string
	Uidlist []string
}

func main() {
	client, err := net.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	uidlist := []string{"bar"}

	args := &Args{"foo", uidlist}
	var reply int
	c := jsonrpc.NewClient(client)
	err = c.Call("GroupManager.AddGroup", args, &reply)
	if err != nil {
		log.Fatal("Add group error:", err)
	}
}
