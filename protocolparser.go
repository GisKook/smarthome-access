package sha

import (
	"bytes"
	"encoding/binary"
)

func CheckProtocol(buffer *bytes.Buffer) uint8 {
	bufferlen := buffer.Len()
	if bufferlen == 0 {
		return Illegal
	}
	if buffer.Bytes()[0] != 0xCE {
		buffer.ReadByte()
		CheckProtocol(buffer)
	} else if bufferlen > 2 {
		var temp [2]byte
		temp[0] = buffer.Bytes()[1]
		temp[1] = buffer.Bytes()[2]
		pkglen := binary.BigEndian.Uint16(temp)
		if pkglen > bufferlen {
			return HalfPack
		} else {
			temp[0] = buffer.Bytes()[3]
			temp[1] = buffer.Bytes()[4]
			cmdid := binary.BigEndian.Uint16(temp)

			return cmdid
		}
	} else {
		return HalfPack
	}
}
