package sha

import (
	"encoding/binary"
)

var (
	Infrared      uint8 = 0
	DoorMagnetic  uint8 = 1
	WarningButton uint8 = 2
)

type Device struct {
	Oid     uint64
	Type    uint8
	Company uint16
	Status  uint8
	Name    string
}

type LoginPacket struct {
	Uid             uint64
	BoxVersion      uint8
	ProtocolVersion uint8
	DeviceList      []Device

	Result uint8
}

func (p *LoginPacket) Serialize() []byte {
	var buf []byte
	buf = append(buf, 0xCE)
	buf = append(buf, 0x00)
	buf = append(buf, 0x0E)
	buf = append(buf, 0x80)
	buf = append(buf, 0x01)
	gatewayid := make([]byte, 8)
	binary.BigEndian.PutUint64(gatewayid, p.Uid)
	buf = append(buf, gatewayid[2:]...)
	buf = append(buf, p.Result)
	buf = append(buf, CheckSum(buf, uint16(len(buf))))
	buf = append(buf, 0xCE)

	return buf
}

func NewLoginPakcet(Uid uint64, BoxVersion uint8, ProtocolVersion uint8, DeviceList []Device, Result uint8) *LoginPacket {
	return &LoginPacket{
		Uid:             Uid,
		BoxVersion:      BoxVersion,
		ProtocolVersion: ProtocolVersion,
		DeviceList:      DeviceList,
		Result:          Result,
	}
}

func ParseLogin(buffer []byte, c *Conn) *LoginPacket {
	gatewayid, reader := GetGatewayID(buffer)
	ok := GetGatewayHub().Check(gatewayid)
	if ok {
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
			deid := []byte{0, 0}
			deid = append(deid, deviceid...)
			did := binary.BigEndian.Uint64(deid)
			devicelist[i].Oid = did
			devicecompany_byte := make([]byte, 2)
			reader.Read(devicecompany_byte)
			devicecompany := binary.BigEndian.Uint16(devicecompany_byte)
			devicelist[i].Company = devicecompany

			status, _ := reader.ReadByte()
			devicelist[i].Status = status
			devicenamelen, _ := reader.ReadByte()
			devicename := make([]byte, devicenamelen)
			reader.Read(devicename)
			devicelist[i].Name = string(devicename)
		}

		c.uid = gatewayid
		c.SetStatus(ConnSuccess)
		NewConns().Add(c)

		return NewLoginPakcet(gatewayid, boxversion, protocolversion, devicelist, 1)
	}

	return NewLoginPakcet(gatewayid, 0, 0, nil, 0)

}
