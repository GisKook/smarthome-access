package sha

import (
	"encoding/binary"
)

type NsqDeviceListPacket struct {
	GatewayID    uint64
	SerialNumber uint32
}

func (p *NsqDeviceListPacket) Serialize() []byte {
	var buf []byte
	buf = append(buf, 0xCE)
	buf = append(buf, 0x00)
	buf = append(buf, 0x0B)
	buf = append(buf, 0x80)
	buf = append(buf, 0x03)
	gatewayid := make([]byte, 8)
	binary.BigEndian.PutUint64(gatewayid, p.GatewayID)
	buf = append(buf, gatewayid[2:]...)
	buf = append(buf, CheckSum(buf, uint16(len(buf))))
	buf = append(buf, 0xCE)

	return buf
}

func ParseNsqDeviceList(gatewayid uint64, serialnum uint32) *NsqDeviceListPacket {
	return &NsqDeviceListPacket{
		GatewayID:    gatewayid,
		SerialNumber: serialnum,
	}
}
