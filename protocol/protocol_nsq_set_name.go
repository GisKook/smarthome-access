package protocol

import (
	"bytes"
	"github.com/giskook/smarthome-access/base"
	"github.com/giskook/smarthome-access/pb"
)

const CMD_SET_DEVICE_NAME uint16 = 0x8008

type Nsq_Set_Devcie_Name_Packet struct {
	SerialNum uint32
	DeviceID  uint64
	Name      string
}

func (p *Nsq_Set_Devcie_Name_Packet) Serialize() []byte {
	var writer bytes.Buffer
	writer.WriteByte(STARTFLAG)
	base.WriteWord(&writer, 0)
	base.WriteWord(&writer, CMD_SET_DEVICE_NAME)
	base.WriteDWord(&writer, p.SerialNum)
	base.WriteQuaWord(&writer, p.DeviceID)
	writer.WriteByte(byte(len(p.Name)))
	writer.WriteString(p.Name)
	base.WriteLength(&writer)
	writer.WriteByte(CheckSum(writer.Bytes(), uint16(writer.Len())))
	writer.WriteByte(ENDFLAG)

	return writer.Bytes()
}

func Parse_NSQ_Set_Device_Name(serialnum uint32, paras []*Report.Command_Param) *Nsq_Set_Devcie_Name_Packet {
	deviceid := paras[0].Npara
	name := paras[1].Strpara

	return &Nsq_Set_Devcie_Name_Packet{
		SerialNum: serialnum,
		DeviceID:  deviceid,
		Name:      name,
	}

}
