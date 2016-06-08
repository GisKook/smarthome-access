package sha

import (
	"github.com/giskook/gotcp"
	"github.com/giskook/smarthome-access/protocol"
)

type ShaPacket struct {
	Type   uint16
	Packet gotcp.Packet
}

func (this *ShaPacket) Serialize() []byte {
	switch this.Type {
	case protocol.Login:
		return this.Packet.(*protocol.LoginPacket).Serialize()
	case protocol.HeartBeat:
		return this.Packet.(*protocol.HeartPacket).Serialize()
	}

	return nil
}

func NewShaPacket(Type uint16, Packet gotcp.Packet) *ShaPacket {
	return &ShaPacket{
		Type:   Type,
		Packet: Packet,
	}
}

type ShaProtocol struct {
}

func (this *ShaProtocol) ReadPacket(c *gotcp.Conn) (gotcp.Packet, error) {
	smconn := c.GetExtraData().(*Conn)
	smconn.UpdateReadflag()

	buffer := smconn.GetBuffer()

	conn := c.GetRawConn()
	for {
		data := make([]byte, 2048)
		readLengh, err := conn.Read(data)

		if err != nil {
			return nil, err
		}

		if readLengh == 0 {
			return nil, gotcp.ErrConnClosing
		} else {
			buffer.Write(data[0:readLengh])
			cmdid, pkglen := protocol.CheckProtocol(buffer)

			pkgbyte := make([]byte, pkglen)
			buffer.Read(pkgbyte)
			switch cmdid {
			case protocol.Login:
				pkg := protocol.ParseLogin(pkgbyte)
				return NewShaPacket(protocol.Login, pkg), nil
			case protocol.HeartBeat:
				pkg := protocol.ParseHeart(pkgbyte)
				return NewShaPacket(protocol.HeartBeat, pkg), nil

			case protocol.Illegal:
			case protocol.HalfPack:
			}
		}
	}

}
