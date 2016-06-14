package protocol

import (
	"github.com/giskook/smarthome-access/base"
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

const ON uint8 = 1
const OFF uint8 = 0

type Feedback_OnOff_Packet struct {
	GatewayID uint64
	DeviceID  uint64
	ClusterID uint16
	Endpoint  uint8
	Action    uint8
	Result    uint8
}

func (p *Feedback_OnOff_Packet) Serialize() []byte {
	para := []*Report.Command_Param{
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(p.Result),
		},
	}

	command := &Report.Command{
		Type:  Report.Command_CMT_REP_ONOFF,
		Paras: para,
	}

	feedback_onoff_pkg := &Report.ControlReport{
		Tid:          p.GatewayID,
		SerialNumber: 0,
		Command:      command,
	}

	data, _ := proto.Marshal(feedback_onoff_pkg)

	return data
}

func Parse_Feedback_Onoff(buffer []byte, id uint64) *Feedback_Deployment_Packet {
	reader := ParseHeader(buffer)
	deviceid := base.ReadQuaWord(reader)
	clusterid := base.ReadWord(reader)
	action, _ := reader.ReadByte()
	result, _ := reader.ReadByte()

	return &Feedback_OnOff_Packet{
		GatewayID: id,
		ClusterID: deviceid,
		Endpoint:  clusterid,
		Action:    action,
		Result:    result,
	}
}
