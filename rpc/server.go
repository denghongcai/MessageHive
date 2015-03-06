package rpc

import (
	"net/http"

	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"

	"github.com/denghongcai/messagehive/onlinetable"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("main")

type Config struct {
	address     string
	onlinetable *onlinetable.Container
}

func NewConfig(address string, onlinetable *onlinetable.Container) Config {
	return Config{
		address:     address,
		onlinetable: onlinetable,
	}
}

func Start(config Config) error {
	gm := &GroupManager{
		onlinetable: config.onlinetable,
	}
	server := rpc.NewServer()
	server.RegisterCodec(json.NewCodec(), "application/json")
	server.RegisterService(gm, "")
	http.Handle("/rpc", server)
	err := http.ListenAndServe(config.address, nil)
	return err
}
