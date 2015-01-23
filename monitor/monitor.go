package monitor

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/denghongcai/generalmessagegate/onlinetable"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("main")

type Config struct {
	address     string
	onlineTable *onlinetable.Container
}

func NewConfig(address string, onlinetable *onlinetable.Container) Config {
	return Config{
		address:     address,
		onlineTable: onlinetable,
	}
}

type onlineTableStat struct {
	TotalOnlineNum int
	OnlineGroupNum int
	OnlineSysNum   int
	OnlineUserNum  int
}

func setUpRoutes(config Config) {
	http.Handle("/stats", newStaticsHandler(config))
}

func Start(config Config) error {
	setUpRoutes(config)
	log.Info("Monitor listen on %s", config.address)
	err := http.ListenAndServe(config.address, nil)
	if err != nil {
		return err
	}
	return nil
}

type staticsHandler struct {
	onlineTable *onlinetable.Container
}

func newStaticsHandler(config Config) staticsHandler {
	handler := staticsHandler{
		onlineTable: config.onlineTable,
	}
	return handler
}

func (handler staticsHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	data, err := json.Marshal(onlineTableStat{1, 1, 1, 1})
	if err != nil {
		log.Error(err.Error())
	}
	writer.Header().Set("Content-Type", "application/json")
	fmt.Fprint(writer, string(data))
}
