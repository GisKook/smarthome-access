package protocol

import (
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

type Notify_Gateway_Offline_Packet struct {
	GatewayID uint64
}

func (p *Notify_Gateway_Offline_Packet) Serialize() []byte {
	command := &Report.Command{
		Type:  Report.Command_CMT_REP_NOTIFY_GATEWAY_OFFLINE,
		Paras: nil,
	}

	notify_gateway_offline_pkg := &Report.ControlReport{
		Tid:          p.GatewayID,
		SerialNumber: 0,
		Command:      command,
	}

	data, _ := proto.Marshal(notify_gateway_offline_pkg)

	return data
}
