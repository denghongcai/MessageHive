package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/denghongcai/messagehive/message"
	"github.com/denghongcai/messagehive/monitor"
	"github.com/denghongcai/messagehive/onlinetable"
	"github.com/denghongcai/messagehive/router"
	"github.com/denghongcai/messagehive/rpc"
	"github.com/denghongcai/messagehive/server"
	"github.com/op/go-logging"
)

var (
	logLevel = flag.String("l", "debug", "log level")
	help     = flag.Bool("h", false, "help")
	port     = flag.Int("p", 1430, "port")
)

var log = logging.MustGetLogger("main")

var format = logging.MustStringFormatter(
	"%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level} %{id}%{color:reset} %{message}",
)

func setLoggerOption(log *logging.Logger, logLevel *string) {
	logBackend := logging.NewLogBackend(os.Stderr, "", 0)

	logBackendFormatter := logging.NewBackendFormatter(logBackend, format)

	logBackendLeveled := logging.AddModuleLevel(logBackendFormatter)

	switch *logLevel {
	case "debug":
		logBackendLeveled.SetLevel(logging.DEBUG, "")
	case "info":
		logBackendLeveled.SetLevel(logging.INFO, "")
	}

	logging.SetBackend(logBackendLeveled)

	log.Info("Current log level is %s", *logLevel)
}

func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: -l (debug|info)\n")
	}
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	setLoggerOption(log, logLevel)

	onlineTable := onlinetable.NewContainer()
	mainChan := make(chan *message.Container, 1024)

	serverAddress := []string{"0.0.0.0", fmt.Sprintf("%d", *port)}
	serverConfig := server.NewConfig(strings.Join(serverAddress, ":"), mainChan, onlineTable)
	go func() {
		if err := server.Handler(serverConfig); err != nil {
			log.Error(err.Error())
		}
	}()

	routerConfig := router.NewConfig(mainChan, onlineTable)
	go func() {
		if err := router.Handler(routerConfig); err != nil {
			log.Error(err.Error())
		}
	}()

	monitorAddress := []string{"0.0.0.0", "8888"}
	monitorConfig := monitor.NewConfig(strings.Join(monitorAddress, ":"), onlineTable)
	go func() {
		if err := monitor.Start(monitorConfig); err != nil {
			log.Error(err.Error())
		}
	}()

	rpcAddress := []string{"0.0.0.0", "9999"}
	rpcConfig := rpc.NewConfig(strings.Join(rpcAddress, ":"), onlineTable)
	if err := rpc.Start(rpcConfig); err != nil {
		log.Error(err.Error())
	}
}
