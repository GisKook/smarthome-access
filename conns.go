package sha

import (
	"sync"
	"sync/atomic"
)

type Conns struct {
	connsindex  map[uint32]*Conn
	mutex_index sync.Mutex
	connsuid    map[uint64]*Conn
	mutex_uid   sync.Mutex
	index       uint32
}

var g_mutex_conns sync.Mutex
var connsInstance *Conns

func NewConns() *Conns {
	defer g_mutex_conns.Unlock()
	g_mutex_conns.Lock()
	if connsInstance == nil {
		connsInstance = &Conns{
			connsindex: make(map[uint32]*Conn),
			connsuid:   make(map[uint64]*Conn),
			index:      0,
		}
	}

	return connsInstance
}

func (cs *Conns) Add(conn *Conn) {
	conn.index = atomic.AddUint32(&cs.index, 1)
	cs.mutex_index.Lock()
	cs.connsindex[conn.index] = conn
	cs.mutex_index.Unlock()
}

func (cs *Conns) SetID(gatewayid uint64, conn *Conn) {
	cs.mutex_uid.Lock()
	cs.connsuid[gatewayid] = conn
	cs.mutex_uid.Unlock()
}

func (cs *Conns) GetConn(uid uint64) *Conn {
	defer cs.mutex_uid.Unlock()
	cs.mutex_uid.Lock()
	return cs.connsuid[uid]
}

func (cs *Conns) Remove(c *Conn) {
	cs.mutex_index.Lock()
	delete(cs.connsindex, c.index)
	cs.mutex_index.Unlock()

	cs.mutex_uid.Lock()
	connuid, ok := cs.connsuid[c.ID]
	if ok && c.index == connuid.index {
		delete(cs.connsuid, c.ID)
	}
	cs.mutex_uid.Unlock()
}

func (cs *Conns) Check(uid uint64) bool {
	cs.mutex_uid.Lock()
	conn, ok := cs.connsuid[uid]
	cs.mutex_uid.Unlock()
	if ok {
		cs.mutex_index.Lock()
		_, realok := cs.connsindex[conn.index]
		cs.mutex_index.Unlock()

		return realok
	}
	return ok
}

func (cs *Conns) GetCount() int {
	defer cs.mutex_index.Unlock()
	cs.mutex_index.Lock()
	return len(cs.connsindex)
}
