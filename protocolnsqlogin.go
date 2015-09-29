package sha

import (
	"github.com/bitly/go-nsq"
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

var (
	Offline uint8 = 0
	Online  uint8 = 1
	UnAuth  uint8 = 2
)

type NsqLoginPacket struct {
	GatewayID    uint64
	SerialNumber uint32
	Result       uint8
}

func (p *NsqLoginPacket) Serialize() []byte {
	rep := &Report.ControlReport{
		Tid:          p.GatewayID,
		SerialNumber: p.SerialNumber,
		Command: &Report.Command{
			Type: Report.Command_CMT_REPLOGIN,
			Paras: []*Report.Command_Param{
				&Report.Command_Param{
					Type:  Report.Command_Param_UINT8,
					Npara: 1,
				},
			},
		},
	}

	repdata, err := proto.Marshal(rep)
	if err != nil {
		log.Println("marshaling error", err)
	}

	return rep
}

func ParseNsqLogin(gatewayid uint64, serialnum uint32) *NsqLoginPacket {
	var result uint8 = Offline
	online := NewConns().Check(gatewayid)
	var indb bool = false
	if !online {
		indb = GetGatewayHub().Check(gatewayid)
	}
	if online {
		result = Online
	} else if !indb {
		result = UnAuth
	}

	return &NsqLoginPacket{
		GatewayID:    gatewayid,
		SerialNumber: serialnum,
		Result:       result,
	}
}
