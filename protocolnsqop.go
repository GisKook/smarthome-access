package sha

import (
	"github.com/bitly/go-nsq"
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

type NsqOpPacket struct {
	GatewayID    uint64
	SerialNumber uint32
	DeviceID     uint64
	Op           uint8
}

func (p *NsqOpPacket) Serialize() []byte {
}

func ParseNsqOp(gatewayid uint64, serialnum uint32) *NsqOpPacket {

}
