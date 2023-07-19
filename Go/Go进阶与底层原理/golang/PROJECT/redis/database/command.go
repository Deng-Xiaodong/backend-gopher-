package database

import "strings"

type command struct {
	executor ExecFunc
	arity    int //allow number of args, arity < 0 means len(args) >= -arity
}

var CmdTable = make(map[string]*command)

func RegisterCommand(name string, executor ExecFunc, arity int) {
	name = strings.ToLower(name)
	CmdTable[name] = &command{
		executor: executor,
		arity:    arity,
	}
}
