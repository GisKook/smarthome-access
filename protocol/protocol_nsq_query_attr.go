package protocol

import (
	"bytes"
	"github.com/giskook/smarthome-access/base"
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

const CMD_QUERY_ATTR uint16 = 0x800b
const CMD_QUERY_ATTR_LEN uint16 = 0x14

type Nsq_Query_Attr_Packet struct {
	GatewayID uint64
	SerialNum uint32
	Status    byte
	DeviceID  uint64
	Endpoint  uint8
}

func (p *Nsq_Query_Attr_Packet) SerializeOnline() []byte {
	var writer bytes.Buffer
	writer.WriteByte(STARTFLAG)
	base.WriteWord(&writer, CMD_QUERY_ATTR_LEN)
	base.WriteWord(&writer, CMD_QUERY_ATTR)
	base.WriteDWord(&writer, p.SerialNum)
	base.WriteQuaWord(&writer, p.DeviceID)
	writer.WriteByte(p.Endpoint)
	writer.WriteByte(CheckSum(writer.Bytes(), uint16(writer.Len())))
	writer.WriteByte(ENDFLAG)

	return writer.Bytes()
}

func (p *Nsq_Query_Attr_Packet) SerializeOffline() []byte {
	para := []*Report.Command_Param{
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(p.Status),
		},
	}

	command := &Report.Command{
		Type:  Report.Command_CMT_REP_DEVICE_ATTR,
		Paras: para,
	}

	feedback_device_attr_pkg := &Report.ControlReport{
		Tid:          p.GatewayID,
		SerialNumber: p.SerialNum,
		Command:      command,
	}

	data, _ := proto.Marshal(feedback_device_attr_pkg)

	return data
}

func (p *Nsq_Query_Attr_Packet) Serialize() []byte {
	if p.Status == GATEWAY_ON_LINE {
		return p.SerializeOnline()
	} else {
		return p.SerializeOffline()
	}
}

func Parse_NSQ_Query_Attr(gatewayid uint64, serialnum uint32, status byte, paras []*Report.Command_Param) *Nsq_Query_Attr_Packet {
	deviceid := paras[0].Npara
	endpint := uint8(paras[1].Npara)

	return &Nsq_Query_Attr_Packet{
		GatewayID: gatewayid,
		SerialNum: serialnum,
		Status:    status,
		DeviceID:  deviceid,
		Endpoint:  endpint,
	}

}
