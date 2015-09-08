package sha

import (
	"github.com/giskook/gotcp"
)

var (
	Login     uint8 = 0
	HeartBeat uint8 = 1
)

type ShaPacket struct {
	Type   uint8
	Packet gotcp.Packet
}

func (this *ShaPacket) Serialize() []byte {
	switch this.Type {
	case Login:
		return this.Packet.(*LoginPacket).Serialize()
	case HeartBeat:
		return this.Packet.(*HeartPacket).Serialize()
	}

	return nil
}

func NewShaPacket(Type uint8, Packet gotcp.Packet) *ShaPacket {
	return &ShaPacket{
		Type:   Type,
		Packet: Packet,
	}
}

type ShaProtocol struct {
}

func (this *ShaProtocol) ReadPacket(c *gotcp.Conn) (gotcp.Packet, error) {
	conn := c.GetRawConn()

	for {
		data := make([]byte, 1024)
		readLengh, err := conn.Read(data)

		if err != nil {
			return nil, err
		}

		var uid []byte
		uid = append(uid, 0x01)
		uid = append(uid, 0x02)
		if readLengh == 0 {
			return nil, gotcp.ErrConnClosing
		} else {
			if data[1] == 0xAA {
				return NewShaPacket(Login, NewLoginPakcet(uid, 1, 1)), nil
			} else {
				return NewShaPacket(HeartBeat, NewHeartPacket(uid)), nil
			}
		}
	}

}
