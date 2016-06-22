package protocol

import (
	"bytes"
	"github.com/giskook/smarthome-access/base"
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

const CMD_ONOFF uint16 = 0x8010
const CMD_ONOFF_LEN uint16 = 0x0015

type Nsq_OnOff_Packet struct {
	GatewayID uint64
	SerialNum uint32
	Status    byte
	DeviceID  uint64
	Endpoint  uint8
	Action    uint8
}

func (p *Nsq_OnOff_Packet) SerializeOnline() []byte {
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

func (p *Nsq_OnOff_Packet) SerializeOffline() []byte {
	para := []*Report.Command_Param{
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(p.Status),
		},
	}

	command := &Report.Command{
		Type:  Report.Command_CMT_REP_ONOFF,
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

func (p *Nsq_OnOff_Packet) Serialize() []byte {
	if p.Status == GATEWAY_ON_LINE {
		return p.SerializeOnline()
	} else {
		return p.SerializeOffline()
	}
}

func Parse_NSQ_OnOff(gatewayid uint64, serialnum uint32, status byte, paras []*Report.Command_Param) *Nsq_OnOff_Packet {
	deviceid := paras[0].Npara
	endpint := uint8(paras[1].Npara)
	action := uint8(paras[2].Npara)

	return &Nsq_OnOff_Packet{
		GatewayID: gatewayid,
		SerialNum: serialnum,
		Status:    status,
		DeviceID:  deviceid,
		Endpoint:  endpint,
		Action:    action,
	}
}
