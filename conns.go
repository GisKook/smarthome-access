package sha

import (
	"sync/atomic"
)

type Conns struct {
	connsindex map[uint32]*Conn
	connsuid   map[string]uint32
	index      uint32
}

func NewConns() {
	return &Conns{
		connsindex: make(map[uint32]*Conn),
		connsuid:   make(map[string]uint32),
		index:      0,
	}
}

func (cs *Conns) Add(conn *Conn) {
	conn.index = atomic.AddUInt32(&cs.index, 1)
}

func (cs *Conns) GetConn(uid string) *Conn {
	return cs.connsindex[connsuid[uid]]
}

func (cs *Conns) Remove(uid string) {
	index := connsuid[uid]
	delete(cs.connsindex, index)
	delete(cs.connsuid, uid)
}
