package sha

import (
	"fmt"
	"github.com/giskook/gotcp"
	"github.com/giskook/smarthome-access/base"
	"github.com/giskook/smarthome-access/protocol"
	"log"
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
	NewConns().Add(conn)

	return true
}

func (this *Callback) OnClose(c *gotcp.Conn) {
	conn := c.GetExtraData().(*Conn)
	conn.Close()
	NewConns().Remove(conn)
}

func on_login(c *gotcp.Conn, p *ShaPacket) {
	conn := c.GetExtraData().(*Conn)
	conn.Status = ConnSuccess
	loginPkg := p.Packet.(*protocol.LoginPacket)
	conn.Gateway = loginPkg.Gateway
	conn.ID = conn.Gateway.ID
	NewConns().SetID(conn.ID, conn)
	c.AsyncWritePacket(p, time.Second)
}

func on_add_del_device(c *gotcp.Conn, p *ShaPacket) {
	conn := c.GetExtraData().(*Conn)
	add_del_device_pkg := p.Packet.(*protocol.Add_Del_Device_Packet)

	if add_del_device_pkg.Action == protocol.ADD_DEVICE {
		if !base.Gateway_Check_Device(conn.Gateway, add_del_device_pkg.Device.ID) {
			base.Gateway_Add_Device(conn.Gateway, add_del_device_pkg.Device)
		} else {
			base.Gateway_Update_Davice(conn.Gateway, add_del_device_pkg.Device)
		}
		GetServer().GetProducer().Send(GetConfiguration().NsqConfig.UpTopic, p.Serialize())
	} else if add_del_device_pkg.Action == protocol.DEL_DEVICE {
		if base.Gateway_Check_Device(conn.Gateway, add_del_device_pkg.Device.ID) {
			base.Gateway_Del_Device(conn.Gateway, conn.ID)
			GetServer().GetProducer().Send(GetConfiguration().NsqConfig.UpTopic, p.Serialize())
		}
	}
}

func on_feedback_set_name(c *gotcp.Conn, p *ShaPacket) {
	conn := c.GetExtraData().(*Conn)
	feedback_set_name_pkg := p.Packet.(*protocol.Feedback_SetName_Packet)
	base.Gateway_Set_Device_Name(conn.Gateway, feedback_set_name_pkg.DeviceID, feedback_set_name_pkg.DeviceName)
	fmt.Printf("%+v\n", conn.Gateway)
	fmt.Printf("dviceid %d\n", feedback_set_name_pkg.DeviceID)
	GetServer().GetProducer().Send(GetConfiguration().NsqConfig.UpTopic, p.Serialize())
}

func on_feedback_del_device(c *gotcp.Conn, p *ShaPacket) {
	conn := c.GetExtraData().(*Conn)
	feedback_del_device_pkg := p.Packet.(*protocol.Feedback_Del_Device_Packet)
	if feedback_del_device_pkg.Result == 0 {
		base.Gateway_Del_Device(conn.Gateway, feedback_del_device_pkg.DeviceID)
	}
	GetServer().GetProducer().Send(GetConfiguration().NsqConfig.UpTopic, p.Serialize())
}

func on_feedback_onoff(c *gotcp.Conn, p *ShaPacket) {
	conn := c.GetExtraData().(*Conn)
	feedback_onoff_pkg := p.Packet.(*protocol.Feedback_OnOff_Packet)
	if feedback_onoff_pkg.Result == 0 {
		base.Gateway_Set_Device_Status(conn.Gateway, feedback_onoff_pkg.DeviceID, feedback_onoff_pkg.Endpoint, feedback_onoff_pkg.Action)
	}
	GetServer().GetProducer().Send(GetConfiguration().NsqConfig.UpTopic, p.Serialize())
}

func on_feedback_level_control(c *gotcp.Conn, p *ShaPacket) {
	conn := c.GetExtraData().(*Conn)
	feedback_level_control_pkg := p.Packet.(*protocol.Feedback_Level_Control_Packet)
	if feedback_level_control_pkg.Result == 0 {
		base.Gateway_Set_Device_Status(conn.Gateway, feedback_level_control_pkg.DeviceID, feedback_level_control_pkg.Endpoint, feedback_level_control_pkg.Action)
	}
	GetServer().GetProducer().Send(GetConfiguration().NsqConfig.UpTopic, p.Serialize())
}

func on_notify_onoff(c *gotcp.Conn, p *ShaPacket) {
	log.Println("on notify onoff")
	GetServer().GetProducer().Send(GetConfiguration().NsqConfig.UpTopic, p.Serialize())
	conn := c.GetExtraData().(*Conn)
	notify_onoff_pkg := p.Packet.(*protocol.Notify_OnOff_Packet)
	base.Gateway_Set_Device_Status(conn.Gateway, notify_onoff_pkg.DeviceID, notify_onoff_pkg.EndPoint, notify_onoff_pkg.Status)
}

func on_notify_online(c *gotcp.Conn, p *ShaPacket) {
	log.Println("on notify online")
	conn := c.GetExtraData().(*Conn)
	notify_online_pkg := p.Packet.(*protocol.Notify_Online_Status_Packet)
	base.Gateway_Set_Device_Online(conn.Gateway, notify_online_pkg.DeviceID, notify_online_pkg.Status)

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
		GetServer().GetProducer().Send(GetConfiguration().NsqConfig.UpTopic, p.Serialize())
	case protocol.Add_Del_Device:
		on_add_del_device(c, shaPacket)
	case protocol.Feedback_SetName:
		on_feedback_set_name(c, shaPacket)
	case protocol.Feedback_Del_Device:
		on_feedback_del_device(c, shaPacket)
	case protocol.Feedback_Query_Attr:
		GetServer().GetProducer().Send(GetConfiguration().NsqConfig.UpTopic, p.Serialize())
	case protocol.Feedback_Depolyment:
		GetServer().GetProducer().Send(GetConfiguration().NsqConfig.UpTopic, p.Serialize())
	case protocol.Feedback_OnOff:
		on_feedback_onoff(c, shaPacket)
	case protocol.Feedback_Level_Control:
		on_feedback_level_control(c, shaPacket)
	case protocol.Notify_OnOff:
		on_notify_onoff(c, shaPacket)
	case protocol.Notify_Level:
		GetServer().GetProducer().Send(GetConfiguration().NsqConfig.UpTopic, p.Serialize())
	case protocol.Notify_Online_Status:
		on_notify_online(c, shaPacket)
	}

	return true
}
