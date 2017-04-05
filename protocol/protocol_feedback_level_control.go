package protocol

import (
	"github.com/giskook/smarthome-access/base"
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

type Feedback_Level_Control_Packet struct {
	GatewayID uint64
	SerialNum uint32
	DeviceID  uint64
	Endpoint  uint8
	Level     uint8
}

func (p *Feedback_Level_Control_Packet) Serialize() []byte {
	para := []*Report.Command_Param{
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(GATEWAY_ON_LINE),
		},
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(p.Level),
		},
	}

	command := &Report.Command{
		Type:  Report.Command_CMT_REP_ONOFF,
		Paras: para,
	}

	feedback_level_control_pkg := &Report.ControlReport{
		Tid:          p.GatewayID,
		SerialNumber: p.SerialNum,
		Command:      command,
	}

	data, _ := proto.Marshal(feedback_level_control_pkg)

	return data
}

func Parse_Feedback_Level_Control(buffer []byte, id uint64) *Feedback_Level_Control_Packet {
	reader := ParseHeader(buffer)
	serialnum := base.ReadDWord(reader)
	deviceid := base.ReadQuaWord(reader)
	endpoint, _ := reader.ReadByte()
	level, _ := reader.ReadByte()

	return &Feedback_Level_Control_Packet{
		GatewayID: id,
		SerialNum: serialnum,
		DeviceID:  deviceid,
		Endpoint:  endpoint,
		Level:     level,
	}
}
