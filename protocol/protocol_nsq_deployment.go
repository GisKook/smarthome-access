package protocol

import (
	"bytes"
	"github.com/giskook/smarthome-access/base"
	"github.com/giskook/smarthome-access/pb"
)

const CMD_DEPLOYMENT uint16 = 0x800f
const CMD_DEPLOYMENT_LEN uint16 = 0x0019

const CMD_DEPLOYMENT_ARM uint8 = 0
const CMD_DEPLOYMENT_DISARM uint8 = 1
const CMD_DEPLOYMENT_ARMTIME uint8 = 2

type Nsq_Deployment_Packet struct {
	SerialNum        uint32
	DeviceID         uint64
	Endpoint         uint8
	ArmMode          uint8
	ArmStartTimeHour uint8
	ArmStartTimeMin  uint8
	ArmEndTImeHour   uint8
	ArmEndTImeMin    uint8
}

func (p *Nsq_Deployment_Packet) Serialize() []byte {
	var writer bytes.Buffer
	writer.WriteByte(STARTFLAG)
	base.WriteWord(&writer, CMD_DEPLOYMENT_LEN)
	base.WriteWord(&writer, CMD_DEPLOYMENT)
	base.WriteQuaWord(&writer, p.DeviceID)
	writer.WriteByte(byte(p.Endpoint))
	base.WriteDWord(&writer, p.SerialNum)
	writer.WriteByte(byte(p.ArmMode))
	writer.WriteByte(byte(p.ArmStartTimeHour))
	writer.WriteByte(byte(p.ArmStartTimeMin))
	writer.WriteByte(byte(p.ArmEndTImeHour))
	writer.WriteByte(byte(p.ArmEndTImeMin))
	writer.WriteByte(CheckSum(writer.Bytes(), uint16(writer.Len())))
	writer.WriteByte(ENDFLAG)

	return writer.Bytes()
}

func Parse_NSQ_Deployment(serialnum uint32, paras []*Report.Command_Param) *Nsq_Deployment_Packet {
	deviceid := paras[0].Npara
	endpint := uint8(paras[1].Npara)
	armmode := uint8(paras[2].Npara)
	armstarthour := uint8(paras[3].Npara)
	armstartmin := uint8(paras[4].Npara)
	armendhour := uint8(paras[5].Npara)
	armendmin := uint8(paras[6].Npara)

	return &Nsq_Deployment_Packet{
		SerialNum:        serialnum,
		DeviceID:         deviceid,
		Endpoint:         endpint,
		ArmMode:          armmode,
		ArmStartTimeHour: armstarthour,
		ArmStartTimeMin:  armstartmin,
		ArmEndTImeHour:   armendhour,
		ArmEndTImeMin:    armendmin,
	}
}
