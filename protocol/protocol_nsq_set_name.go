package protocol

import (
	"bytes"
	"fmt"
	"github.com/giskook/smarthome-access/base"
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

const CMD_SET_DEVICE_NAME uint16 = 0x8008

type Nsq_Set_Devcie_Name_Packet struct {
	GatewayID uint64
	Status    byte
	SerialNum uint32
	DeviceID  uint64
	Name      string
}

func (p *Nsq_Set_Devcie_Name_Packet) SerializeOnline() []byte {
	fmt.Printf("%+v\n", p)
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

func (p *Nsq_Set_Devcie_Name_Packet) SerializeOffline() []byte {
	para := []*Report.Command_Param{
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(p.Status),
		},
	}

	command := &Report.Command{
		Type:  Report.Command_CMT_REP_SETNAME,
		Paras: para,
	}

	nsq_setname_pkg := &Report.ControlReport{
		Tid:          p.GatewayID,
		SerialNumber: p.SerialNum,
		Command:      command,
	}

	data, _ := proto.Marshal(nsq_setname_pkg)

	return data
}

func (p *Nsq_Set_Devcie_Name_Packet) Serialize() []byte {
	if p.Status == GATEWAY_OFF_LINE {
		return p.SerializeOffline()
	} else {
		return p.SerializeOnline()
	}
}

func Parse_NSQ_Set_Device_Name(gatewayid uint64, serialnum uint32, status byte, paras []*Report.Command_Param) *Nsq_Set_Devcie_Name_Packet {
	deviceid := paras[0].Npara
	fmt.Printf("%d\n", deviceid)
	name := paras[1].Strpara

	return &Nsq_Set_Devcie_Name_Packet{
		GatewayID: gatewayid,
		SerialNum: serialnum,
		Status:    status,
		DeviceID:  deviceid,
		Name:      name,
	}

}
