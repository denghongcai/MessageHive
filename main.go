// 主包，负责功能模块初始化
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/denghongcai/MessageHive/message"
	"github.com/denghongcai/MessageHive/monitor"
	"github.com/denghongcai/MessageHive/onlinetable"
	"github.com/denghongcai/MessageHive/router"
	"github.com/denghongcai/MessageHive/rpc"
	"github.com/denghongcai/MessageHive/server"
	"github.com/op/go-logging"
)

var (
	logLevel = flag.String("l", "debug", "log level")
	help     = flag.Bool("h", false, "help")
	port     = flag.Int("p", 1430, "port")
)

var log = logging.MustGetLogger("main")

// 日志输出格式
var format = logging.MustStringFormatter(
	"%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level} %{id}%{color:reset} %{message}",
)

// 日志配置
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

	// 连接监听模块初始化
	serverAddress := []string{"0.0.0.0", fmt.Sprintf("%d", *port)}
	serverConfig := server.NewConfig(strings.Join(serverAddress, ":"), mainChan, onlineTable)
	go func() {
		if err := server.Handler(serverConfig); err != nil {
			log.Error(err.Error())
		}
	}()

	// 路由模块初始化
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
