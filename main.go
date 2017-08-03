package main

import (
	"leafmock/msg"
	"leafmock/network"

	"github.com/golang/protobuf/proto"
)

func main() {
	var mc = network.NewMockClient("127.0.0.1:3563", 1024)

	mc.Register(&msg.Login_S2C{}, handleLoginRes)
	mc.Register(&msg.RegisterName_S2C{}, handleRegisterNameRes)
	mc.Run()

	// login request
	mc.WriteMsg(&msg.Login_C2S{
		Sdk:   proto.String("mocksdk"),
		Uid:   proto.String("db6b6843e2f4"),
		Token: proto.String("260c74fd-4781-405f-9410-08beaa6264be"),
	})

}

func handleLoginRes(c interface{}, m interface{}) {
	a := c.(*network.MockClient)
	res := m.(*msg.Login_S2C)
	if res.GetCode() == msg.Login_S2C_YES {
		// if login success, register nickname
		a.WriteMsg(&msg.RegisterName_C2S{
			Nickname: proto.String("猴子请来的救比"),
		})
	}
}

func handleRegisterNameRes(c interface{}, m interface{}) {
	a := c.(*network.MockClient)
	res := m.(*msg.RegisterName_S2C)

	if res.GetCode() == msg.RegisterName_S2C_YES {
		// do something
		// a.Close()
		_ = a
	}
}
