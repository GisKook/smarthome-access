package protocol

import (
	"github.com/giskook/smarthome-access/base"
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

type Notify_Upgrade_Status_Packet struct {
	GatewayID       uint64
	SerialID        uint32
	ProtocolVersion uint8
	HardwareVersion uint8
	Result          uint8
}

func (p *Notify_Upgrade_Status_Packet) Serialize() []byte {
	para := []*Report.Command_Param{
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(p.ProtocolVersion),
		},
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(p.HardwareVersion),
		},
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(p.Result),
		},
	}

	command := &Report.Command{
		Type:  Report.Command_CMT_REP_NOTIFY_UPGRADE_STATUS,
		Paras: para,
	}

	notify_temperature := &Report.ControlReport{
		Tid:          p.GatewayID,
		SerialNumber: p.SerialID,
		Command:      command,
	}

	data, _ := proto.Marshal(notify_temperature)

	return data
}

func Parse_Notify_Upgrade_Status(buffer []byte) *Notify_Upgrade_Status_Packet {
	reader := ParseHeader(buffer)
	serialid := base.ReadDWord(reader)
	gateway_id := base.ReadMac(reader)
	protocol_version, _ := reader.ReadByte()
	hardware_version, _ := reader.ReadByte()
	result, _ := reader.ReadByte()

	return &Notify_Upgrade_Status_Packet{
		GatewayID:       gateway_id,
		SerialID:        serialid,
		ProtocolVersion: protocol_version,
		HardwareVersion: hardware_version,
		Result:          result,
	}
}
