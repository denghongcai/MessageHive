package command

import (
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
	serverCmd.Run = server
}

func server(cmd *cobra.Command, args []string) {
	InitializeConfig()

	viper.SetDefault("Port", 1430)
	viper.SetDefault("Bind", "0.0.0.0")

	if cmd.Flags().Lookup("port").Changed {
		viper.Set("Port", serverPort)
	}

	if cmd.Flags().Lookup("bind").Changed {
		viper.Set("Bind", serverInterface)
	}
}
