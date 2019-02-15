package api

import (
	"time"
)

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
