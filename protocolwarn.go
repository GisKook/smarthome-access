package sha

import (
	"encoding/binary"
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

type WarnPacket struct {
	GatewayID  uint64
	Time       uint64
	DeviceType uint8
	DeviceID   uint64
	WarnType   uint8
}

func (p *WarnPacket) Serialize() []byte {
	para := []*Report.Command_Param{
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT64,
			Npara: p.Time,
		},
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(p.DeviceType),
		},
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT64,
			Npara: p.DeviceID,
		},
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(p.WarnType),
		},
	}

	command := &Report.Command{
		Type:  Report.Command_CMT_REPWARNUP,
		Paras: para,
	}

	warn := &Report.ControlReport{
		Tid:          p.GatewayID,
		SerialNumber: 0,
		Command:      command,
	}

	data, _ := proto.Marshal(warn)

	return data
}

func ParseWarn(buffer []byte) *WarnPacket {
	gatewayid, reader := GetGatewayID(buffer)
	devicetype, _ := reader.ReadByte()
	deviceidlen, _ := reader.ReadByte()
	deviceid_byte := make([]byte, deviceidlen)
	reader.Read(deviceid_byte)
	tmp := []byte{0, 0}
	deviceid_byte = append(tmp, deviceid_byte...)
	deviceid := binary.BigEndian.Uint64(deviceid_byte)
	warntime_byte := make([]byte, 8)
	reader.Read(warntime_byte)
	warntime := binary.BigEndian.Uint64(warntime_byte)
	warntype, _ := reader.ReadByte()

	return &WarnPacket{
		GatewayID:  gatewayid,
		Time:       warntime,
		DeviceType: devicetype,
		DeviceID:   deviceid,
		WarnType:   warntype,
	}
}
