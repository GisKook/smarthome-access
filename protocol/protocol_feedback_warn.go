package protocol

import (
	"github.com/giskook/smarthome-access/base"
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

type Feedback_Warn_Packet struct {
	GatewayID uint64
	SerialNum uint32
	DeviceID  uint64
	Endpoint  uint8
	CommandID uint8
	Status    uint8
}

func (p *Feedback_Warn_Packet) Serialize() []byte {
	para := []*Report.Command_Param{
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(GATEWAY_ON_LINE),
		},
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(p.Status),
		},
	}
	command := &Report.Command{
		Type:  Report.Command_CMT_REP_FEEDBACK_WARN,
		Paras: para,
	}

	feedback_warn_pkg := &Report.ControlReport{
		Tid:          p.GatewayID,
		SerialNumber: p.SerialNum,
		Command:      command,
	}

	data, _ := proto.Marshal(feedback_warn_pkg)

	return data
}

func Parse_Feedback_Warn(buffer []byte, id uint64) *Feedback_Warn_Packet {

	reader := ParseHeader(buffer)
	serialnum := base.ReadDWord(reader)
	deviceid := base.ReadQuaWord(reader)
	endpoint, _ := reader.ReadByte()
	commandid, _ := reader.ReadByte()
	status, _ := reader.ReadByte()

	return &Feedback_Warn_Packet{
		GatewayID: id,
		SerialNum: serialnum,
		DeviceID:  deviceid,
		Endpoint:  endpoint,
		CommandID: commandid,
		Status:    status,
	}
}
