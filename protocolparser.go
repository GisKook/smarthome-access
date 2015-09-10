package sha

import (
	"bytes"
	"encoding/binary"
)

func CheckSum(cmd []byte, cmdlen uint16) byte {
	temp := cmd[0]
	for i := 1; i < cmdlen; i++ {
		temp ^= cmd[i]
	}

	return temp
}

func CheckProtocol(buffer *bytes.Buffer) (uint8, uint16) {
	bufferlen := buffer.Len()
	if bufferlen == 0 {
		return Illegal, 0
	}
	if buffer.Bytes()[0] != 0xCE {
		buffer.ReadByte()
		CheckProtocol(buffer)
	} else if bufferlen > 2 {
		var temp [2]byte
		temp[0] = buffer.Bytes()[1]
		temp[1] = buffer.Bytes()[2]
		pkglen := binary.BigEndian.Uint16(temp)
		if pkglen < 8 || pkglen > 2048 { // flag + messagelen + cmdid + checksum + flag = 7  2048 is a magic number
			buffer.ReadByte()
			CheckProtocol(buffer)
		}
		if pkglen > bufferlen {
			return HalfPack, 0
		} else {
			if CheckSum(buffer.Bytes()[0], pkglen-2) == buffer.Bytes()[pkglen-2] && buffer.Bytes()[pkglen-1] == 0xCE {
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
}
