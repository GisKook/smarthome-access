package sha

import (
	"bytes"
	"encoding/binary"
)

var (
	Infrared      uint8 = 0
	DoorMagnetic  uint8 = 1
	WarningButton uint8 = 2
)

type Device struct {
	Oid     []byte
	Type    uint8
	Company uint16
}

type LoginPacket struct {
	Uid             []byte
	BoxVersion      uint8
	ProtocolVersion uint8
	DeviceList      []Device

	Result uint8
}

func (this *LoginPacket) Serialize() []byte {
	var buf []byte
	buf = append(buf, 0xCE)
	buf = append(buf, 0xAA)
	buf = append(buf, 0xCE)

	return buf
}

func NewLoginPakcet(Uid []byte, BoxVersion uint8, ProtocolVersion uint8, DeviceList []Device) *LoginPacket {
	return &LoginPacket{
		Uid:             Uid,
		BoxVersion:      BoxVersion,
		ProtocolVersion: ProtocolVersion,
		DeviceList:      DeviceList,
		Result:          0,
	}
}

func ParseLogin(buffer []byte) *LoginPacket {
	reader := bytes.NewReader(buffer)
	reader.Seek(5, 0)
	uid := make([]byte, 6)
	reader.Read(uid)
	boxversion, _ := reader.ReadByte()
	protocolversion, _ := reader.ReadByte()
	devicecount_byte := make([]byte, 2)
	reader.Read(devicecount_byte)
	devicecount := binary.BigEndian.Uint16(devicecount_byte)
	devicelist := make([]Device, devicecount)
	for i := 0; uint16(i) < devicecount; i++ {
		devicelist[i].Type, _ = reader.ReadByte()
		deviceidlength, _ := reader.ReadByte()
		deviceid := make([]byte, deviceidlength)
		reader.Read(deviceid)
		devicelist[i].Oid = deviceid
		devicecompany_byte := make([]byte, 2)
		devicecompany := binary.BigEndian.Uint16(devicecompany_byte)
		devicelist[i].Company = devicecompany
	}

	return NewLoginPakcet(uid, boxversion, protocolversion, devicelist)

}
