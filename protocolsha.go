package sha

import (
	"github.com/giskook/gotcp"
	"github.com/giskook/smarthome-access/protocol"
	"log"
	"sync"
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
	case protocol.Notification:
		return this.Packet.(*protocol.Notification_Packet).Serialize()
	case protocol.Add_Del_Device:
		return this.Packet.(*protocol.Add_Del_Device_Packet).Serialize()
	case protocol.Feedback_SetName:
		return this.Packet.(*protocol.Feedback_SetName_Packet).Serialize()
	case protocol.Feedback_Del_Device:
		return this.Packet.(*protocol.Feedback_Del_Device_Packet).Serialize()
	case protocol.Feedback_Query_Attr:
		return this.Packet.(*protocol.Feedback_Query_Attr_Packet).Serialize()
	case protocol.Feedback_Depolyment:
		return this.Packet.(*protocol.Feedback_Deployment_Packet).Serialize()
	case protocol.Feedback_OnOff:
		return this.Packet.(*protocol.Feedback_OnOff_Packet).Serialize()
	case protocol.Feedback_Level_Control:
		return this.Packet.(*protocol.Feedback_Level_Control_Packet).Serialize()
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
	var once sync.Once
	once.Do(smconn.UpdateReadflag)

	buffer := smconn.GetBuffer()

	conn := c.GetRawConn()
	for {
		if smconn.ReadMore {
			data := make([]byte, 2048)
			readLengh, err := conn.Read(data)
			log.Printf("<IN>    %x\n", data[0:readLengh])
			if err != nil {
				return nil, err
			}

			if readLengh == 0 {
				return nil, gotcp.ErrConnClosing
			}
			buffer.Write(data[0:readLengh])
		}

		cmdid, pkglen := protocol.CheckProtocol(buffer)

		pkgbyte := make([]byte, pkglen)
		buffer.Read(pkgbyte)
		switch cmdid {
		case protocol.Login:
			pkg := protocol.ParseLogin(pkgbyte)
			smconn.ReadMore = false
			return NewShaPacket(protocol.Login, pkg), nil
		case protocol.HeartBeat:
			pkg := protocol.ParseHeart(pkgbyte, smconn.ID)
			smconn.ReadMore = false
			return NewShaPacket(protocol.HeartBeat, pkg), nil
		case protocol.Add_Del_Device:
			pkg := protocol.Parse_Add_Del_Device(pkgbyte, smconn.ID)
			smconn.ReadMore = false
			return NewShaPacket(protocol.Add_Del_Device, pkg), nil
		case protocol.Notification:
			pkg := protocol.Parse_Notification(pkgbyte, smconn.ID)
			smconn.ReadMore = false
			return NewShaPacket(protocol.Notification, pkg), nil
		case protocol.Feedback_SetName:
			pkg := protocol.Parse_Feedback_SetName(pkgbyte, smconn.ID)
			smconn.ReadMore = false
			return NewShaPacket(protocol.Feedback_SetName, pkg), nil
		case protocol.Feedback_Del_Device:
			pkg := protocol.Parse_Feedback_Del_Device(pkgbyte, smconn.ID)
			smconn.ReadMore = false
			return NewShaPacket(protocol.Feedback_Del_Device, pkg), nil
		case protocol.Feedback_Query_Attr:
			pkg := protocol.Parse_Feedback_Query_Attr(pkgbyte, smconn.ID)
			smconn.ReadMore = false
			return NewShaPacket(protocol.Feedback_Query_Attr, pkg), nil
		case protocol.Feedback_Depolyment:
			pkg := protocol.Parse_Feedback_Deployment(pkgbyte, smconn.ID)
			smconn.ReadMore = false
			return NewShaPacket(protocol.Feedback_Depolyment, pkg), nil
		case protocol.Feedback_OnOff:
			pkg := protocol.Parse_Feedback_Onoff(pkgbyte, smconn.ID)
			smconn.ReadMore = false
			return NewShaPacket(protocol.Feedback_OnOff, pkg), nil

		case protocol.Feedback_Level_Control:
			pkg := protocol.Parse_Feedback_Level_Control(pkgbyte, smconn.ID)
			smconn.ReadMore = false
			return NewShaPacket(protocol.Feedback_Level_Control, pkg), nil

		case protocol.Illegal:
			smconn.ReadMore = true
		case protocol.HalfPack:
			smconn.ReadMore = true
		}
	}

}
