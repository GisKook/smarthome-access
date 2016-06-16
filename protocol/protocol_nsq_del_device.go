package protocol

import (
	"bytes"
	"github.com/giskook/smarthome-access/base"
	"github.com/giskook/smarthome-access/pb"
)

const CMD_DEL_DEVICE uint16 = 0x800a
const CMD_DEL_DEVICE_LEN uint16 = 0x13

type Nsq_Del_Devcie_Packet struct {
	SerialNum uint32
	DeviceID  uint64
}

func (p *Nsq_Del_Devcie_Packet) Serialize() []byte {
	var writer bytes.Buffer
	writer.WriteByte(STARTFLAG)
	base.WriteWord(&writer, CMD_DEL_DEVICE_LEN)
	base.WriteWord(&writer, CMD_DEL_DEVICE)
	base.WriteDWord(&writer, p.SerialNum)
	base.WriteQuaWord(&writer, p.DeviceID)
	writer.WriteByte(CheckSum(writer.Bytes(), uint16(writer.Len())))
	writer.WriteByte(ENDFLAG)

	return writer.Bytes()
}

func Parse_NSQ_Del_Device(serialnum uint32, paras []*Report.Command_Param) *Nsq_Del_Devcie_Packet {
	deviceid := paras[0].Npara

	return &Nsq_Del_Devcie_Packet{
		SerialNum: serialnum,
		DeviceID:  deviceid,
	}

}
