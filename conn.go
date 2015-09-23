package sha

import (
	"bytes"
	"github.com/giskook/gotcp"
	"log"
	"time"
)

var ConnSuccess uint8 = 0
var ConnUnauth uint8 = 1

type ConnConfig struct {
	HeartBeat    time.Duration
	ReadLimit    int64
	WriteLimit   int64
	NsqChanLimit int32
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
	uid                  uint64
	status               uint8
}

func NewConn(conn *gotcp.Conn, config *ConnConfig) *Conn {
	return &Conn{
		conn:                 conn,
		recieveBuffer:        bytes.NewBuffer([]byte{}),
		config:               config,
		readflag:             time.Now().Unix(),
		writeflag:            time.Now().Unix(),
		ticker:               time.NewTicker(config.HeartBeat * 1e9),
		packetNsqReceiveChan: make(chan gotcp.Packet, config.NsqChanLimit),
		closeChan:            make(chan bool),
		index:                0,
		status:               ConnUnauth,
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

func (c *Conn) SetStatus(status uint8) {
	c.status = status
}

func (c *Conn) checkHeart() {
	defer func() {
		c.conn.Close()
	}()

	for {
		<-c.ticker.C
		now := time.Now().Unix()
		if now-c.readflag > c.config.ReadLimit {
			log.Println("read linmit")
			return
		}
		if now-c.writeflag > c.config.WriteLimit {
			log.Println("write limit")
			return
		}
		if c.status == ConnUnauth {
			log.Println("status")
			return
		}
		if <-c.closeChan {
			log.Println("close status")
			return
		}
	}
}

func (c *Conn) Do() {
	go c.writeToclientLoop()
	go c.checkHeart()
}

type Callback struct{}

func (this *Callback) OnConnect(c *gotcp.Conn) bool {
	config := &ConnConfig{
		HeartBeat:  6,
		ReadLimit:  600,
		WriteLimit: 600,
	}
	conn := NewConn(c, config)

	c.PutExtraData(conn)

	conn.Do()

	return true
}

func (this *Callback) OnClose(c *gotcp.Conn) {
	conn := c.GetExtraData().(*Conn)
	conn.Close()
}

func (this *Callback) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	shaPacket := p.(*ShaPacket)
	switch shaPacket.Type {
	case Login:
		c.AsyncWritePacket(shaPacket, time.Second)
	case HeartBeat:
		c.AsyncWritePacket(shaPacket, time.Second)
	case SendDeviceList:
		GetServer().GetProducer().Send(GetServer().GetTopic(), p.Serialize())
	case OperateFeedback:
		GetServer().GetProducer().Send(GetServer().GetTopic(), p.Serialize())
	case Warn:
		GetServer().GetProducer().Send(GetServer().GetTopic(), p.Serialize())
	case AddDelDevice:
		GetServer().GetProducer().Send(GetServer().GetTopic(), p.Serialize())

	}

	return true
}
