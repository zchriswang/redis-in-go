package database

import (
	"redis-in-go/datastructure/dict"
	"redis-in-go/interface/database"
	"redis-in-go/interface/resp"
	"redis-in-go/resp/reply"
	"strings"
)

type DB struct {
	index int
	data  dict.Dict
}

func makeDB() *DB {
	db := &DB{
		data: dict.MakeSyncDict(),
	}
	return db
}

type ExecFunc func(db *DB, args [][]byte) resp.Reply

type CmdLine = [][]byte

func (d *DB) Exec(c resp.Connection, cmdLine CmdLine) resp.Reply {
	// PING SET SETNX
	cmdName := strings.ToLower(string(cmdLine[0]))
	cmd, ok := cmdTable[cmdName]
	if !ok {
		return reply.MakeStandardErrReply("ERR unknown command " + cmdName)
	}
	if !validateArity(cmd.arity, cmdLine) {
		return reply.MakeArgNumErrReply(cmdName)
	}
	fun := cmd.executor
	// SET K V
	return fun(d, cmdLine[1:])
}

// SET K V -> arity = 3
// EXISTS k1 k2 k3... arity = -2
func validateArity(arity int, cmdArgs [][]byte) bool {
	argNum := len(cmdArgs)
	if arity >= 0 {
		return argNum == arity
	}
	return argNum >= -arity
}

func (d *DB) GetEntity(key string) (*database.DataEntity, bool) {
	raw, ok := d.data.Get(key)
	if !ok {
		return nil, false
	}
	entity, _ := raw.(*database.DataEntity)
	return entity, true
}

func (db *DB) PutEntity(key string, entity *database.DataEntity) int {
	return db.data.Put(key, entity)
}

func (db *DB) PutIfExists(key string, entity *database.DataEntity) int {
	return db.data.PutIfExists(key, entity)
}

func (db *DB) PutIfAbsent(key string, entity *database.DataEntity) int {
	return db.data.PutIfAbsent(key, entity)
}

func (db *DB) Remove(key string) {
	db.data.Remove(key)
}

func (db *DB) Removes(keys ...string) int {
	deleted := 0
	for _, key := range keys {
		_, exists := db.data.Get(key)
		if exists {
			db.Remove(key)
			deleted++
		}
	}
	return deleted
}

func (db *DB) Flush() {
	db.data.Clear()
}
