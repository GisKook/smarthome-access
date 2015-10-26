package sha

import (
	"encoding/binary"
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

var (
	DeviceNotExist uint8 = 2
)

type NsqDelDevicePacket struct {
	GatewayID    uint64
	SerialNumber uint32
	DeviceID     uint64
	Exist        bool
}

func (p *NsqDelDevicePacket) Serialize() []byte {
	if p.Exist {
		buf := []byte{0xCE, 0x00, 0x18, 0x80, 0x0A}
		gatewayid := make([]byte, 8)
		binary.BigEndian.PutUint64(gatewayid, p.GatewayID)
		buf = append(buf, gatewayid[2:]...)
		serialnum := make([]byte, 4)
		binary.BigEndian.PutUint32(serialnum, p.SerialNumber)
		buf = append(buf, serialnum...)
		buf = append(buf, 0x06)
		deviceid := make([]byte, 8)
		binary.BigEndian.PutUint64(deviceid, p.DeviceID)
		buf = append(buf, deviceid[2:]...)
		buf = append(buf, CheckSum(buf, uint16(len(buf))))
		buf = append(buf, 0xCE)

		return buf
	} else {
		para := []*Report.Command_Param{
			&Report.Command_Param{
				Type:  Report.Command_Param_UINT8,
				Npara: uint64(DeviceNotExist),
			},
		}

		command := &Report.Command{
			Type:  Report.Command_CMT_REPDELDEVICE,
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
}

func ParseNsqDelDevice(gatewayid uint64, serialnum uint32, command *Report.Command) (*NsqDelDevicePacket, bool) {
	commandparam := command.GetParas()
	exist := NewGatewayHub().CheckDevice(gatewayid, commandparam[0].Npara)

	return &NsqDelDevicePacket{
		GatewayID:    gatewayid,
		SerialNumber: serialnum,
		DeviceID:     commandparam[0].Npara,
		Exist:        exist,
	}, exist
}
