package sha

import (
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
	"log"
)

var (
	Offline     uint8 = 0
	Online      uint8 = 1
	UnAuth      uint8 = 2
	PasswdError uint8 = 3
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
					Npara: uint64(p.Result),
				},
			},
		},
	}

	repdata, err := proto.Marshal(rep)
	if err != nil {
		log.Println("marshaling error", err)
	}

	return repdata
}

func ParseNsqLogin(gatewayid uint64, serialnum uint32, command *Report.Command) *NsqLoginPacket {
	var result uint8 = Offline
	cmdparam := command.GetParas()
	if GetUserPasswdHub().Check(gatewayid, cmdparam[0].Strpara) {
		online := NewConns().Check(gatewayid)
		var indb bool = false
		if !online {
			indb = NewGatewayHub().Check(gatewayid)
		}
		if online {
			result = Online
		} else if !indb {
			result = UnAuth
		}
	} else {
		result = PasswdError
	}

	return &NsqLoginPacket{
		GatewayID:    gatewayid,
		SerialNumber: serialnum,
		Result:       result,
	}
}
