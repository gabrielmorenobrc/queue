package api

import (
	"github.com/coding-boot-camp/nexus/services/tkt"
	"time"
)

type Config struct {
	DatabaseConfig tkt.DatabaseConfig `json:"databaseConfig"`
	WorkerInterval int64              `json:"workerInterval"`
	MaxErrorCount  int                `json:"maxErrorCount"`
	LogToConsole   bool               `json:"logToConsole"`
	LogTags        []string           `json:"logTags"`
}

type Entry struct {
	Id         *int64     `json:"id"`
	Context    *string    `json:"context"`
	Date       *time.Time `json:"date"`
	Data       []byte     `json:"data"`
	ErrorCount *int       `json:"errorCount"`
}

type Success struct {
	Id      *int64     `json:"id"`
	EntryId *int64     `json:"entryId" sql:"entry_id"`
	Date    *time.Time `json:"date"`
}

type Error struct {
	Id      *int64     `json:"id"`
	EntryId *int64     `json:"entryId" sql:"entry_id"`
	Date    *time.Time `json:"date"`
	Data    []byte     `json:"data"`
}
