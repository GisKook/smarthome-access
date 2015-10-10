package sha

import (
	"encoding/binary"
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

type AddDelDevicePacket struct {
	GatewayID  uint64
	Action     uint8
	DeviceType uint8
	DeviceID   uint64
	Company    uint16
	Status     uint8
}

func (p *AddDelDevicePacket) Serialize() []byte {
	para := []*Report.Command_Param{
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(p.Action),
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
			Type:  Report.Command_Param_UINT16,
			Npara: uint64(p.Company),
		},

		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(p.Status),
		},
	}

	command := &Report.Command{
		Type:  Report.Command_CMT_REPADDDELDEVICE,
		Paras: para,
	}

	adddeldevice := &Report.ControlReport{
		Tid:          p.GatewayID,
		SerialNumber: 0,
		Command:      command,
	}

	data, _ := proto.Marshal(adddeldevice)

	return data
}

func ParseAddDelDevice(buffer []byte) *AddDelDevicePacket {
	gatewayid, reader := GetGatewayID(buffer)
	action, _ := reader.ReadByte()
	devicetype, _ := reader.ReadByte()
	deviceidlen, _ := reader.ReadByte()
	deviceid_byte := make([]byte, deviceidlen)
	reader.Read(deviceid_byte)
	tmp := []byte{0, 0}
	deviceid_byte = append(tmp, deviceid_byte...)
	deviedid := binary.BigEndian.Uint64(deviceid_byte)
	company_byte := make([]byte, 2)
	company := binary.BigEndian.Uint16(company_byte)
	status, _ := reader.ReadByte()

	return &AddDelDevicePacket{
		GatewayID:  gatewayid,
		DeviceType: devicetype,
		DeviceID:   deviedid,
		Action:     action,
		Company:    company,
		Status:     status,
	}
}
