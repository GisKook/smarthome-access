package sha

import (
	"encoding/binary"
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

type FeedbackSetDevicenamePacket struct {
	GatewayID    uint64
	SerialNumber uint32
	Result       uint8
	DeviceID     uint64
	DeviceName   string
}

func (p *FeedbackSetDevicenamePacket) Serialize() []byte {
	para := []*Report.Command_Param{
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(p.Result),
		},
	}

	command := &Report.Command{
		Type:  Report.Command_CMT_REPSETDEVICENAME,
		Paras: para,
	}
	feedback := &Report.ControlReport{
		Tid:          p.GatewayID,
		SerialNumber: p.SerialNumber,
		Command:      command,
	}

	data, _ := proto.Marshal(feedback)

	return data
}

func ParseFeedbackSetDevicename(buffer []byte) *FeedbackSetDevicenamePacket {
	gatewayid, reader := GetGatewayID(buffer)
	serialnumber_byte := make([]byte, 4)
	reader.Read(serialnumber_byte)
	serialnumber := binary.BigEndian.Uint32(serialnumber_byte)

	result, _ := reader.ReadByte()
	deviceidlen, _ := reader.ReadByte()
	deviceid_byte := make([]byte, deviceidlen)
	reader.Read(deviceid_byte)
	did := []byte{0, 0}
	did = append(did, deviceid_byte...)
	deviceid := binary.BigEndian.Uint64(did)
	devicenamelen, _ := reader.ReadByte()
	devicename := make([]byte, devicenamelen)
	if result == 1 {
		NewGatewayHub().Setname(gatewayid, deviceid, string(devicename))
	}

	return &FeedbackSetDevicenamePacket{
		GatewayID:    gatewayid,
		SerialNumber: serialnumber,
		Result:       result,
		DeviceID:     deviceid,
		DeviceName:   string(devicename),
	}
}
