package protocol

import (
	"bytes"
	"encoding/binary"

	"errors"
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
	"log"
)

const STARTFLAG byte = 0xCE
const ENDFLAG byte = 0xCE
const GATEWAY_OFF_LINE byte = 252
const GATEWAY_ON_LINE byte = 251

var (
	Illegal  uint16 = 0
	HalfPack uint16 = 255

	Login                  uint16 = 1
	HeartBeat              uint16 = 2
	Add_Del_Device         uint16 = 5
	Notification           uint16 = 6
	Feedback_SetName       uint16 = 8
	Feedback_Del_Device    uint16 = 10
	Feedback_Query_Attr    uint16 = 11
	Feedback_Depolyment    uint16 = 15
	Feedback_OnOff         uint16 = 16
	Feedback_Warn          uint16 = 0x000d
	Notify_OnOff           uint16 = 0x0017
	Feedback_Level_Control uint16 = 0x0011
	Notify_Level           uint16 = 0x0019
	Notify_Online_Status   uint16 = 0x0020
	Feedback_Upgrade       uint16 = 0x0022
	Notify_Temperature     uint16 = 0x0024
	Notify_Humidity        uint16 = 0x0025
	Notify_Security_Aids   uint16 = 0x0030
)

func ParseHeader(buffer []byte) *bytes.Reader {
	reader := bytes.NewReader(buffer)
	reader.Seek(5, 0)

	return reader
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

	return HalfPack, 0
}

func CheckNsqProtocol(message []byte) (uint64, uint32, *Report.Command, error) {
	command := &Report.ControlReport{}
	err := proto.Unmarshal(message, command)
	if err != nil {
		log.Println("unmarshal error")
		return 0, 0, nil, errors.New("unmarshal error")
	} else {
		gatewayid := command.Tid
		serialnum := command.SerialNumber
		cmd := command.GetCommand()

		return gatewayid, serialnum, cmd, nil
	}
}
