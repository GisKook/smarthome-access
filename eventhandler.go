package sha

import (
	"fmt"
	"github.com/giskook/gotcp"
	"github.com/giskook/smarthome-access/base"
	"github.com/giskook/smarthome-access/protocol"
	"time"
)

type Callback struct{}

func (this *Callback) OnConnect(c *gotcp.Conn) bool {
	checkinterval := GetConfiguration().GetServerConnCheckInterval()
	readlimit := GetConfiguration().GetServerReadLimit()
	writelimit := GetConfiguration().GetServerWriteLimit()
	config := &ConnConfig{
		ConnCheckInterval: uint16(checkinterval),
		ReadLimit:         uint16(readlimit),
		WriteLimit:        uint16(writelimit),
	}
	conn := NewConn(c, config)

	c.PutExtraData(conn)

	conn.Do()

	return true
}

func (this *Callback) OnClose(c *gotcp.Conn) {
	conn := c.GetExtraData().(*Conn)
	conn.Close()
	NewConns().Remove(conn.ID)
}

func on_login(c *gotcp.Conn, p *ShaPacket) {
	conn := c.GetExtraData().(*Conn)
	conn.Status = ConnSuccess
	loginPkg := p.Packet.(*protocol.LoginPacket)
	conn.Gateway = loginPkg.Gateway
	conn.ID = conn.Gateway.ID
}

func on_add_del_device(c *gotcp.Conn, p *ShaPacket) {
	conn := c.GetExtraData().(*Conn)
	add_del_device_pkg := p.Packet.(*protocol.Add_Del_Device_Packet)

	if add_del_device_pkg.Action == protocol.ADD_DEVICE {
		base.Gateway_Add_Device(conn.Gateway, add_del_device_pkg.Device)
	} else if add_del_device_pkg.Action == protocol.DEL_DEVICE {
		base.Gateway_Del_Device(conn.Gateway, conn.ID)
	}
	fmt.Printf("%+v\n", *conn.Gateway)
	GetServer().GetProducer().Send(NSQ_CONTROL_PUB_TOPIC, p.Serialize())
}

func on_feedback_set_name(c *gotcp.Conn, p *ShaPacket) {
	conn := c.GetExtraData().(*Conn)
	feedback_set_name_pkg := p.Packet.(*protocol.Feedback_SetName_Packet)
	base.Gateway_Set_Device_Name(conn.Gateway, feedback_set_name_pkg.DeviceID, feedback_set_name_pkg.DeviceName)
	GetServer().GetProducer().Send(NSQ_CONTROL_PUB_TOPIC, p.Serialize())
}

func (this *Callback) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	shaPacket := p.(*ShaPacket)
	switch shaPacket.Type {
	case protocol.Login:
		//	c.AsyncWritePacket(shaPacket, time.Second)
		on_login(c, shaPacket)
	case protocol.HeartBeat:
		c.AsyncWritePacket(shaPacket, time.Second)
		//GetServer().GetProducer().Send(GetServer().GetTopic(), p.Serialize())
	case protocol.Notification:
		GetServer().GetProducer().Send(NSQ_CONTROL_PUB_TOPIC, p.Serialize())
	case protocol.Add_Del_Device:
		on_add_del_device(c, shaPacket)
	case protocol.Feedback_SetName:
		on_feedback_set_name(c, shaPacket)
	case protocol.Feedback_Del_Device:
		GetServer().GetProducer().Send(NSQ_CONTROL_PUB_TOPIC, p.Serialize())
	case protocol.Feedback_Query_Attr:
		GetServer().GetProducer().Send(NSQ_CONTROL_PUB_TOPIC, p.Serialize())
	case protocol.Feedback_Depolyment:
		GetServer().GetProducer().Send(NSQ_CONTROL_PUB_TOPIC, p.Serialize())
	case protocol.Feedback_OnOff:
		GetServer().GetProducer().Send(NSQ_CONTROL_PUB_TOPIC, p.Serialize())

	}

	return true
}
