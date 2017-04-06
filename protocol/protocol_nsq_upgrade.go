package protocol

import (
	"bytes"
	"github.com/giskook/smarthome-access/base"
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

const CMD_UPGRADE uint16 = 0x8022

type Nsq_Upgrade_Packet struct {
	GatewayID uint64
	SerialNum uint32
	Status    uint8
}

func (p *Nsq_Upgrade_Packet) SerializeOnline() []byte {
	var writer bytes.Buffer
	writer.WriteByte(STARTFLAG)
	base.WriteWord(&writer, 0)
	base.WriteWord(&writer, CMD_UPGRADE)
	base.WriteDWord(&writer, p.SerialNum)
	base.WriteQuaWord(&writer, p.GatewayID)
	writer.WriteByte(0)
	base.WriteLength(&writer)
	writer.WriteByte(CheckSum(writer.Bytes(), uint16(writer.Len())))
	writer.WriteByte(ENDFLAG)

	return writer.Bytes()
}

func (p *Nsq_Upgrade_Packet) SerializeOffline() []byte {
	para := []*Report.Command_Param{
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(GATEWAY_OFF_LINE),
		},
	}

	command := &Report.Command{
		Type:  Report.Command_CMT_REQ_UPGRADE,
		Paras: para,
	}

	feedback_onoff_pkg := &Report.ControlReport{
		Tid:          p.GatewayID,
		SerialNumber: p.SerialNum,
		Command:      command,
	}

	data, _ := proto.Marshal(feedback_onoff_pkg)

	return data
}

func (p *Nsq_Upgrade_Packet) Serialize() []byte {
	if p.Status == GATEWAY_ON_LINE {
		return p.SerializeOnline()
	} else {
		return p.SerializeOffline()
	}
}

func Parse_NSQ_Upgrade(gatewayid uint64, serialnum uint32, status uint8) *Nsq_Upgrade_Packet {
	return &Nsq_Upgrade_Packet{
		GatewayID: gatewayid,
		SerialNum: serialnum,
		Status:    status,
	}
}
