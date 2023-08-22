package database

import (
	"redis/datastruct/dict"
	"redis/interface/database"
	"redis/interface/resp"
	"redis/resp/reply"
	"strings"
)

type DB struct {
	index  int
	data   dict.Dict
	addAof func(line CmdLine)
}

//有等号和无等号
type ExecFunc func(db *DB, args [][]byte) resp.Reply
type CmdLine = [][]byte

func makeDB() *DB {
	return &DB{
		data:   dict.MakeSyncDict(),
		addAof: func(line CmdLine) {},
	}
}

func (db *DB) Exec(cmdLine CmdLine) resp.Reply {
	cmdName := strings.ToLower(string(cmdLine[0]))
	cmd, ok := CmdTable[cmdName]
	if !ok {
		return reply.MakeErrReply("ERR unknown command '" + cmdName + "'")
	}
	if !validateArity(cmd.arity, cmdLine) {
		return reply.MakeArgNumErrReply(cmdName)
	}
	fun := cmd.executor
	return fun(db, cmdLine[1:])
}
func validateArity(arity int, cmdArgs [][]byte) bool {
	argNum := len(cmdArgs)
	if arity >= 0 {
		return argNum == arity
	}
	return argNum >= -arity
}

/* ---- data Access ----- */
func (db *DB) GetEntity(key string) (*database.DataEntity, bool) {
	raw, ok := db.data.Get(key)
	if !ok {
		return nil, false
	}
	return raw.(*database.DataEntity), true
}

func (db *DB) PutEntity(key string, val *database.DataEntity) int {
	return db.data.Put(key, val)
}

func (db *DB) PutIfExists(key string, val *database.DataEntity) int {
	return db.data.PutIfExists(key, val)
}

func (db *DB) PutIfAbsent(key string, val *database.DataEntity) int {
	return db.data.PutIfAbsent(key, val)
}

func (db *DB) Remove(key string) bool {
	return db.data.Remove(key)
}

func (db *DB) Removes(keys ...string) int {
	count := 0
	for _, key := range keys {
		if db.data.Remove(key) {
			count++
		}
	}
	return count
}
func (db *DB) Flush() {
	db.data.Clear()
}
