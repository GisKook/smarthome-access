package sha

import (
	"sync"
	"sync/atomic"
)

type Conns struct {
	connsindex map[uint32]*Conn
	connsuid   map[uint64]*Conn
	index      uint32
	mutex      *sync.RWMutex
}

var connsInstance *Conns

func NewConns() *Conns {
	if connsInstance == nil {
		connsInstance = &Conns{
			connsindex: make(map[uint32]*Conn),
			connsuid:   make(map[uint64]*Conn),
			index:      0,
			mutex:      new(sync.RWMutex),
		}
	}

	return connsInstance
}

func (cs *Conns) Add(conn *Conn) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	conn.index = atomic.AddUint32(&cs.index, 1)
	cs.connsindex[conn.index] = conn
}

func (cs *Conns) SetID(gatewayid uint64, conn *Conn) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	cs.connsuid[gatewayid] = conn
}

func (cs *Conns) GetConn(uid uint64) *Conn {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	return cs.connsuid[uid]
}

func (cs *Conns) Remove(c *Conn) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	delete(cs.connsindex, c.index)

	connuid, ok := cs.connsuid[c.ID]
	if ok && c.index == connuid.index {
		delete(cs.connsuid, c.ID)
	}
}

func (cs *Conns) Check(uid uint64) bool {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	conn, ok := cs.connsuid[uid]
	if ok {
		_, realok := cs.connsindex[conn.index]

		return realok
	}
	return ok
}

func (cs *Conns) GetCount() int {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	return len(cs.connsindex)
}
