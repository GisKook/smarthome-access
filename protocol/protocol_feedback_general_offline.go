package protocol

import (
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

type Feedback_General_Offline_Package struct {
	GatewayID uint64
	SerialNum uint32
}

func (p *Feedback_General_Offline_Package) Serialize() []byte {

	command := &Report.Command{
		Type:  Report.Command_CMT_REP_GENERAL_OFFLINE,
		Paras: nil,
	}

	feedback_onoff_pkg := &Report.ControlReport{
		Tid:          p.GatewayID,
		SerialNumber: p.SerialNum,
		Command:      command,
	}

	data, _ := proto.Marshal(feedback_onoff_pkg)

	return data
}
