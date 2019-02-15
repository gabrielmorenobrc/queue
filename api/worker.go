package api

import (
	"github.com/coding-boot-camp/nexus/services/tkt"
	"runtime/debug"
	"time"
)

type Callback func(txCtx *tkt.TxCtx, entry Entry)

type Worker struct {
	config   Config
	context  string
	callback Callback
}

func (o *Worker) Start() {

	ticker := time.NewTicker(time.Second * time.Duration(o.config.WorkerInterval))
	go func() {
		for range ticker.C {
			o.processTicker()
		}
	}()

}

func (o *Worker) processTicker() {
	defer catchPanic()
	tkt.Logger("debug").Printf("Processing %s...", o.context)
	entries := o.listPending()
	o.processEntries(entries)
	tkt.Logger("debug").Printf("Finished.")
}

func (o *Worker) processEntries(entries []Entry) {
	for i := range entries {
		entry := entries[i]
		tkt.Logger("debug").Printf("Processing entry %d", *entry.Id)
		o.processEntry(entry)
	}
}

func (o *Worker) processEntry(entry Entry) {
	defer o.processPanic(entry)
	tkt.ExecuteTransactional(o.config.DatabaseConfig, func(txCtx *tkt.TxCtx, args ...interface{}) interface{} {
		o.callback(txCtx, entry)
		NewApi(txCtx).RegisterSuccess(*entry.Id)
		return nil
	})
}

func (o *Worker) processPanic(entry Entry) {
	if r := recover(); r != nil {
		tkt.ExecuteTransactional(o.config.DatabaseConfig, func(txCtx *tkt.TxCtx, args ...interface{}) interface{} {
			NewApi(txCtx).RegisterError(*entry.Id, r)
			return nil
		})
	}
}

func (o *Worker) listPending() []Entry {
	return tkt.ExecuteTransactional(o.config.DatabaseConfig, func(txCtx *tkt.TxCtx, args ...interface{}) interface{} {
		return NewApi(txCtx).ListPending(o.context, o.config.MaxErrorCount)
	}).([]Entry)
}

func NewWorker(config Config, context string, callback Callback) *Worker {
	return &Worker{config: config, context: context, callback: callback}
}

func catchPanic() {
	if r := recover(); r != nil {
		tkt.Logger("error").Printf("%s", r)
		debug.PrintStack()
	}
}
