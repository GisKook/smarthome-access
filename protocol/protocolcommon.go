package protocol

import (
	"bytes"
	"encoding/binary"
)

var (
	Illegal  uint16 = 0
	HalfPack uint16 = 255

	Login          uint16 = 1
	HeartBeat      uint16 = 2
	Add_Del_Device uint16 = 5
)

func GetGatewayID(buffer []byte) (uint64, *bytes.Reader) {
	reader := bytes.NewReader(buffer)
	reader.Seek(5, 0)
	uid := make([]byte, 6)
	reader.Read(uid)
	gid := []byte{0, 0}
	gid = append(gid, uid...)
	return binary.BigEndian.Uint64(gid), reader
}

func CheckSum(cmd []byte, cmdlen uint16) byte {
	temp := cmd[0]
	for i := uint16(1); i < cmdlen; i++ {
		temp ^= cmd[i]
	}

	return temp
}
func CheckProtocol(buffer *bytes.Buffer) (uint16, uint16) {
	//log.Printf("check protocol %x\n", buffer.Bytes())
	bufferlen := buffer.Len()
	if bufferlen == 0 {
		return Illegal, 0
	}
	if buffer.Bytes()[0] != 0xCE {
		buffer.ReadByte()
		CheckProtocol(buffer)
	} else if bufferlen > 2 {
		temp := make([]byte, 2)
		temp[0] = buffer.Bytes()[1]
		temp[1] = buffer.Bytes()[2]
		pkglen := binary.BigEndian.Uint16(temp)
		if pkglen < 8 || pkglen > 2048 { // flag + messagelen + cmdid + checksum + flag = 7  2048 is a magic number
			buffer.ReadByte()
			CheckProtocol(buffer)
		}
		if int(pkglen) > bufferlen {
			return HalfPack, 0
		} else {
			checksum := CheckSum(buffer.Bytes(), pkglen-2)
			if checksum == buffer.Bytes()[pkglen-2] && buffer.Bytes()[pkglen-1] == 0xCE {
				temp[0] = buffer.Bytes()[3]
				temp[1] = buffer.Bytes()[4]
				cmdid := binary.BigEndian.Uint16(temp)
				return cmdid, pkglen
			} else {
				buffer.ReadByte()
				CheckProtocol(buffer)
			}
		}
	} else {
		return HalfPack, 0
	}

	return Illegal, 0
}
