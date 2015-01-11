package main

import (
	"flag"
	"fmt"
	"github.com/denghongcai/generalmessagegate/message"
	"github.com/denghongcai/generalmessagegate/onlinetable"
	"github.com/denghongcai/generalmessagegate/router"
	"github.com/denghongcai/generalmessagegate/server"
	"github.com/op/go-logging"
	"os"
	"strings"
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

	address := []string{"0.0.0.0", fmt.Sprintf("%d", *port)}
	serverConfig := server.NewConfig(strings.Join(address, ":"), mainChan, onlineTable)
	go func() {
		if err := server.Handler(serverConfig); err != nil {
			log.Error(err.Error())
		}
	}()

	routerConfig := router.NewConfig(mainChan, onlineTable)
	if err := router.Handler(routerConfig); err != nil {
		log.Error(err.Error())
	}
}
