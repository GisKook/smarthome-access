package protocol

import (
	"bytes"
	"github.com/giskook/smarthome-access/base"
	"github.com/giskook/smarthome-access/pb"
)

const CMD_WARN uint16 = 0x800d
const CMD_WARN_LEN uint16 = 0x1b

type Nsq_Warn_Packet struct {
	SerialNum       uint32
	DeviceID        uint64
	Endpoint        uint8
	WarningDuration uint16
	WarningMode     uint8
	Strobe          uint8
	SirenLevel      uint8
	StrobeLevel     uint8
	StrobeDutyCycle uint8
}

func (p *Nsq_Warn_Packet) Serialize() []byte {
	var writer bytes.Buffer
	writer.WriteByte(STARTFLAG)
	base.WriteWord(&writer, CMD_WARN_LEN)
	base.WriteWord(&writer, CMD_WARN)
	base.WriteDWord(&writer, p.SerialNum)
	base.WriteQuaWord(&writer, p.DeviceID)
	writer.WriteByte(p.Endpoint)
	base.WriteWord(&writer, p.WarningDuration)
	writer.WriteByte(p.WarningMode)
	writer.WriteByte(p.Strobe)
	writer.WriteByte(p.SirenLevel)
	writer.WriteByte(p.StrobeLevel)
	writer.WriteByte(p.StrobeDutyCycle)
	writer.WriteByte(CheckSum(writer.Bytes(), uint16(writer.Len())))
	writer.WriteByte(ENDFLAG)

	return writer.Bytes()
}

func Parse_NSQ_Warn(serialnum uint32, paras []*Report.Command_Param) *Nsq_Warn_Packet {
	deviceid := paras[0].Npara
	endpint := uint8(paras[1].Npara)
	warningduration := uint16(paras[2].Npara)
	warningmode := uint8(paras[3].Npara)
	strobe := uint8(paras[4].Npara)
	sirenlevel := uint8(paras[5].Npara)
	strobelevel := uint8(paras[6].Npara)
	strobedutycycle := uint8(paras[7].Npara)

	return &Nsq_Warn_Packet{
		SerialNum:       serialnum,
		DeviceID:        deviceid,
		Endpoint:        endpint,
		WarningDuration: warningduration,
		WarningMode:     warningmode,
		Strobe:          strobe,
		SirenLevel:      sirenlevel,
		StrobeLevel:     strobelevel,
		StrobeDutyCycle: strobedutycycle,
	}
}
