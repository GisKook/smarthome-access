package protocol

import (
	"encoding/binary"
	"github.com/giskook/smarthome-access/base"
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

type LoginPacket struct {
	Gateway *base.Gateway
}

func (p *LoginPacket) Serialize() []byte {
	command := &Report.Command{
		Type:  Report.Command_CMT_REQ_LOGIN,
		Paras: nil,
	}

	login := &Report.ControlReport{
		Tid:          p.Gateway.ID,
		SerialNumber: 0,
		Command:      command,
	}

	data, _ := proto.Marshal(login)

	return data
}

func ParseLogin(buffer []byte) *LoginPacket {
	reader := ParseHeader(buffer)
	gatewayid := base.ReadMac(reader)

	gatewaynamelen, _ := reader.ReadByte()
	gatewayname_byte := make([]byte, gatewaynamelen)
	reader.Read(gatewayname_byte)
	boxversion, _ := reader.ReadByte()
	protocolversion, _ := reader.ReadByte()
	devicecount_byte := make([]byte, 2)
	reader.Read(devicecount_byte)
	devicecount := binary.BigEndian.Uint16(devicecount_byte)
	devicelist := make([]base.Device, devicecount)
	for i := 0; uint16(i) < devicecount; i++ {
		devicelist[i].ID = base.ReadQuaWord(reader)
		device_name_len, _ := reader.ReadByte()
		device_name_byte := make([]byte, device_name_len)
		reader.Read(device_name_byte)
		endpoint_count, _ := reader.ReadByte()
		endpoints := make([]base.Endpoint, endpoint_count)
		for j := 0; byte(j) < endpoint_count; j++ {
			endpoints[j].Endpoint, _ = reader.ReadByte()
			endpoints[j].DeviceTypeID = base.ReadWord(reader)
			if endpoints[j].DeviceTypeID == 0x0402 {
				endpoints[j].Zonetype = base.ReadWord(reader)
			}
		}
		devicelist[i].Endpoints = endpoints
	}

	return &LoginPacket{
		Gateway: &base.Gateway{
			ID:              gatewayid,
			Name:            string(gatewayname_byte),
			BoxVersion:      boxversion,
			ProtocolVersion: protocolversion,
			Devices:         devicelist,
		},
	}

}
