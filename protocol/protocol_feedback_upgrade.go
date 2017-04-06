package protocol

import (
	"github.com/giskook/smarthome-access/base"
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

type Feedback_Upgrade_Packet struct {
	GatewayID uint64
	SerialNum uint32
}

func (p *Feedback_Upgrade_Packet) Serialize() []byte {
	para := []*Report.Command_Param{
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(GATEWAY_ON_LINE),
		},
	}

	command := &Report.Command{
		Type:  Report.Command_CMT_REP_FEEDBACK_UPGRADE,
		Paras: para,
	}

	feedback_onoff_pkg := &Report.ControlReport{
		Tid:          p.GatewayID,
		SerialNumber: p.SerialNum,
		Command:      command,
	}

	data, _ := proto.Marshal(feedback_onoff_pkg)

	return data
}

func Parse_Feedback_Upgrade(buffer []byte, id uint64) *Feedback_Upgrade_Packet {
	reader := ParseHeader(buffer)
	serialnum := base.ReadDWord(reader)

	return &Feedback_Upgrade_Packet{
		GatewayID: id,
		SerialNum: serialnum,
	}
}
