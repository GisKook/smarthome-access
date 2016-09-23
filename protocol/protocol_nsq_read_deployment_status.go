package protocol

import (
	"bytes"
	"github.com/giskook/smarthome-access/base"
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

const CMD_READ_DEPLOYMENT_STATUS uint16 = 0x801f
const CMD_READ_DEPLOYMENT_STATUS_LEN uint16 = 0x0014

type Nsq_Read_Deployment_Status_Packet struct {
	GatewayID uint64
	SerialNum uint32
	Status    byte
	DeviceID  uint64
	Endpoint  uint8
}

func (p *Nsq_Read_Deployment_Status_Packet) SerializeOnline() []byte {
	var writer bytes.Buffer
	writer.WriteByte(STARTFLAG)
	base.WriteWord(&writer, CMD_READ_DEPLOYMENT_STATUS_LEN)
	base.WriteWord(&writer, CMD_READ_DEPLOYMENT_STATUS)
	base.WriteDWord(&writer, p.SerialNum)
	base.WriteQuaWord(&writer, p.DeviceID)
	writer.WriteByte(CheckSum(writer.Bytes(), uint16(writer.Len())))
	writer.WriteByte(ENDFLAG)

	return writer.Bytes()
}

func (p *Nsq_Read_Deployment_Status_Packet) SerializeOffline() []byte {
	para := []*Report.Command_Param{
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(p.Status),
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

func (p *Nsq_Read_Deployment_Status_Packet) Serialize() []byte {
	if p.Status == GATEWAY_ON_LINE {
		return p.SerializeOnline()
	} else {
		return p.SerializeOffline()
	}
}

func Parse_NSQ_Read_Deployment_Status(gatewayid uint64, serialnum uint32, status byte, paras []*Report.Command_Param) *Nsq_Read_Deployment_Status_Packet {
	deviceid := paras[0].Npara
	endpint := uint8(paras[1].Npara)

	return &Nsq_Read_Deployment_Status_Packet{
		GatewayID: gatewayid,
		SerialNum: serialnum,
		Status:    status,
		DeviceID:  deviceid,
		Endpoint:  endpint,
	}
}
