package main

import (
	"bytes"
	"flag"
	"github.com/coding-boot-camp/nexus/services/tkt"
	"github.com/gabrielmorenobrc/queue/api"
	"net/http"
)

type Config struct {
	DatabaseConfig tkt.DatabaseConfig `json:"databaseConfig"`
	WorkerInterval int                `json:"workerInterval"`
	MaxErrorCount  int                `json:"maxErrorCount"`
	LogToConsole   bool               `json:"logToConsole"`
	LogTags        []string           `json:"logTags"`
}

var config Config

var conf = flag.String("conf", "conf.json", "Config")

func main() {
	flag.Parse()
	tkt.LoadConfig(*conf, &config)

	tkt.ConfigLoggers("lss.log", 2000000, 10, config.LogToConsole, config.LogTags...)

	tkt.ExecuteTransactional(config.DatabaseConfig, func(txCtx *tkt.TxCtx, args ...interface{}) interface{} {
		api := queue.NewApi(txCtx)
		api.Post("harness", "Right")
		api.Post("harness", "Wrong")
		return nil
	})

	worker := queue.NewWorker(config.DatabaseConfig, "harness", config.MaxErrorCount, config.WorkerInterval, func(txCtx *tkt.TxCtx, entry queue.Entry) {
		value := ""
		buffer := bytes.NewBuffer(entry.Data)
		tkt.JsonDecode(&value, buffer)
		println(value)
		if value == "Right" {
			println("Hi mom")
		} else {
			panic("The wrong one")
		}

	})
	worker.Start()

	http.ListenAndServe(":9999", nil)

}
