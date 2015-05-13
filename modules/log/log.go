package log

import "github.com/op/go-logging"

var logLevels = map[string]logging.Level{
	"Notice":   logging.NOTICE,
	"Debug":    logging.DEBUG,
	"Info":     logging.INFO,
	"Warn":     logging.WARNING,
	"Error":    logging.ERROR,
	"Critical": logging.CRITICAL,
}
var (
	loggers []*logging.Logger
)

type loggerType func(string) (*logging.Logger, error)

var adapters = make(map[string]loggerType)

func Register(name string, log loggerType) {
	if log == nil {
		panic("log: register provider is nil")
	}
	if _, dup := adapters[name]; dup {
		panic("log: register called twice for provider \"" + name + "\"")
	}
	adapters[name] = log
}
func NewLogger(mode, config string) {
	if log, ok := adapters[mode]; ok {
		logger, err := log(config)
		if err != nil {
			panic("log: failed on adding \"" + mode + "\" provider")
		}
		loggers = append(loggers, logger)
	}
}
func Debug(format string, v ...interface{}) {
	for _, logger := range loggers {
		logger.Debug(format, v...)
	}
}
func Notice(format string, v ...interface{}) {
	for _, logger := range loggers {
		logger.Notice(format, v...)
	}
}
func Info(format string, v ...interface{}) {
	for _, logger := range loggers {
		logger.Info(format, v...)
	}
}
func Warn(format string, v ...interface{}) {
	for _, logger := range loggers {
		logger.Warning(format, v...)
	}
}
func Error(format string, v ...interface{}) {
	for _, logger := range loggers {
		logger.Error(format, v...)
	}
}
func Critical(format string, v ...interface{}) {
	for _, logger := range loggers {
		logger.Critical(format, v...)
	}
}
func Fatal(v ...interface{}) {
	for _, logger := range loggers {
		logger.Fatal(v...)
	}
}
func Fatalf(format string, v ...interface{}) {
	for _, logger := range loggers {
		logger.Fatalf(format, v...)
	}
}
