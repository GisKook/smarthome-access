package protocol

import (
	"github.com/giskook/smarthome-access/base"
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

type Notify_OnOff_Packet struct {
	GatewayID uint64
	SerialID  uint32
	DeviceID  uint64
	Endpoint  uint8
	Status    uint8
}

func (p *Notify_OnOff_Packet) Serialize() []byte {
	para := []*Report.Command_Param{
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT32,
			Npara: uint64(p.SerialID),
		},
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT64,
			Npara: p.DeviceID,
		},
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(p.Endpoint),
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

func Parse_Notify_OnOff(buffer []byte, id uint64) *Notify_OnOff_Packet {
	reader := ParseHeader(buffer)
	serialid := base.ReadDWord(reader)
	deviceid := base.ReadQuaWord(reader)
	endpoint, _ := reader.ReadByte()
	status, _ := reader.ReadByte()

	return &Notify_OnOff_Packet{
		GatewayID: id,
		SerialID:  serialid,
		DeviceID:  deviceid,
		Endpoint:  endpoint,
		Status:    status,
	}
}
