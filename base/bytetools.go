package base

import (
	"bytes"
	"encoding/binary"
)

func ReadWord(reader *bytes.Reader) uint16 {
	word_byte := make([]byte, 2)
	reader.Read(word_byte)

	return binary.BigEndian.Uint16(word_byte)
}

func ReadDWord(reader *bytes.Reader) uint32 {
	dword_byte := make([]byte, 4)
	reader.Read(dword_byte)

	return binary.BigEndian.Uint32(dword_byte)
}

func ReadQuaWord(reader *bytes.Reader) uint64 {
	qword_byte := make([]byte, 8)
	reader.Read(qword_byte)

	return binary.BigEndian.Uint64(qword_byte)
}

func WriteMac(mac uint64) []byte {
	mac_byte := make([]byte, 8)
	binary.BigEndian.PutUint64(mac_byte, mac)

	return mac_byte[2:]
}
