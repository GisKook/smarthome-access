package base

const SS_Device_DeviceTypeID uint16 = 0x0402

type Endpoint struct {
	Endpoint     uint8
	DeviceTypeID uint16
	Zonetype     uint16
}

type Device struct {
	ID        uint64
	Name      string
	Endpoints []Endpoint
}

type Gateway struct {
	ID              uint64
	Name            string
	BoxVersion      byte
	ProtocolVersion byte
	Devices         []Device
}

func Gateway_Add_Device(gateway *Gateway, device *Device) {
	gateway.Devices = append(gateway.Devices, *device)
}

func Gateway_Del_Device(gateway *Gateway, deviceid uint64) {
	devicecount := len(gateway.Devices)
	var i int = 0
	for i = 0; i < devicecount; i++ {
		if gateway.Devices[i].ID == deviceid {
			break
		}
	}

	gateway.Devices = append(gateway.Devices[:i], gateway.Devices[i+1:]...)
}

func Gateway_Set_Device_Name(gateway *Gateway, deviceid uint64, name string) {
	device_count := len(gateway.Devices)
	for i := 0; i < device_count; i++ {
		if deviceid == gateway.Devices[i].ID {
			gateway.Devices[i].Name = name
		}
	}
}

func Gateway_Check_Device(gateway *Gateway, deviceid uint64) bool {
	device_count := len(gateway.Devices)
	for i := 0; i < device_count; i++ {
		if deviceid == gateway.Devices[i].ID {
			return true
		}
	}

	return false
}
