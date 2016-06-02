package protocol

import (
	"encoding/binary"
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

type LoginPacket struct {
	Gateway *sha.Gateway
}

func (p *LoginPacket) Serialize() []byte {
	para := []*Report.Command_Param{
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT64,
			Npara: p.Gateway.ID,
		},
	}

	command := &Report.Command{
		Type:  Report.Command_CMT_REQ_LOGIN,
		Paras: para,
	}

	login := &Report.ControlReport{
		Tid:          p.Gateway.ID,
		SerialNumber: 0,
		Command:      command,
	}

	data, _ := proto.Marshal(login)

	return data
}

func ParseLogin(buffer []byte, c *Conn) *LoginPacket {
	gatewayid, reader := sha.GetGatewayID(buffer)

	gatewaynamelen, _ := reader.ReadByte()
	gatewayname_byte := make([]byte, gatewaynamelen)
	reader.Read(gatewayname_byte)
	boxversion, _ := reader.ReadByte()
	protocolversion, _ := reader.ReadByte()
	devicecount_byte := make([]byte, 2)
	reader.Read(devicecount_byte)
	devicecount := binary.BigEndian.Uint16(devicecount_byte)
	devicelist := make([]sha.Device, devicecount)
	for i := 0; uint16(i) < devicecount; i++ {
		devicelist[i].ID = sha.ReadQuaWord(reader)
		device_name_len, _ = reader.ReadByte()
		device_name_byte := make([]byte, device_name_len)
		reader.Read(device_name_byte)
		endpoint_count := read.ReadByte()
		endpoints := make(&sha.Endpoint, endpoint_count)
		for j := 0; byte(j) < endpoint_count; j++ {
			endpoints[j].Endpoint, _ = read.ReadByte()
			endpoints[j].DeviceTypeID = sha.ReadWord()
			if endpoints[j].DeviceTypeID == 0x0402 {
				endpoints[j].Zonetype = sha.ReadWord()
			}
		}
		devicelist[i].Endpoints = endpoints
	}

	c.ID = gatewayid
	NewConns().Add(c)

	return &LoginPacket{
		Gateway: &sha.Gateway{
			ID:              gatewayid,
			Name:            string(gatewayname_byte),
			BoxVersion:      boxversion,
			ProtocolVersion: protocolversion,
			Devices:         DeviceList,
		},
	}

}
