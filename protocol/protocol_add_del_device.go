package protocol

import (
	"github.com/giskook/smarthome-access/base"
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

const ADD_DEVICE uint8 = 1
const DEL_DEVICE uint8 = 0

type Add_Del_Device_Packet struct {
	GatewayID uint64
	Action    uint8
	Device    *base.Device
}

func (p *Add_Del_Device_Packet) Serialize() []byte {
	para := []*Report.Command_Param{
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT64,
			Npara: p.Device.ID,
		},
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(p.Action),
		},
	}
	if p.Action == ADD_DEVICE {
		endpoint_count := uint8(len(p.Device.Endpoints))
		para = append(para, &Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(endpoint_count),
		})
		for i := uint8(0); i < endpoint_count; i++ {
			para = append(para, &Report.Command_Param{
				Type:  Report.Command_Param_UINT8,
				Npara: uint64(p.Device.Endpoints[i].Endpoint),
			})
			para = append(para, &Report.Command_Param{
				Type:  Report.Command_Param_UINT16,
				Npara: uint64(p.Device.Endpoints[i].DeviceTypeID),
			})
			if p.Device.Endpoints[i].DeviceTypeID == base.SS_Device_DeviceTypeID {
				para = append(para, &Report.Command_Param{
					Type:  Report.Command_Param_UINT16,
					Npara: uint64(p.Device.Endpoints[i].Zonetype),
				})
			}

		}

	}

	command := &Report.Command{
		Type:  Report.Command_CMT_REP_ADD_DEL_DEVICE,
		Paras: para,
	}

	add_del_device_pkg := &Report.ControlReport{
		Tid:          p.GatewayID,
		SerialNumber: 0,
		Command:      command,
	}

	data, _ := proto.Marshal(add_del_device_pkg)

	return data
}

func Parse_Add_Del_Device(buffer []byte, id uint64) *Add_Del_Device_Packet {
	reader := ParseHeader(buffer)
	action, _ := reader.ReadByte()
	deviceid := base.ReadQuaWord(reader)
	device_type_count, _ := reader.ReadByte()
	var device base.Device
	device.ID = deviceid
	for i := 0; byte(i) < device_type_count; i++ {
		endpoint, _ := reader.ReadByte()
		devicetypeid := base.ReadWord(reader)
		var zonetype uint16 = 0
		if devicetypeid == base.SS_Device_DeviceTypeID {
			zonetype = base.ReadWord(reader)
		}
		device.Endpoints = append(device.Endpoints, base.Endpoint{
			Endpoint:     endpoint,
			DeviceTypeID: devicetypeid,
			Zonetype:     zonetype,
		})
	}

	return &Add_Del_Device_Packet{
		GatewayID: id,
		Action:    action,
		Device:    &device,
	}
}
