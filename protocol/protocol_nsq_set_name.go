package protocol

import (
	"github.com/giskook/smarthome-access/pb"
)

type Nsq_Set_Devcie_Name_Packet struct {
	SerailNum uint32
	DeviceID  uint64
	Name      string
}

func Parse_NSQ_Set_Device_Name(serialnum uint32,paras Report.Command.Paras) *Nsq_Set_Devcie_Name_Packet { 
	deviceid := paras[0].Npara
	name := paras[1].Strpara

	return & Nsq_Set_Devcie_Name_Packet{
		SerailNum:serialnum,
		DeviceID:deviceid,
		Name:name,
	}

}
