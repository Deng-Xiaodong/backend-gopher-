package parser

import (
	"bufio"
	"errors"
	"io"
	"redis/interface/resp"
	"redis/lib/logger"
	"redis/resp/reply"
	"runtime/debug"
	"strconv"
	"strings"
)

type Payload struct {
	Data resp.Reply
	Err  error
}

var cache = 100

// ParseStream reads data from io.Reader and send payloads through channel
func ParseStream(reader io.Reader) chan *Payload {
	ch := make(chan *Payload, cache)
	go parse0(reader, ch)
	return ch
}

type readState struct {
	readingMultiLine  bool
	expectedArgsCount int
	msgType           byte
	args              [][]byte
	bulkLen           int64
}

func (s *readState) finished() bool {
	return s.expectedArgsCount > 0 && len(s.args) == s.expectedArgsCount
}

func parse0(reader io.Reader, ch chan *Payload) {

	defer func() {
		if err := recover(); err != nil {
			logger.Error(string(debug.Stack()))
		}
	}()
	bufReader := bufio.NewReader(reader)
	var state readState
	var err error
	var msg []byte
	for {
		//read line
		var ioErr bool
		msg, ioErr, err = readLine(bufReader, &state)
		if err != nil {
			// encounter io err, stop read
			if ioErr {
				ch <- &Payload{
					Err: err,
				}
				close(ch)
				return
			}
			// protocol err, reset read state
			ch <- &Payload{
				Err: err,
			}
			state = readState{}
			continue
		}
		//parse line
		if !state.readingMultiLine {
			//如果不处于多行模式则有三种可能
			//- 确实这一行只有一个\r\n(+ - :)
			//- 多行字符串模式未初始化 $
			//- 数组模式未初始化 *
			if msg[0] == '*' {
				//进行数组模式初始化
				err = parseMultiBulkHeader(msg, &state)

				//协议错误
				if err != nil {
					ch <- &Payload{
						Err: err,
					}
					state = readState{}
					continue
				}
				//特殊情况
				if state.expectedArgsCount == 0 {
					ch <- &Payload{
						//这里注意一个细节：Payload的Data是一个接口类型，
						//它能接收所有实现了它的结构体，但是如果接收者是该结构体的指针时必须赋值为该结构体的指针
						Data: &reply.EmptyMultiBulkReply{},
					}
					state = readState{}
					continue
				}
			} else if msg[0] == '$' {
				//进行多行字符串模式初始化
				err = parseBulkHeader(msg, &state)
				//协议错误
				if err != nil {
					ch <- &Payload{
						Err: err,
					}
					state = readState{}
					continue
				}
				//特殊情况
				if state.bulkLen == -1 {
					ch <- &Payload{
						Data: &reply.NullBulkReply{},
					}
					state = readState{}
					continue
				}
			} else {
				//单行模式
				result, err := parseSingleLineReply(msg)
				ch <- &Payload{
					Err:  err,
					Data: result,
				}
				state = readState{} // reset state
				continue
			}
		} else {
			//进入多行或者数组模式
			err = readBody(msg, &state)
			if err != nil {
				ch <- &Payload{
					Err: errors.New("protocol error: " + string(msg)),
				}
				state = readState{} // reset state
				continue
			}
			// if sending finished
			if state.finished() {
				var result resp.Reply
				if state.msgType == '*' {
					result = reply.MakeMultiBulkReply(state.args)
				} else if state.msgType == '$' {
					result = reply.MakeBulkReply(state.args[0])
				}
				ch <- &Payload{
					Data: result,
					Err:  err,
				}
				state = readState{}
			}
		}
	}

}

func readLine(bufReader *bufio.Reader, state *readState) ([]byte, bool, error) {
	var msg []byte
	var err error
	if state.bulkLen == 0 { // read normal line
		msg, err = bufReader.ReadBytes('\n')
		if err != nil {
			return nil, true, err
		}
		if len(msg) == 0 || msg[len(msg)-2] != '\r' {
			return nil, false, errors.New("protocol error: " + string(msg))
		}
	} else { // read bulk line (binary safe)
		msg = make([]byte, state.bulkLen+2)
		_, err = io.ReadFull(bufReader, msg)
		if err != nil {
			return nil, true, err
		}
		if len(msg) == 0 ||
			msg[len(msg)-2] != '\r' ||
			msg[len(msg)-1] != '\n' {
			return nil, false, errors.New("protocol error: " + string(msg))
		}
		state.bulkLen = 0
	}
	return msg, false, nil
}

//初始化解析器，再根据解析器的状态去得到最终reply
//"*"初始化
func parseMultiBulkHeader(msg []byte, state *readState) error {
	var err error
	var expectedLine int64
	expectedLine, err = strconv.ParseInt(string(msg[1:len(msg)-2]), 10, 32)
	if err != nil {
		return errors.New("protocol error: " + string(msg))
	}
	if expectedLine == 0 {
		state.expectedArgsCount = 0
		return nil
	} else if expectedLine > 0 {
		// first line of multi bulk reply
		state.msgType = msg[0]
		state.readingMultiLine = true
		state.expectedArgsCount = int(expectedLine)
		state.args = make([][]byte, 0, expectedLine)
		return nil
	} else {
		return errors.New("protocol error: " + string(msg))
	}
}

//"$"初始化
func parseBulkHeader(msg []byte, state *readState) error {
	var err error
	state.bulkLen, err = strconv.ParseInt(string(msg[1:len(msg)-2]), 10, 64)
	if err != nil {
		return errors.New("protocol error: " + string(msg))
	}
	if state.bulkLen == -1 {
		return nil
	}
	if state.bulkLen > 0 {
		state.msgType = msg[0]
		state.readingMultiLine = false
		state.expectedArgsCount = 1
		state.args = make([][]byte, 0, 1)
		return nil
	} else {
		return errors.New("protocol error: " + string(msg))
	}

}

//"+","-",":" 这种单行的直接解析为最终reply返回
func parseSingleLineReply(msg []byte) (resp.Reply, error) {
	str := strings.TrimSuffix(string(msg), "\r\n")
	var result resp.Reply
	switch msg[0] {
	case '+': // status reply
		result = reply.MakeStatusReply(str[1:])
	case '-': // err reply
		result = reply.MakeErrReply(str[1:])
	case ':': // int reply
		val, err := strconv.ParseInt(str[1:], 10, 64)
		if err != nil {
			return nil, errors.New("protocol error: " + string(msg))
		}
		result = reply.MakeIntReply(val)
	default:
		// parse as text protocol
		strs := strings.Split(str, " ")
		args := make([][]byte, len(strs))
		for i, s := range strs {
			args[i] = []byte(s)
		}
		result = reply.MakeMultiBulkReply(args)
	}
	return result, nil
}

func readBody(msg []byte, state *readState) error {
	line := msg[0 : len(msg)-2]
	var err error

	if line[0] == '$' {
		state.bulkLen, err = strconv.ParseInt(string(line[1:]), 10, 64)
		if err != nil {
			return errors.New("protocol error: " + string(msg))
		}
		if state.bulkLen <= 0 {
			state.args = append(state.args, []byte{})
			state.bulkLen = 0
		}
	} else {
		state.args = append(state.args, line)
	}
	return nil
}
