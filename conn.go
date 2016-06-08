package sha

import (
	"bytes"
	"github.com/giskook/gotcp"
	"github.com/giskook/smarthome-access/base"
	"github.com/giskook/smarthome-access/protocol"
	"log"
	"time"
)

var ConnSuccess uint8 = 0
var ConnUnauth uint8 = 1

type ConnConfig struct {
	ConnCheckInterval uint16
	ReadLimit         uint16
	WriteLimit        uint16
	NsqChanLimit      uint16
}

type Conn struct {
	conn                 *gotcp.Conn
	config               *ConnConfig
	recieveBuffer        *bytes.Buffer
	ticker               *time.Ticker
	readflag             int64
	writeflag            int64
	packetNsqReceiveChan chan gotcp.Packet
	closeChan            chan bool
	index                uint32
	ID                   uint64
	Status               uint8
	Gateway              *base.Gateway
}

func NewConn(conn *gotcp.Conn, config *ConnConfig) *Conn {
	return &Conn{
		conn:                 conn,
		recieveBuffer:        bytes.NewBuffer([]byte{}),
		config:               config,
		readflag:             time.Now().Unix(),
		writeflag:            time.Now().Unix(),
		ticker:               time.NewTicker(time.Duration(config.ConnCheckInterval) * time.Second),
		packetNsqReceiveChan: make(chan gotcp.Packet, config.NsqChanLimit),
		closeChan:            make(chan bool),
		index:                0,
		Status:               ConnUnauth,
	}
}

func (c *Conn) Close() {
	c.closeChan <- true
	c.ticker.Stop()
	c.recieveBuffer.Reset()
	close(c.packetNsqReceiveChan)
	close(c.closeChan)
}

func (c *Conn) GetBuffer() *bytes.Buffer {
	return c.recieveBuffer
}

func (c *Conn) writeToclientLoop() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case p := <-c.packetNsqReceiveChan:
			if p != nil {
				c.conn.GetRawConn().Write(p.Serialize())
			}
		case <-c.closeChan:
			return
		}
	}
}

func (c *Conn) SendToGateway(p gotcp.Packet) {
	c.packetNsqReceiveChan <- p
}

func (c *Conn) UpdateReadflag() {
	c.readflag = time.Now().Unix()
}

func (c *Conn) UpdateWriteflag() {
	c.writeflag = time.Now().Unix()
}

func (c *Conn) checkHeart() {
	defer func() {
		c.conn.Close()
	}()

	var now int64
	for {
		select {
		case <-c.ticker.C:
			now = time.Now().Unix()
			if now-c.readflag > int64(c.config.ReadLimit) {
				log.Println("read linmit")
				return
			}
			if now-c.writeflag > int64(c.config.WriteLimit) {
				log.Println("write limit")
				return
			}
			if c.Status == ConnUnauth {
				log.Printf("unauth's gateway gatewayid %x\n", c.ID)
				return
			}
		case <-c.closeChan:
			return
		}
	}
}

func (c *Conn) Do() {
	go c.checkHeart()
	go c.writeToclientLoop()
}

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
	case protocol.HeartBeat:
		//GetServer().GetProducer().Send(GetServer().GetTopic(), p.Serialize())

	}

	return true
}
