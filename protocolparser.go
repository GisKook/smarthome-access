package sha

import (
	"errors"
	"log"

	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

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
