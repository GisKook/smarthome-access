package protocol

import (
	"github.com/giskook/smarthome-access/base"
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

type Feedback_Query_Attr_Packet struct {
	GatewayID          uint64
	SerialNum          uint32
	DeviceID           uint64
	Endpoint           uint8
	ShortAddr          uint16
	ProfileID          uint16
	ZclVersion         uint8
	ApplicationVersion uint8
	StackVersion       uint8
	HwVersion          uint8
	ManufacturerName   string
	ModelIdentifier    string
	DateCode           string
	PowerSource        uint8
}

func (p *Feedback_Query_Attr_Packet) Serialize() []byte {
	para := []*Report.Command_Param{
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT16,
			Npara: uint64(p.ShortAddr),
		},
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT16,
			Npara: uint64(p.ProfileID),
		},
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(p.ZclVersion),
		},
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(p.ApplicationVersion),
		},
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(p.StackVersion),
		},
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(p.HwVersion),
		},
		&Report.Command_Param{
			Type:    Report.Command_Param_STRING,
			Strpara: p.ManufacturerName,
		},
		&Report.Command_Param{
			Type:    Report.Command_Param_STRING,
			Strpara: p.ModelIdentifier,
		},
		&Report.Command_Param{
			Type:    Report.Command_Param_STRING,
			Strpara: p.DateCode,
		},
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(p.PowerSource),
		},
	}

	command := &Report.Command{
		Type:  Report.Command_CMT_REP_DEVICE_ATTR,
		Paras: para,
	}

	login := &Report.ControlReport{
		Tid:          p.GatewayID,
		SerialNumber: p.SerialNum,
		Command:      command,
	}

	data, _ := proto.Marshal(login)

	return data
}

func Parse_Feedback_Query_Attr(buffer []byte, id uint64) *Feedback_Query_Attr_Packet {
	reader := ParseHeader(buffer)
	serialnum := base.ReadDWord(reader)
	deviceid := base.ReadQuaWord(reader)
	endpoint, _ := reader.ReadByte()
	shortaddr := base.ReadWord(reader)
	profileid := base.ReadWord(reader)
	zclversion, _ := reader.ReadByte()
	applicationversion, _ := reader.ReadByte()
	stackversion, _ := reader.ReadByte()
	hwversion, _ := reader.ReadByte()
	manufacturernamelen := reader.ReadByte()
	manufacturername := base.ReadString(reader, manufacturernamelen)
	modelidentifierlen, _ := reader.ReadByte()
	modelidentifier := base.ReadString(reader, modelidentifierlen)
	datecodelen, _ := reader.ReadByte()
	datecode := base.ReadWord(reader, datecodelen)
	powersouce, _ := reader.ReadByte()

	return &Feedback_Query_Attr_Packet{
		GatewayID:          id,
		SerialNum:          serialnum,
		DeviceID:           deviceid,
		Endpoint:           endpoint,
		ShortAddr:          shortaddr,
		ProfileID:          profileid,
		ZclVersion:         zclversion,
		ApplicationVersion: applicationversion,
		StackVersion:       stackversion,
		HwVersion:          hwversion,
		ManufacturerName:   manufacturername,
		ModelIdentifier:    modelidentifier,
		DateCode:           datecode,
		PowerSource:        powersouce,
	}
}
