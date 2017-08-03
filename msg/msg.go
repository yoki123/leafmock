package msg

import (
	"github.com/name5566/leaf/network/protobuf"
)

var Processor = protobuf.NewProcessor()

func init() {
	Processor.Register(&Login_C2S{})
	Processor.Register(&Login_S2C{})
	Processor.Register(&RegisterName_C2S{})
	Processor.Register(&RegisterName_S2C{})
}
