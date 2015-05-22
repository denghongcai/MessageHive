package command

import (
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/denghongcai/MessageHive/modules/auth"
	"github.com/denghongcai/MessageHive/modules/log"
	"github.com/denghongcai/MessageHive/modules/message"
	"github.com/denghongcai/MessageHive/modules/onlinetable"
	"github.com/denghongcai/MessageHive/modules/router"
	"github.com/denghongcai/MessageHive/modules/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serverPort int
var serverInterface string

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run MessageHive server",
	Long:  "Run Messagehive server",
}

func init() {
	serverCmd.Flags().IntVarP(&serverPort, "port", "p", 1430, "port on which the server will listen")
	serverCmd.Flags().StringVar(&serverInterface, "bind", "0.0.0.0", "interface to which the server will bind")
	serverCmd.Run = serve
}

func serve(cmd *cobra.Command, args []string) {
	InitializeConfig()

	viper.SetDefault("port", 1430)
	viper.SetDefault("bind", "0.0.0.0")

	if cmd.Flags().Lookup("port").Changed {
		viper.Set("port", serverPort)
	}
	serverPort = viper.GetInt("port")

	if cmd.Flags().Lookup("bind").Changed {
		viper.Set("bind", serverInterface)
	}
	serverInterface = viper.GetString("bind")

	auth.SetAuthHandler(viper.GetString("auth.adapter"), `{}`) // 设置认证方式

	onlineTable := onlinetable.NewContainer()       // 在线表初始化
	mainChan := make(chan *message.Container, 1024) // 主内部消息队列

	serverAddress := []string{serverInterface, fmt.Sprintf("%d", serverPort)}
	serverConfig := server.NewConfig(strings.Join(serverAddress, ":"), mainChan, onlineTable)
	go func() {
		if err := server.Handler(serverConfig); err != nil {
			log.Fatal("Server: %s", err.Error())
		}
	}()

	routerConfig := router.NewConfig(mainChan, onlineTable)
	go func() {
		if err := router.Handler(routerConfig); err != nil {
			log.Fatal("Server: %s", err.Error())
		}
	}()

	cleanupDone := make(chan bool)
	go func(cleanupDone chan bool) {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt)
		for _ = range signalChan {
			cleanupDone <- true
		}
	}(cleanupDone)
	<-cleanupDone

	log.Info("Server: gracefully shutdown...")
}
