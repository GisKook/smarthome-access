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
			Npara: p.GatewayID,
		},
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT64,
			Npara: p.Device.ID,
		},
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: p.Action,
		},
	}
	if p.Action == ADD_DEVICE {
		endpoint_count := uint8(len(p.Device.Endpoints))
		Append(para, &Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: endpoint_count,
		})
		for i := uint8(0); i < endpoint_count; i++ {
			Append(para, &Report.Command_Param{
				Type:  Report.Command_Param_UINT8,
				Npara: p.Device.Endpoints[i].Endpoint,
			})
			Append(para, &Report.Command_Param{
				Type:  Report.Command_Param_UINT16,
				Npara: p.Device.Endpoints[i].DeviceTypeID,
			})
			if p.Device.Endpoints[i].DeviceTypeID == base.SS_Device_DeviceTypeID {
				Append(para, &Report.Command_Param{
					Type:  Report.Command_Param_UINT16,
					Npara: p.Device.Endpoints[i].Zonetype,
				})
			}

		}

	}

	command := &Report.Command{
		Type:  Report.Command_CMT_REP_ADD_DEL_DEVICE,
		Paras: para,
	}

	login := &Report.ControlReport{
		Tid:          p.GatewayID,
		SerialNumber: 0,
		Command:      command,
	}

	data, _ := proto.Marshal(login)

	return data
}

func Parse_Add_Del_Device(buffer []byte, uint64 id) *Add_Del_Device_Packet {
	gatewayid, reader := sha.GetGatewayID(buffer)

	action, _ := reader.ReadByte()
	deviceid := sha.ReadQuaWord(reader)
	device_type_count, _ := reader.ReadByte()
	var device base.Device
	device.ID = id
	for i := 0; uint16(i) < device_type_count; i++ {
		endpoint, _ := reader.ReadByte()
		devicetypeid := sha.ReadWord(reader)
		var zonetype uint16 = 0
		if devicetypeid == base.SS_Device_DeviceTypeID {
			zonetype = sha.ReadWord(reader)
		}
		Append(device.Endpoints, &Endpoint{
			Endpoint:     endpoint,
			DeviceTypeID: devicetypeid,
			Zonetype:     zonetype,
		})
	}

	return &Add_Del_Device_Packet{
		GatewayID: id,
		DeviceID:  deviceid,
		Action:    action,
		Device:    &device,
	}
}
