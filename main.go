// 主包，负责功能模块初始化
package main

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/denghongcai/MessageHive/modules/command"
	"github.com/denghongcai/MessageHive/modules/message"
	"github.com/denghongcai/MessageHive/modules/server"
	"github.com/denghongcai/MessageHive/modules/onlinetable"
	"github.com/denghongcai/MessageHive/modules/router"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	command.Execute()

	onlineTable := onlinetable.NewContainer()
	mainChan := make(chan *message.Container, 1024)

	// 连接监听模块初始化
	port := 1430
	serverAddress := []string{"0.0.0.0", fmt.Sprintf("%d", port)}
	serverConfig := server.NewConfig(strings.Join(serverAddress, ":"), mainChan, onlineTable)
	go func() {
		if err := server.Handler(serverConfig); err != nil {
			//			log.Error(err.Error())
		}
	}()

	// 路由模块初始化
	routerConfig := router.NewConfig(mainChan, onlineTable)
	go func() {
		if err := router.Handler(routerConfig); err != nil {
			//			log.Error(err.Error())
		}
	}()

	/*
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
	*/
}
