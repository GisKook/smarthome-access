package protocol

import (
	"github.com/giskook/smarthome-access/base"
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

type Notify_Level_Packet struct {
	GatewayID uint64
	SerialID  uint32
	DeviceID  uint64
	Endpoint  uint8
	Value     uint8
}

func (p *Notify_Level_Packet) Serialize() []byte {
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
			Npara: uint64(p.Value),
		},
	}

	command := &Report.Command{
		Type:  Report.Command_CMT_REP_NOTIFY_LEVEL,
		Paras: para,
	}

	notify_level := &Report.ControlReport{
		Tid:          p.GatewayID,
		SerialNumber: p.SerialID,
		Command:      command,
	}

	data, _ := proto.Marshal(notify_level)

	return data
}

func Parse_Notify_Level(buffer []byte, id uint64) *Notify_Level_Packet {
	reader := ParseHeader(buffer)
	serialid := base.ReadDWord(reader)
	deviceid := base.ReadQuaWord(reader)
	endpoint, _ := reader.ReadByte()
	value, _ := reader.ReadByte()

	return &Notify_Level_Packet{
		GatewayID: id,
		SerialID:  serialid,
		DeviceID:  deviceid,
		Endpoint:  endpoint,
		Value:     value,
	}
}
