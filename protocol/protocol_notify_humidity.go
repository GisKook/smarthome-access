package protocol

import (
	"github.com/giskook/smarthome-access/base"
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

type Notify_Humidity_Packet struct {
	GatewayID uint64
	SerialID  uint32
	TimeStamp uint32
	DeviceID  uint64
	Endpoint  uint8
	Humidity  uint16
}

func (p *Notify_Humidity_Packet) Serialize() []byte {
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
			Type:  Report.Command_Param_UINT32,
			Npara: uint64(p.TimeStamp),
		},
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT16,
			Npara: uint64(p.Humidity),
		},
	}

	command := &Report.Command{
		Type:  Report.Command_CMT_REP_NOTIFY_HUMIDITY,
		Paras: para,
	}

	notify_temperature := &Report.ControlReport{
		Tid:          p.GatewayID,
		SerialNumber: p.SerialID,
		Command:      command,
	}

	data, _ := proto.Marshal(notify_temperature)

	return data
}

func Parse_Notify_Humidity(buffer []byte, id uint64) *Notify_Humidity_Packet {
	reader := ParseHeader(buffer)
	serialid := base.ReadDWord(reader)
	time_stamp := base.ReadDWord(reader)
	deviceid := base.ReadQuaWord(reader)
	endpoint, _ := reader.ReadByte()
	humidity := base.ReadWord(reader)

	return &Notify_Humidity_Packet{
		GatewayID: id,
		SerialID:  serialid,
		TimeStamp: time_stamp,
		DeviceID:  deviceid,
		Endpoint:  endpoint,
		Humidity:  humidity,
	}
}
