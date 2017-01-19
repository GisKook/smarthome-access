package protocol

import (
	"github.com/giskook/smarthome-access/base"
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

type Notify_Online_Status_Packet struct {
	GatewayID uint64
	SerialID  uint32
	DeviceID  uint64
	Status    uint8
}

func (p *Notify_Online_Status_Packet) Serialize() []byte {
	para := []*Report.Command_Param{
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT64,
			Npara: p.DeviceID,
		},
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(p.Status),
		},
	}

	command := &Report.Command{
		Type:  Report.Command_CMT_REP_NOTIFY_ONOFF,
		Paras: para,
	}

	notify_onoff := &Report.ControlReport{
		Tid:          p.GatewayID,
		SerialNumber: p.SerialID,
		Command:      command,
	}

	data, _ := proto.Marshal(notify_onoff)

	return data
}

func Parse_Notify_Online_Status(buffer []byte, id uint64) *Notify_Online_Status_Packet {
	reader := ParseHeader(buffer)
	serialid := base.ReadDWord(reader)
	deviceid := base.ReadQuaWord(reader)
	status, _ := reader.ReadByte()

	return &Notify_Online_Status_Packet{
		GatewayID: id,
		SerialID:  serialid,
		DeviceID:  deviceid,
		Status:    status,
	}
}
