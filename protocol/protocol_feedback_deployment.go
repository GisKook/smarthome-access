package protocol

import (
	"github.com/giskook/smarthome-access/base"
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

type Feedback_Deployment_Packet struct {
	GatewayID uint64
	DeviceID  uint64
	SerialNum uint32
	Result    uint8
}

func (p *Feedback_Deployment_Packet) Serialize() []byte {
	para := []*Report.Command_Param{
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(GATEWAY_ON_LINE),
		},
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(p.Result),
		},
	}

	command := &Report.Command{
		Type:  Report.Command_CMT_REP_DEPLOYMENT,
		Paras: para,
	}

	feedback_depolyment_pkg := &Report.ControlReport{
		Tid:          p.GatewayID,
		SerialNumber: p.SerialNum,
		Command:      command,
	}

	data, _ := proto.Marshal(feedback_depolyment_pkg)

	return data
}

func Parse_Feedback_Deployment(buffer []byte, id uint64) *Feedback_Deployment_Packet {
	reader := ParseHeader(buffer)
	serialnum := base.ReadDWord(reader)
	deviceid := base.ReadQuaWord(reader)
	result, _ := reader.ReadByte()

	return &Feedback_Deployment_Packet{
		GatewayID: id,
		DeviceID:  deviceid,
		SerialNum: serialnum,
		Result:    result,
	}
}
