package protocol

import (
	"github.com/giskook/smarthome-access/base"
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

type Notification_Packet struct {
	GatewayID        uint64
	DeviceID         uint64
	NotificationTime uint64
	Endpoint         uint8
	DeviceTypeID     uint16
	Zonetype         uint16
	ZoneStatus       uint16
}

func (p *Notification_Packet) Serialize() []byte {
	para := []*Report.Command_Param{
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT64,
			Npara: p.DeviceID,
		},
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(p.Endpoint),
		},
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT64,
			Npara: p.NotificationTime,
		},
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT16,
			Npara: uint64(p.Zonetype),
		},
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT16,
			Npara: uint64(p.ZoneStatus),
		},
	}

	command := &Report.Command{
		Type:  Report.Command_CMT_REP_NOTIFICATION,
		Paras: para,
	}

	login := &Report.ControlReport{
		Tid:          p.GatewayID,
		SerialNumber: 0,
		Command:      command,
	}

	data, _ := proto.Marshal(login)

	return data
}

func Parse_Notification(buffer []byte, id uint64) *Notification_Packet {
	reader := ParseHeader(buffer)
	deviceid := base.ReadQuaWord(reader)
	notification_time := base.ReadQuaWord(reader) + 8*60*60
	endpoint, _ := reader.ReadByte()
	devicetypeid := base.ReadWord(reader)
	zonetype := base.ReadWord(reader)
	zonestatus := base.ReadWord(reader)

	return &Notification_Packet{
		GatewayID:        id,
		DeviceID:         deviceid,
		NotificationTime: notification_time,
		Endpoint:         endpoint,
		DeviceTypeID:     devicetypeid,
		Zonetype:         zonetype,
		ZoneStatus:       zonestatus,
	}
}
