package database

import (
	"redis/interface/resp"
	"redis/resp/reply"
	"regexp"
)

func execDel(db *DB, args [][]byte) resp.Reply {

	keys := make([]string, len(args))
	for i, arg := range args {
		keys[i] = string(arg)
	}
	deleted := db.Removes(keys...)
	return reply.MakeIntReply(int64(deleted))
}

// execExists checks if a is existed in db
func execExists(db *DB, args [][]byte) resp.Reply {
	result := int64((0))
	for _, key := range args {
		if _, ok := db.GetEntity(string(key)); ok {
			result++
		}
	}
	return reply.MakeIntReply(result)
}
func execFlushDB(db *DB, args [][]byte) resp.Reply {

	db.Flush()
	return &reply.OkReply{}
}
func execType(db *DB, args [][]byte) resp.Reply {
	entity, ok := db.GetEntity(string(args[0]))
	if !ok {
		return reply.MakeStatusReply("none")
	}
	switch entity.Data.(type) {
	case []byte:
		return reply.MakeStatusReply("string")
	}
	return &reply.UnknownErrReply{}
}
func execRename(db *DB, args [][]byte) resp.Reply {
	old := string(args[0])
	entity, ok := db.GetEntity(old)
	if !ok {
		return reply.MakeErrReply("no such key")

	}
	db.Remove(old)
	new := string(args[1])
	db.PutEntity(new, entity)
	return reply.MakeIntReply(int64(1))
}

// execRenameNx a key, only if the new key does not exist
func execRenameNx(db *DB, args [][]byte) resp.Reply {
	src := string(args[0])
	dest := string(args[1])

	_, ok := db.GetEntity(dest)
	if ok {
		return reply.MakeIntReply(0)
	}

	entity, ok := db.GetEntity(src)
	if !ok {
		return reply.MakeErrReply("no such key")
	}
	db.Removes(src, dest) // clean src and dest with their ttl
	db.PutEntity(dest, entity)
	return reply.MakeIntReply(1)
}

// execKeys returns all keys matching the given pattern
func execKeys(db *DB, args [][]byte) resp.Reply {

	result := make([][]byte, 0)
	pattern := string(args[0])
	var match bool
	for i := 1; i < len(args); i++ {
		match, _ = regexp.MatchString(pattern, string(args[i]))
		if match {
			result = append(result, args[i])
		}
	}
	return reply.MakeMultiBulkReply(result)

}

func init() {
	RegisterCommand("Del", execDel, -2)
	RegisterCommand("Exists", execExists, -2)
	RegisterCommand("Keys", execKeys, 2)
	RegisterCommand("FlushDB", execFlushDB, -1)
	RegisterCommand("Type", execType, 2)
	RegisterCommand("Rename", execRename, 3)
	RegisterCommand("RenameNx", execRenameNx, 3)
}
