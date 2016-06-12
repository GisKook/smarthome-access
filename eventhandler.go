package sha

import (
	"github.com/giskook/gotcp"
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

func (this *Callback) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	shaPacket := p.(*ShaPacket)
	switch shaPacket.Type {
	case protocol.Login:
		//	c.AsyncWritePacket(shaPacket, time.Second)
		conn := c.GetExtraData().(*Conn)
		conn.Status = ConnSuccess
	case protocol.HeartBeat:
		c.AsyncWritePacket(shaPacket, time.Second)
		//GetServer().GetProducer().Send(GetServer().GetTopic(), p.Serialize())
	case protocol.Add_Del_Device:

		GetServer().GetProducer().Send(NSQ_CONTROL_PUB_TOPIC, p.Serialize())
	}

	return true
}
