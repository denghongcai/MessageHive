package command

import (
	"github.com/denghongcai/MessageHive/modules/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var MessageHiveCmd = &cobra.Command{
	Use:   "MessageHive",
	Short: "MessageHive is a expressive, fast, full featured message gate.",
	Long:  "A expressivee, fast, full featured message gate lovely built by Hongcai Deng",
	Run: func(cmd *cobra.Command, args []string) {
		InitializeConfig()
	},
}

var messagehiveCmdV *cobra.Command

var CfgFile, LogLevel string

func Execute() {
	AddCommands()
	MessageHiveCmd.Execute()
}

func AddCommands() {
	MessageHiveCmd.AddCommand(serverCmd)
}

func init() {
	MessageHiveCmd.PersistentFlags().StringVarP(&CfgFile, "config", "c", "", "config file (default is path/config.yaml|json|toml)")
	MessageHiveCmd.PersistentFlags().StringVar(&LogLevel, "logLevel", "Info", "logout put level")
	messagehiveCmdV = MessageHiveCmd
}

func InitializeConfig() {
	viper.SetConfigFile(CfgFile)
	err := viper.ReadInConfig()
	if err != nil {
		//panic("Unable to locate Config file.Perhaps you need to create a config from example")
	}
	viper.SetDefault("LogLevel", "Info")

	if messagehiveCmdV.PersistentFlags().Lookup("logLevel").Changed {
		viper.Set("LogLevel", LogLevel)
	}

	log.NewLogger("console", `{"level": "`+viper.GetString("LogLevel")+`"}`)

	log.Info("Using config file: %s", viper.ConfigFileUsed())
}
