package protocol

import (
	"github.com/giskook/smarthome-access/base"
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

type Feedback_SetName_Packet struct {
	GatewayID  uint64
	DeviceID   uint64
	SerialNum  uint32
	Result     uint8
	DeviceName string
}

func (p *Feedback_SetName_Packet) Serialize() []byte {
	para := []*Report.Command_Param{
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(GATEWAY_ON_LINE),
		},
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(p.Result),
		},
	}

	command := &Report.Command{
		Type:  Report.Command_CMT_REP_SETNAME,
		Paras: para,
	}

	feedback_setname_pkg := &Report.ControlReport{
		Tid:          p.GatewayID,
		SerialNumber: p.SerialNum,
		Command:      command,
	}

	data, _ := proto.Marshal(feedback_setname_pkg)

	return data
}

func Parse_Feedback_SetName(buffer []byte, id uint64) *Feedback_SetName_Packet {
	reader := ParseHeader(buffer)
	serialnum := base.ReadDWord(reader)
	deviceid := base.ReadQuaWord(reader)
	result, _ := reader.ReadByte()
	devicename_length, _ := reader.ReadByte()
	devicename := base.ReadString(reader, devicename_length)

	return &Feedback_SetName_Packet{
		GatewayID:  id,
		DeviceID:   deviceid,
		SerialNum:  serialnum,
		Result:     result,
		DeviceName: devicename,
	}
}
