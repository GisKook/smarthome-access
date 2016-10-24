package protocol

import (
	"bytes"
	"encoding/binary"
	"github.com/giskook/smarthome-access/base"
	//"github.com/giskook/smarthome-access/pb"
	//"github.com/golang/protobuf/proto"
	"time"
)

const CMD_LOGIN uint16 = 0x8001
const CMD_LOGIN_LEN uint16 = 0x12

type LoginPacket struct {
	Gateway *base.Gateway
}

func (p *LoginPacket) Serialize() []byte {
	//	command := &Report.Command{
	//		Type:  Report.Command_CMT_REQ_LOGIN,
	//		Paras: nil,
	//	}
	//
	//	login := &Report.ControlReport{
	//		Tid:          p.Gateway.ID,
	//		SerialNumber: 0,
	//		Command:      command,
	//	}
	//
	//	data, _ := proto.Marshal(login)

	//      return data
	var writer bytes.Buffer
	writer.WriteByte(STARTFLAG)
	base.WriteWord(&writer, CMD_LOGIN_LEN)
	base.WriteWord(&writer, CMD_LOGIN)
	base.WriteMacBytes(&writer, p.Gateway.ID)
	t := uint32(time.Now().Unix())
	base.WriteDWord(&writer, t)
	writer.WriteByte(0)
	writer.WriteByte(CheckSum(writer.Bytes(), uint16(writer.Len())))
	writer.WriteByte(ENDFLAG)

	return writer.Bytes()
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
		device_name := base.ReadString(reader, device_name_len)
		devicelist[i].Name = device_name
		endpoint_count, _ := reader.ReadByte()
		endpoints := make([]base.Endpoint, endpoint_count)
		for j := 0; byte(j) < endpoint_count; j++ {
			endpoints[j].Endpoint, _ = reader.ReadByte()
			endpoints[j].DeviceTypeID = base.ReadWord(reader)
			if endpoints[j].DeviceTypeID == base.SS_Device_DeviceTypeID {
				endpoints[j].Zonetype = base.ReadWord(reader)
			} else if endpoints[j].DeviceTypeID == base.MPO_Device_DeviceTypeID || endpoints[j].DeviceTypeID == base.Shade_Device_DeviceTypeID || endpoints[j].DeviceTypeID == base.HA_Device_ON_OFF_Output_DeviceTypeID {
				endpoints[j].Status, _ = reader.ReadByte()
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
