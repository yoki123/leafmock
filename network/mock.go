package network

import (
	"net"
	"reflect"
	"fmt"
	"github.com/name5566/leaf/log"
)

type MockClient struct {
	conn            *TCPConn
	processor       Processor
	functions       map[interface{}]interface{}
	address         string
	pendingWriteNum int
}

func NewMockClient(address string, pendingWriteNum int) *MockClient {
	mc := new(MockClient)
	mc.address = address
	mc.pendingWriteNum = pendingWriteNum
	mc.functions = make(map[interface{}]interface{})
	return mc
}

func (a *MockClient) SetProcessor(p Processor) {
	a.processor = p
}

func (a *MockClient) Connect() {
	conn, err := net.Dial("tcp", a.address)
	if err != nil {
		panic(err)
	}
	tcpConn := newTCPConn(conn, a.pendingWriteNum, NewMsgParser())
	a.conn = tcpConn
}

func (a *MockClient) Run() {
	for {
		data, err := a.conn.ReadMsg()
		if err != nil {
			log.Debug("read message: %v", err)
			break
		}

		if a.processor != nil {
			msg, err := a.processor.Unmarshal(data)
			if err != nil {
				log.Debug("unmarshal message error: %v", err)
				break
			}
			err = a.call(msg)
		}
	}
}

func (a *MockClient) WriteMsg(msg interface{}) {
	if a.processor != nil {
		data, err := a.processor.Marshal(msg)
		if err != nil {
			log.Error("marshal message %v error: %v", reflect.TypeOf(msg), err)
			return
		}
		err = a.conn.WriteMsg(data...)
		if err != nil {
			log.Error("write message %v error: %v", reflect.TypeOf(msg), err)
		}
	}
}

func (a *MockClient) Close() {
	a.conn.Close()
}

func (a *MockClient) Destroy() {
	a.conn.Destroy()
}

func (a *MockClient) call(arg interface{}) (err error) {
	id := reflect.TypeOf(arg)
	f := a.functions[id]
	if f == nil {
		err = fmt.Errorf("function id %v: function not registered", id)
		return
	}
	f.(func(interface{}, interface{}))(a, arg)
	return nil
}

func (a *MockClient) Register(m interface{}, f interface{}) {
	switch f.(type) {
	case func(interface{}, interface{}):
	default:
		panic(fmt.Sprintf("function id %v: definition of function is invalid", m))
	}
	id := reflect.TypeOf(m)

	if _, ok := a.functions[id]; ok {
		panic(fmt.Sprintf("function id %v: already registered", id))
	}
	a.functions[id] = f
}
