package api

import (
	"bytes"
	"github.com/coding-boot-camp/nexus/services/tkt"
	"time"
)

type Api struct {
	txCtx *tkt.TxCtx
}

func (o *Api) Post(context string, data interface{}) *Entry {
	id := o.txCtx.Seq().Next("queue.entry")
	var json []byte
	buffer := bytes.Buffer{}
	tkt.JsonEncode(data, &buffer)
	json = buffer.Bytes()
	entity := Entry{Id: &id, Context: &context, Date: tkt.PTime(time.Now()), Data: json, ErrorCount: tkt.PInt(0)}
	o.txCtx.InsertEntity("queue", entity, false)
	return &entity
}

func (o *Api) ListPending(context string, maxErrorCount int) []Entry {
	return o.txCtx.QueryStruct(Entry{}, `select e.* from queue.entry e
		left outer join queue.success s on s.entry_id = e.id
		where context = $1 and errorcount < $2 and s.id is null order by e.date`, context, maxErrorCount).([]Entry)
}

func (o *Api) RegisterSuccess(entryId int64) {
	id := o.txCtx.Seq().Next("queue.success")
	entity := Success{Id: &id, EntryId: &entryId, Date: tkt.PTime(time.Now())}
	o.txCtx.InsertEntity("queue", entity, false)
}

func (o *Api) RegisterError(entryId int64, data interface{}) {
	id := o.txCtx.Seq().Next("queue.error")
	buffer := bytes.Buffer{}
	tkt.JsonEncode(data, &buffer)
	entity := Error{Id: &id, EntryId: &entryId, Date: tkt.PTime(time.Now()), Data: buffer.Bytes()}
	o.txCtx.InsertEntity("queue", entity, false)
	entry := o.txCtx.FindStruct(Entry{}, `select * from queue.entry where id = $1`, entryId).(*Entry)
	entry.ErrorCount = tkt.PInt(*entry.ErrorCount + 1)
	o.txCtx.UpdateEntity("queue", *entry)
}

func NewApi(txCtx *tkt.TxCtx) *Api {
	return &Api{txCtx: txCtx}
}
