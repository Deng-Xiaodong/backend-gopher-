package aof

import (
	"io"
	"os"
	"redis/config"
	"redis/interface/database"
	"redis/lib/logger"
	"redis/lib/utils"
	"redis/resp/connection"
	"redis/resp/parser"
	"redis/resp/reply"
	"strconv"
)

const (
	aofQueueSize = 1 << 16
)

type CmdLine = [][]byte
type payload struct {
	cmdLine CmdLine
	dbIndex int
}
type AofHandler struct {
	mdb       database.Database
	currentDB int
	aofFile   *os.File
	fileName  string
	aofChan   chan *payload
}

func NewAofHandler(mdb database.Database) (*AofHandler, error) {

	handle := &AofHandler{}
	handle.mdb = mdb
	handle.fileName = config.Properties.AppendFilename
	handle.loadAof()
	aofFile, err := os.OpenFile(handle.fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	handle.aofFile = aofFile
	handle.aofChan = make(chan *payload, aofQueueSize)
	go func() {
		handle.handleAof()
	}()
	return handle, nil
}

func (aof *AofHandler) loadAof() {
	file, err := os.Open(aof.fileName)
	if err != nil {
		logger.Warn(err)
		return
	}
	defer file.Close()

	ch := parser.ParseStream(file)
	fakeConn := &connection.Connection{}
	for p := range ch {
		go func() {
			aof.exec(p, fakeConn)
		}()
	}
}

func (aof *AofHandler) exec(p *parser.Payload, fakeConn *connection.Connection) {
	if p.Err != nil {
		if p.Err == io.EOF {
			return
		}
		logger.Error("parse error: " + p.Err.Error())

	}
	if p.Data == nil {
		logger.Error("empty payload")

	}
	r, ok := p.Data.(*reply.MultiBulkReply)
	if !ok {
		logger.Error("require multi bulk reply")

	}
	ret := aof.mdb.Exec(fakeConn, r.Args)
	if reply.IsErrorReply(ret) {
		logger.Error("exec err")
	}
}

func (aof *AofHandler) AddAof(dbIndex int, line CmdLine) {
	if config.Properties.AppendOnly && aof.aofChan != nil {
		aof.aofChan <- &payload{
			cmdLine: line,
			dbIndex: dbIndex,
		}
	}
}

func (aof *AofHandler) handleAof() {
	aof.currentDB = 0
	for p := range aof.aofChan {
		if p.dbIndex != aof.currentDB {
			data := reply.MakeMultiBulkReply(utils.ToCmdLine("SELECT", strconv.Itoa(p.dbIndex))).ToBytes()
			_, err := aof.aofFile.Write(data)
			if err != nil {
				logger.Warn(err)
				continue
			}
			aof.currentDB = p.dbIndex
		}
		data := reply.MakeMultiBulkReply(p.cmdLine).ToBytes()
		_, err := aof.aofFile.Write(data)
		if err != nil {
			logger.Warn(err)
		}
	}
}
