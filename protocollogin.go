package sha

var (
	Infrared      uint8 = 0
	DoorMagnetic  uint8 = 1
	WarningButton uint8 = 2
)

//type Device struct {
//	Oid  []byte
//	Type uint8
//}

type LoginPacket struct {
	Uid             []byte
	BoxVersion      uint8
	ProtocolVersion uint8
	//DeviceCount     uint16
	//DeviceList      []Device

	Result uint8
}

func (this *LoginPacket) Serialize() []byte {
	var buf []byte
	buf = append(buf, 0xCE)
	buf = append(buf, 0xAA)
	buf = append(buf, 0xCE)

	return buf
}

func NewLoginPakcet(Uid []byte, BoxVersion uint8, ProtocolVersion uint8) *LoginPacket {
	return &LoginPacket{
		Uid:             Uid,
		BoxVersion:      BoxVersion,
		ProtocolVersion: ProtocolVersion,
		Result:          0,
	}
}
