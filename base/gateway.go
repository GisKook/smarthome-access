package base

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
