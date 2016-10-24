package protocol

import (
	"github.com/giskook/smarthome-access/base"
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

const READON uint8 = 1
const READOFF uint8 = 0

type Feedback_Read_OnOff_Status_Packet struct {
	GatewayID uint64
	SerialNum uint32
	DeviceID  uint64
	Endpoint  uint8
	Action    uint8
	Result    uint8
}

func (p *Feedback_Read_OnOff_Status_Packet) Serialize() []byte {
	para := []*Report.Command_Param{
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(GATEWAY_ON_LINE),
		},
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(p.Action),
		},
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(p.Result),
		},
	}

	command := &Report.Command{
		Type:  Report.Command_CMT_REP_NOTIFY_ONOFF,
		Paras: para,
	}

	feedback_read_onoff_status_pkg := &Report.ControlReport{
		Tid:          p.GatewayID,
		SerialNumber: p.SerialNum,
		Command:      command,
	}

	data, _ := proto.Marshal(feedback_read_onoff_status_pkg)

	return data
}

func Parse_Feedback_Read_OnOff_Status(buffer []byte, id uint64) *Feedback_Read_OnOff_Status_Packet {
	reader := ParseHeader(buffer)
	serialnum := base.ReadDWord(reader)
	deviceid := base.ReadQuaWord(reader)
	endpoint, _ := reader.ReadByte()
	action, _ := reader.ReadByte()
	result, _ := reader.ReadByte()

	return &Feedback_Read_OnOff_Status_Packet{
		GatewayID: id,
		SerialNum: serialnum,
		DeviceID:  deviceid,
		Endpoint:  endpoint,
		Action:    action,
		Result:    result,
	}
}
