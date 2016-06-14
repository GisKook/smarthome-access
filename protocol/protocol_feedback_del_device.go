package protocol

import (
	"github.com/giskook/smarthome-access/base"
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

type Feedback_Del_Device_Packet struct {
	GatewayID uint64
	SerialNum uint32
	DeviceID  uint64
	Result    uint8
}

func (p *Feedback_Del_Device_Packet) Serialize() []byte {
	para := []*Report.Command_Param{
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(p.Result),
		},
	}

	command := &Report.Command{
		Type:  Report.Command_CMT_REP_DEL_DEVICE,
		Paras: para,
	}

	feedback_del_device_pkg := &Report.ControlReport{
		Tid:          p.GatewayID,
		SerialNumber: p.SerialNum,
		Command:      command,
	}

	data, _ := proto.Marshal(feedback_del_device_pkg)

	return data
}

func Parse_Feedback_Del_Device(buffer []byte, id uint64) *Feedback_Del_Device_Packet {
	reader := ParseHeader(buffer)
	serialnum := base.ReadDWord(reader)
	deviceid := base.ReadQuaWord(reader)
	result, _ := reader.ReadByte()

	return &Feedback_Del_Device_Packet{
		GatewayID: id,
		SerialNum: serialnum,
		DeviceID:  deviceid,
		Result:    result,
	}
}
