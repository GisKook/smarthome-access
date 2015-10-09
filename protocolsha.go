package sha

import (
	"github.com/giskook/gotcp"
	"log"
)

var (
	Illegal  uint16 = 0
	HalfPack uint16 = 255

	Login                 uint16 = 1
	HeartBeat             uint16 = 2
	SendDeviceList        uint16 = 3
	OperateFeedback       uint16 = 4
	AddDelDevice          uint16 = 5
	Warn                  uint16 = 6
	SetDevicenameFeedback uint16 = 7
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
	case AddDelDevice:
		return this.Packet.(*AddDelDevicePacket).Serialize()
	case SetDevicenameFeedback:
		return this.Packet.(*FeedbackSetDevicenamePacket).Serialize()
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
			log.Printf("recv box cmd %d \n", cmdid)

			pkgbyte := make([]byte, pkglen)
			buffer.Read(pkgbyte)
			switch cmdid {
			case Login:
				pkg := ParseLogin(pkgbyte, smconn)
				return NewShaPacket(Login, pkg), nil
			case HeartBeat:
				pkg := ParseHeart(pkgbyte)
				return NewShaPacket(HeartBeat, pkg), nil
			case SendDeviceList:
				pkg := ParseDeviceList(pkgbyte, smconn)
				return NewShaPacket(SendDeviceList, pkg), nil
			case OperateFeedback:
				pkg := ParseFeedback(pkgbyte)
				return NewShaPacket(OperateFeedback, pkg), nil
			case Warn:
				pkg := ParseWarn(pkgbyte)
				return NewShaPacket(Warn, pkg), nil
			case AddDelDevice:
				pkg := ParseAddDelDevice(pkgbyte)
				return NewShaPacket(AddDelDevice, pkg), nil
			case SetDevicenameFeedback:
				pkg := ParseFeedbackSetDevicename(pkgbyte)
				return NewShaPacket(SetDevicenameFeedback, pkg), nil

			case Illegal:
			case HalfPack:
			}
		}
	}

}
