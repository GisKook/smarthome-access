package protocol

import (
	"github.com/giskook/smarthome-access/base"
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

type Feedback_Deployment_Packet struct {
	GatewayID     uint64
	SerialNum     uint32
	DeviceID      uint64
	EndPoint      uint8
	ArmModel      uint8
	StartTimeHour uint8
	StartTimeMin  uint8
	EndTimeHour   uint8
	EndTimeMin    uint8
	Result        uint8
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
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(p.ArmModel),
		},
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(p.StartTimeHour),
		},
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(p.StartTimeMin),
		},
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(p.EndTimeHour),
		},
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(p.EndTimeMin),
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
	endpoint, _ := reader.ReadByte()
	armmodel, _ := reader.ReadByte()
	start_time_hour, _ := reader.ReadByte()
	start_time_min, _ := reader.ReadByte()
	end_time_hour, _ := reader.ReadByte()
	end_time_min, _ := reader.ReadByte()
	result, _ := reader.ReadByte()

	return &Feedback_Deployment_Packet{
		GatewayID:     id,
		SerialNum:     serialnum,
		DeviceID:      deviceid,
		EndPoint:      endpoint,
		ArmModel:      armmodel,
		StartTimeHour: start_time_hour,
		StartTimeMin:  start_time_min,
		EndTimeHour:   end_time_hour,
		EndTimeMin:    end_time_min,
		Result:        result,
	}
}
