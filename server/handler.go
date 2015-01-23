package server

import (
	"crypto/tls"

	"github.com/denghongcai/generalmessagegate/client"
	"github.com/denghongcai/generalmessagegate/message"
	"github.com/denghongcai/generalmessagegate/onlinetable"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("main")

type Config struct {
	address     string
	mainchan    chan *message.Container
	onlinetable *onlinetable.Container
}

func NewConfig(address string, mainchan chan *message.Container, onlinetable *onlinetable.Container) Config {
	return Config{
		address:     address,
		mainchan:    mainchan,
		onlinetable: onlinetable,
	}
}

func Handler(config Config) error {
	tlsconfig, err := tlsConfig()
	if err != nil {
		return err
	}
	listener, err := tls.Listen("tcp", config.address, &tlsconfig)
	if err != nil {
		return err
	}
	log.Info("Server listen on: %s", config.address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		instance := client.NewInstance(conn, config.mainchan, config.onlinetable)
		instance.Handler()
	}
}
