package sha

import (
	"bytes"
	"github.com/giskook/gotcp"
)

var (
	Login           uint8 = 0
	HeartBeat       uint8 = 1
	SendDeviceList  uint8 = 2
	OperateFeedback uint8 = 3
	AddDevice       uint8 = 4
	DelDevice       uint8 = 5
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
	smconn := c.GetExtraData()
	element, _ := smconn.(Conn)
	buffer := element.GetBuffer()

	for {
		var data [2048]byte
		readLengh, err := conn.Read(data)

		if err != nil {
			return nil, err
		}

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
