package sha

import (
	"bytes"
	"fmt"
	"github.com/giskook/gotcp"
	"time"
)

var ConnSuccess uint8 = 0
var ConnNormal uint8 = 1

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
	closeChan            chan struct{}
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
		closeChan:            make(chan struct{}),
		index:                0,
		status:               ConnNormal,
	}
}

func (c *Conn) Close() {
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
		case <-c.closeChan:
			return
		case p := <-c.packetNsqReceiveChan:
			if _, err := c.conn.GetRawConn().Write(p.Serialize()); err != nil {
				return
			}
		}
	}
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

	for {
		<-c.ticker.C
		now := time.Now().Unix()
		fmt.Println(now - c.readflag)
		if now-c.readflag > c.config.ReadLimit {
			return
		}
		if now-c.writeflag > c.config.WriteLimit {
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
		HeartBeat:  1,
		ReadLimit:  6,
		WriteLimit: 6,
	}
	conn := NewConn(c, config)

	c.PutExtraData(conn)

	conn.Do()

	NewConns().Add(conn)

	return true
}

func (this *Callback) OnClose(c *gotcp.Conn) {
	conn := c.GetExtraData()
	element, _ := conn.(*Conn)

	element.Close()
}

func (this *Callback) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	shaPacket := p.(*ShaPacket)
	c.AsyncWritePacket(shaPacket, time.Second)

	return true
}
