package protocol

//import (
//	"github.com/giskook/smarthome-access/pb"
//	"github.com/golang/protobuf/proto"
//)
//
//type DeviceTypeID struct {
//	Endpoint     uint8
//	DeviceTypeID uint16
//}
//
//type Add_Del_Device_Packet struct {
//	DeviceID      uint64
//	Action        uint8
//	DeviceTypeIDs []DeviceTypeID
//}
//
//func (p *Add_Del_Device_Packet) Serialize() []byte {
//	para := []*Report.Command_Param{
//		&Report.Command_Param{
//			Type:  Report.Command_Param_UINT64,
//			Npara: p.Gateway.ID,
//		},
//	}
//
//	command := &Report.Command{
//		Type:  Report.Command_CMT_REQ_LOGIN,
//		Paras: para,
//	}
//
//	login := &Report.ControlReport{
//		Tid:          p.Gateway.ID,
//		SerialNumber: 0,
//		Command:      command,
//	}
//
//	data, _ := proto.Marshal(login)
//
//	return data
//}
//
//func Parse_Add_Del_Device(buffer []byte) *Add_Del_Device_Packet {
//	gatewayid, reader := sha.GetGatewayID(buffer)
//
//	action, _ := reader.ReadByte()
//	deviceid := sha.ReadQuaWord(reader)
//	device_type_count, _ := reader.ReadByte()
//	device_type_ids := make(DeviceTypeID, device_type_count)
//	for i := 0; uint16(i) < device_type_count; i++ {
//		device_type_ids[i].Endpoint, _ = reader.ReadByte()
//		device_type_ids[i].DeviceTypeID = sha.ReadWord(reader)
//	}
//
//	return &Add_Del_Device_Packet{
//		DeviceID:      deviceid,
//		Action:        action,
//		DeviceTypeIDs: device_type_ids,
//	}
//}
