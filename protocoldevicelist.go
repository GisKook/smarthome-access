package sha

import (
	"encoding/binary"
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

type DeviceListPacket struct {
	GatewayID    uint64
	SerialNumber uint32
	DeviceCount  uint16
	DeviceList   []Device
}

func (p *DeviceListPacket) Serialize() []byte {
	para := []*Report.Command_Param{
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT16,
			Npara: uint64(p.DeviceCount),
		},
	}

	if p.DeviceCount > 0 {
		for i := uint16(0); i < p.DeviceCount; i++ {
			para = append(para, &Report.Command_Param{
				Type:  Report.Command_Param_UINT64,
				Npara: p.DeviceList[i].Oid,
			})
			para = append(para, &Report.Command_Param{
				Type:  Report.Command_Param_UINT8,
				Npara: uint64(p.DeviceList[i].Type),
			})
			para = append(para, &Report.Command_Param{
				Type:  Report.Command_Param_UINT16,
				Npara: uint64(p.DeviceList[i].Company),
			})

		}
	}
	command := &Report.Command{
		Type:  Report.Command_CMT_REPDEVICELIST,
		Paras: para,
	}
	devicelist := &Report.ControlReport{
		Tid:          p.GatewayID,
		SerialNumber: p.SerialNumber,
		Command:      command,
	}

	data, _ := proto.Marshal(devicelist)

	return data
}

func ParseDeviceList(buffer []byte, c *Conn) *DeviceListPacket {
	gatewayid, reader := GetGatewayID(buffer)
	serialnumber_byte := make([]byte, 4)
	reader.Read(serialnumber_byte)
	serialnumber := binary.BigEndian.Uint32(serialnumber_byte)

	devicecount_byte := make([]byte, 2)
	reader.Read(devicecount_byte)
	devicecount := binary.BigEndian.Uint16(devicecount_byte)

	var DeviceList []Device
	for i := uint16(0); i < devicecount; i++ {
		devicetype, _ := reader.ReadByte()
		deviceidlen, _ := reader.ReadByte()
		deviceid_byte := make([]byte, deviceidlen)
		reader.Read(deviceid_byte)
		deviceid := binary.BigEndian.Uint64(deviceid_byte)
		company_byte := make([]byte, 2)
		reader.Read(company_byte)
		company := binary.BigEndian.Uint16(company_byte)
		DeviceList = append(DeviceList, Device{
			Oid:     deviceid,
			Type:    devicetype,
			Company: company,
		})
	}

	return &DeviceListPacket{
		GatewayID:    gatewayid,
		SerialNumber: serialnumber,
		DeviceCount:  devicecount,
		DeviceList:   DeviceList,
	}
}
