package protocol

import (
	"github.com/giskook/smarthome-access/base"
)

type HeartPacket struct {
	ID uint64
}

func (p *HeartPacket) Serialize() []byte {
	var buf []byte
	buf = append(buf, 0xCE)
	buf = append(buf, 0x00)
	buf = append(buf, 0x0D)
	buf = append(buf, 0x80)
	buf = append(buf, 0x02)
	buf = append(buf, base.WriteMac(p.ID)...)
	buf = append(buf, CheckSum(buf, uint16(len(buf))))
	buf = append(buf, 0xCE)

	return buf
}

func ParseHeart(buffer []byte, gatewayid uint64) *HeartPacket {
	return &HeartPacket{
		ID: gatewayid,
	}
}
