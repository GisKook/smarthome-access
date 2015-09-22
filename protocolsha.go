package sha

import (
	"github.com/giskook/gotcp"
	"log"
)

var (
	Illegal  uint16 = 0
	HalfPack uint16 = 255

	Login           uint16 = 1
	HeartBeat       uint16 = 2
	SendDeviceList  uint16 = 3
	OperateFeedback uint16 = 4
	AddDevice       uint16 = 5
	DelDevice       uint16 = 5
	Warn            uint16 = 7
)

type ShaPacket struct {
	Type   uint16
	Packet gotcp.Packet
}

func (this *ShaPacket) Serialize() []byte {
	switch this.Type {
	case Login:
		return this.Packet.(*LoginPacket).Serialize()
	case HeartBeat:
		return this.Packet.(*HeartPacket).Serialize()
	case SendDeviceList:
		return this.Packet.(*DeviceListPacket).Serialize()
	case OperateFeedback:
		return this.Packet.(*FeedbackPacket).Serialize()
	case Warn:
		return this.Packet.(*WarnPacket).Serialize()
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
			cmdid, pkglen := CheckProtocol(buffer)
			log.Println(cmdid)

			switch cmdid {
			case Login:
				pkgbyte := make([]byte, pkglen)
				buffer.Read(pkgbyte)
				pkg := ParseLogin(pkgbyte, smconn)
				return NewShaPacket(Login, pkg), nil
			case HeartBeat:
				pkgbyte := make([]byte, pkglen)
				buffer.Read(pkgbyte)
				pkg := ParseHeart(pkgbyte)
				return NewShaPacket(HeartBeat, pkg), nil
			case SendDeviceList:
				pkgbyte := make([]byte, pkglen)
				buffer.Read(pkgbyte)
				pkg := ParseDeviceList(pkgbyte, smconn)
				return NewShaPacket(SendDeviceList, pkg), nil
			case OperateFeedback:
				pkgbyte := make([]byte, pkglen)
				buffer.Read(pkgbyte)
				pkg := ParseFeedback(pkgbyte)
				return NewShaPacket(OperateFeedback, pkg), nil
			case Warn:
				pkgbyte := make([]byte, pkglen)
				buffer.Read(pkgbyte)
				pkg := ParseWarn(pkgbyte)
				return NewShaPacket(Warn, pkg), nil
			case Illegal:
			case HalfPack:
			}
		}
	}

}
