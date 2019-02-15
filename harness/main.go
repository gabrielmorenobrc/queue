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
		api.NewApi(txCtx).Post("harness", "Right")
		api.NewApi(txCtx).Post("harness", "Wrong")
		return nil
	})

	worker := api.NewWorker(config.DatabaseConfig, "harness", config.MaxErrorCount, config.WorkerInterval, func(txCtx *tkt.TxCtx, entry api.Entry) {
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
