package yeelight

import (
	"encoding/json"
	"fmt"
	"math"
	"net"
	"time"
)

func (y *Bulb) ExecuteCommand(name string, params ...interface{}) (*CommandResult, error) {
	return y.execute(y.newCommand(name, params))
}

func (y *Bulb) newCommand(name string, params []interface{}) *Command {
	if len(params) > 0 {
		switch v := params[0].(type) {
		case []interface{}:
			params = v
		case []string:
			s := make([]interface{}, len(v))
			for i, val := range v {
				s[i] = val
			}
			params = s
		default:
		}
	}

	return &Command{
		Method: name,
		ID:     y.getCmdId(),
		Params: params,
	}
}

func (y *Bulb) execute(cmd *Command) (*CommandResult, error) {

	if y.conn == nil {
		conn, err := net.Dial("tcp", y.addr)

		if nil != err {
			return nil, fmt.Errorf("cannot open connection to %s. %s", y.addr, err)
		}
		y.conn = conn
	}
	//time.Sleep(time.Second)
	y.conn.SetReadDeadline(time.Now().Add(timeout))
	//write request/command
	b, _ := json.Marshal(cmd)
	fmt.Fprint(y.conn, string(b)+crlf)
	//wait and read for response
	_, err := y.conn.Write(b)
	if err != nil {
		return nil, fmt.Errorf("cannot read command result %s", err)
	}

	return nil, nil
}

func (y *Bulb) getCmdId() int {
	if y.cmdId == math.MaxInt32 {
		y.cmdId = 0
	}
	currentId := y.cmdId
	y.cmdId += 1
	return currentId
}
