package protocol

import (
	"github.com/giskook/smarthome-access/base"
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

const GATEWAY_ONLINE uint8 = 1
const GATEWAY_OFFLINE uint8 = 0

type Nsq_Device_List_Packet struct {
	Status  uint8
	Gateway *base.Gateway
}

func (p *Nsq_Device_List_Packet) Serialize() []byte {
	para := []*Report.Command_Param{
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(p.Status),
		},
	}
	if p.Status == GATEWAY_ONLINE {
		para = append(para, &Report.Command_Param{
			Type:    Report.Command_Param_STRING,
			Strpara: p.Gateway.Name,
		})
		device_count := len(p.Gateway.Devices)
		para = append(para, &Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(device_count),
		})
		for i := 0; i < device_count; i++ {
			para = append(para, &Report.Command_Param{
				Type:  Report.Command_Param_UINT64,
				Npara: p.Gateway.Devices[i].ID,
			})
			para = append(para, &Report.Command_Param{
				Type:    Report.Command_Param_STRING,
				Strpara: p.Gateway.Devices[i].Name,
			})
			endpoint_count := len(p.Gateway.Devices[i].Endpoints)
			para = append(para, &Report.Command_Param{
				Type:  Report.Command_Param_UINT8,
				Npara: uint64(endpoint_count),
			})
			for j := 0; j < endpoint_count; j++ {
				para = append(para, &Report.Command_Param{
					Type:  Report.Command_Param_UINT8,
					Npara: uint64(p.Gateway.Devices[i].Endpoints[j].Endpoint),
				})
				devicetypeid := p.Gateway.Devices[i].Endpoints[j].DeviceTypeID
				para = append(para, &Report.Command_Param{
					Type:  Report.Command_Param_UINT16,
					Npara: uint64(devicetypeid),
				})
				if devicetypeid == 0x0402 {
					para = append(para, &Report.Command_Param{
						Type:  Report.Command_Param_UINT16,
						Npara: uint64(p.Gateway.Devices[i].Endpoints[j].Zonetype),
					})
				}
			}
		}
	}
	command := &Report.Command{
		Type:  Report.Command_CMT_REP_DEVICELIST,
		Paras: para,
	}

	device_list_pkg := &Report.ControlReport{
		Tid:          p.Gateway.ID,
		SerialNumber: 0,
		Command:      command,
	}

	data, _ := proto.Marshal(device_list_pkg)

	return data
}

func Parse_Device_List(status uint8, gateway *base.Gateway) *Nsq_Device_List_Packet {
	return &Nsq_Device_List_Packet{
		Status:  status,
		Gateway: gateway,
	}
}
