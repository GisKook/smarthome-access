package base

const SS_Device_DeviceTypeID uint16 = 0x0402
const MPO_Device_DeviceTypeID uint16 = 0x0009 // Mains Power Outlet
const Shade_Device_DeviceTypeID uint16 = 0x0200
const HA_Device_ON_OFF_Output_DeviceTypeID uint16 = 0x0002

const ONLINE uint8 = 1
const OFFLINE uint8 = 0

type Endpoint struct {
	Endpoint     uint8
	DeviceTypeID uint16
	Zonetype     uint16
	Status       uint8
}

type Device struct {
	ID        uint64
	Name      string
	Status    uint8
	Endpoints []Endpoint
}

type Gateway struct {
	ID              uint64
	Name            string
	BoxVersion      byte
	ProtocolVersion byte
	Devices         []Device
}

func _gateway_get_device(gateway *Gateway, deviceid uint64) *Device {
	devicecount := len(gateway.Devices)
	var i int = 0
	for i = 0; i < devicecount; i++ {
		if gateway.Devices[i].ID == deviceid {
			return &gateway.Devices[i]
		}
	}

	return nil
}

func _device_get_endpoint(device *Device, endpoint uint8) *Endpoint {
	epcount := len(device.Endpoints)
	var i int = 0
	for i = 0; i < epcount; i++ {
		if device.Endpoints[i].Endpoint == endpoint {
			return &device.Endpoints[i]
		}
	}

	return nil
}

func Gateway_Add_Device(gateway *Gateway, device *Device) {
	gateway.Devices = append(gateway.Devices, *device)
}

func Gateway_Update_Davice(gateway *Gateway, device *Device) {
	Gateway_Del_Device(gateway, device.ID)
	Gateway_Add_Device(gateway, device)
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

func Gateway_Set_Device_Status(gateway *Gateway, deviceid uint64, endpoint uint8, status uint8) {
	device := _gateway_get_device(gateway, deviceid)
	if device != nil {
		ep := _device_get_endpoint(device, endpoint)
		if ep != nil {
			ep.Status = status
		}
	}
}

func Gateway_Set_Device_Online(gateway *Gateway, deviceid uint64, status uint8) {
	device := _gateway_get_device(gateway, deviceid)
	if device != nil {
		device.Status = status
	}
}
