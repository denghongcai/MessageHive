package log

import (
	"encoding/json"
	"github.com/op/go-logging"
	"os"
)

type consoleLoggerConfig struct {
	Level string `json:"level"`
}

func setLoggerOption(config consoleLoggerConfig) *logging.Logger {
	log := logging.MustGetLogger("console")
	format := logging.MustStringFormatter("%{color}%{time:15:04:05.000} ▶ %{level} %{id}%{color:reset} %{message}")
	logBackend := logging.NewLogBackend(os.Stderr, "", 0)
	logBackendFormatter := logging.NewBackendFormatter(logBackend, format)
	logBackendLeveled := logging.AddModuleLevel(logBackendFormatter)
	logBackendLeveled.SetLevel(logLevels[config.Level], "")
	logging.SetBackend(logBackendLeveled)
	return log
}
func NewConsole(config string) (*logging.Logger, error) {
	loggerConfig := &consoleLoggerConfig{}
	err := json.Unmarshal([]byte(config), loggerConfig)
	if err != nil {
		return nil, err
	}
	return setLoggerOption(*loggerConfig), nil
}
func init() {
	Register("console", NewConsole)
}
