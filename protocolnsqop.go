package sha

import (
	"encoding/binary"
	"github.com/giskook/smarthome-access/pb"
)

type NsqOpPacket struct {
	GatewayID    uint64
	SerialNumber uint32
	DeviceID     uint64
	Op           uint8
}

func (p *NsqOpPacket) Serialize() []byte {
	buf := []byte{0xCE, 0x00, 0x17, 0x08, 0x04}
	gatewayid := make([]byte, 8)
	binary.BigEndian.PutUint64(gatewayid, p.GatewayID)
	buf = append(buf, gatewayid[2:]...)
	buf = append(buf, CheckSum(buf, uint16(len(buf))))
	buf = append(buf, 0xCE)

	return buf

}

func ParseNsqOp(gatewayid uint64, serialnum uint32, command *Report.Command) *NsqOpPacket {
	commandparam := command.GetParas()

	return &NsqOpPacket{
		GatewayID:    gatewayid,
		SerialNumber: serialnum,
		DeviceID:     commandparam[0].Npara,
		Op:           uint8(commandparam[1].Npara),
	}
}
