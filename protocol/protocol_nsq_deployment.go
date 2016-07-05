package protocol

import (
	"bytes"
	"github.com/giskook/smarthome-access/base"
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

const CMD_DEPLOYMENT uint16 = 0x800f
const CMD_DEPLOYMENT_LEN uint16 = 0x0019

const CMD_DEPLOYMENT_ARM uint8 = 0
const CMD_DEPLOYMENT_DISARM uint8 = 1
const CMD_DEPLOYMENT_ARMTIME uint8 = 2

type Nsq_Deployment_Packet struct {
	GatewayID        uint64
	SerialNum        uint32
	Status           byte
	DeviceID         uint64
	Endpoint         uint8
	ArmMode          uint8
	ArmStartTimeHour uint8
	ArmStartTimeMin  uint8
	ArmEndTImeHour   uint8
	ArmEndTImeMin    uint8
}

func (p *Nsq_Deployment_Packet) SerializeOnline() []byte {
	var writer bytes.Buffer
	writer.WriteByte(STARTFLAG)
	base.WriteWord(&writer, CMD_DEPLOYMENT_LEN)
	base.WriteWord(&writer, CMD_DEPLOYMENT)
	base.WriteDWord(&writer, p.SerialNum)
	base.WriteQuaWord(&writer, p.DeviceID)
	writer.WriteByte(byte(p.Endpoint))
	writer.WriteByte(byte(p.ArmMode))
	writer.WriteByte(byte(p.ArmStartTimeHour))
	writer.WriteByte(byte(p.ArmStartTimeMin))
	writer.WriteByte(byte(p.ArmEndTImeHour))
	writer.WriteByte(byte(p.ArmEndTImeMin))
	writer.WriteByte(CheckSum(writer.Bytes(), uint16(writer.Len())))
	writer.WriteByte(ENDFLAG)

	return writer.Bytes()
}

func (p *Nsq_Deployment_Packet) SerializeOffline() []byte {
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

func (p *Nsq_Deployment_Packet) Serialize() []byte {
	if p.Status == GATEWAY_ON_LINE {
		return p.SerializeOnline()
	} else {
		return p.SerializeOffline()
	}
}

func Parse_NSQ_Deployment(gatewayid uint64, serialnum uint32, status byte, paras []*Report.Command_Param) *Nsq_Deployment_Packet {
	deviceid := paras[0].Npara
	endpint := uint8(paras[1].Npara)
	armmode := uint8(paras[2].Npara)
	armstarthour := (uint8(paras[3].Npara) + 24 - 8) % 24
	armstartmin := uint8(paras[4].Npara)
	armendhour := (uint8(paras[5].Npara) + 24 - 8) % 24
	armendmin := uint8(paras[6].Npara)

	return &Nsq_Deployment_Packet{
		GatewayID:        gatewayid,
		SerialNum:        serialnum,
		Status:           status,
		DeviceID:         deviceid,
		Endpoint:         endpint,
		ArmMode:          armmode,
		ArmStartTimeHour: armstarthour,
		ArmStartTimeMin:  armstartmin,
		ArmEndTImeHour:   armendhour,
		ArmEndTImeMin:    armendmin,
	}
}
