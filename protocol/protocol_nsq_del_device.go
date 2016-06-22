package protocol

import (
	"bytes"
	"github.com/giskook/smarthome-access/base"
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

const CMD_DEL_DEVICE uint16 = 0x800a
const CMD_DEL_DEVICE_LEN uint16 = 0x13

type Nsq_Del_Devcie_Packet struct {
	GatewayID uint64
	Status    byte
	SerialNum uint32
	DeviceID  uint64
}

func (p *Nsq_Del_Devcie_Packet) SerializeOnline() []byte {
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

func (p *Nsq_Del_Devcie_Packet) SerializeOffline() []byte {
	para := []*Report.Command_Param{
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(p.Status),
		},
	}

	command := &Report.Command{
		Type:  Report.Command_CMT_REP_DEL_DEVICE,
		Paras: para,
	}

	feedback_del_device_pkg := &Report.ControlReport{
		Tid:          p.GatewayID,
		SerialNumber: p.SerialNum,
		Command:      command,
	}

	data, _ := proto.Marshal(feedback_del_device_pkg)

	return data
}

func (p *Nsq_Del_Devcie_Packet) Serialize() []byte {
	if p.Status == GATEWAY_ON_LINE {
		return p.SerializeOnline()
	} else {
		return p.SerializeOffline()
	}
}

func Parse_NSQ_Del_Device(gatewayid uint64, serialnum uint32, status byte, paras []*Report.Command_Param) *Nsq_Del_Devcie_Packet {
	deviceid := paras[0].Npara

	return &Nsq_Del_Devcie_Packet{
		GatewayID: gatewayid,
		SerialNum: serialnum,
		Status:    status,
		DeviceID:  deviceid,
	}

}
