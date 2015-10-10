package sha

import (
	"fmt"
	"sync"
)

type GatewayProperty struct {
	Uid         uint64
	Name        string
	Devicecount uint16
	Devicelist  []Device
}

type GatewayHub struct {
	Gateway map[uint64]*GatewayProperty

	waitGroup *sync.WaitGroup
}

var gatewayhub *GatewayHub

func (g *GatewayHub) Add(gatewayid uint64, deviceid uint64, devicetype uint8, company uint16, status uint8, name string) {
	_, ok := g.Gateway[gatewayid]
	if !ok {
		g.Gateway[gatewayid] = &GatewayProperty{
			Uid:         gatewayid,
			Devicecount: 1,
			Devicelist: []Device{
				{
					Oid:     deviceid,
					Type:    devicetype,
					Company: company,
					Status:  status,
					Name:    name,
				},
			},
		}
	} else {
		g.Gateway[gatewayid].Devicelist = append(g.Gateway[gatewayid].Devicelist, Device{
			Oid:     deviceid,
			Type:    devicetype,
			Company: company,
			Status:  status,
			Name:    name,
		})
		g.Gateway[gatewayid].Devicecount++
	}
}

func (g *GatewayHub) Insert(gatewayid uint64, name string, devicecount uint16, devicelist []Device) {
	g.Gateway[gatewayid] = &GatewayProperty{
		Uid:         gatewayid,
		Devicecount: devicecount,
		Name:        name,
		Devicelist:  devicelist,
	}
}

func NewGatewayHub() *GatewayHub {
	if gatewayhub == nil {
		gatewayhub = &GatewayHub{
			Gateway:   make(map[uint64]*GatewayProperty),
			waitGroup: &sync.WaitGroup{},
		}
	}

	return gatewayhub
}

func (g *GatewayHub) Del(gatewayid uint64, deviceid uint64) {
	_, ok := g.Gateway[gatewayid]
	if ok {
		devicelist := g.Gateway[gatewayid].Devicelist
		var i uint16
		for i = 0; i < g.Gateway[gatewayid].Devicecount; i++ {
			if deviceid == devicelist[i].Oid {
				devicelist = append(devicelist[:i], devicelist[i+1:]...)
				g.Gateway[gatewayid].Devicecount--
				g.Gateway[gatewayid].Devicelist = devicelist
				break
			}
		}
	}
	fmt.Println(g.Gateway[gatewayid].Devicelist)
}

func (g *GatewayHub) Remove(gatewayid uint64) {
	delete(g.Gateway, gatewayid)
}

func (g *GatewayHub) Setname(gatewayid uint64, deviceid uint64, name string) {
	_, ok := g.Gateway[gatewayid]
	if ok {
		devicelist := g.Gateway[gatewayid].Devicelist
		var i uint16
		for i = 0; i < g.Gateway[gatewayid].Devicecount; i++ {
			if deviceid == devicelist[i].Oid {
				devicelist[i].Name = name
				break
			}
		}

	}
	fmt.Println(g.Gateway[gatewayid].Devicelist)
}

func (g *GatewayHub) GetGateway(gatewayid uint64) *GatewayProperty {
	gateway, _ := g.Gateway[gatewayid]

	return gateway
}

func (g *GatewayHub) Check(gatewayid uint64) bool {
	_, ok := g.Gateway[gatewayid]

	return ok
}
