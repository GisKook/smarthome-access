package protocol

import (
	"bytes"
	"github.com/giskook/smarthome-access/base"
	"github.com/giskook/smarthome-access/pb"
)

const CMD_ONOFF uint16 = 0x8010
const CMD_ONOFF_LEN uint16 = 0x0015

type Nsq_OnOff_Packet struct {
	SerialNum uint32
	DeviceID  uint64
	Endpoint  uint8
	Action    uint8
}

func (p *Nsq_OnOff_Packet) Serialize() []byte {
	var writer bytes.Buffer
	writer.WriteByte(STARTFLAG)
	base.WriteWord(&writer, CMD_ONOFF_LEN)
	base.WriteWord(&writer, CMD_ONOFF)
	base.WriteDWord(&writer, p.SerialNum)
	base.WriteQuaWord(&writer, p.DeviceID)
	writer.WriteByte(byte(p.Endpoint))
	writer.WriteByte(byte(p.Action))
	writer.WriteByte(CheckSum(writer.Bytes(), uint16(writer.Len())))
	writer.WriteByte(ENDFLAG)

	return writer.Bytes()
}

func Parse_NSQ_OnOff(serialnum uint32, paras []*Report.Command_Param) *Nsq_OnOff_Packet {
	deviceid := paras[0].Npara
	endpint := uint8(paras[1].Npara)
	action := uint8(paras[2].Npara)

	return &Nsq_OnOff_Packet{
		SerialNum: serialnum,
		DeviceID:  deviceid,
		Endpoint:  endpint,
		Action:    action,
	}
}
