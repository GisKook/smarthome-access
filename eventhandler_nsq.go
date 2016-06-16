package sha

import (
	"github.com/giskook/smarthome-access/pb"
	"github.com/giskook/smarthome-access/protocol"
)

func Nsq_EventHandler(gatewayid uint64, serialnum uint32, command *Report.Command) {
	switch command.Type {
	case Report.Command_CMT_REQ_SETNAME:
		pkg := protocol.Parse_NSQ_Set_Device_Name(serialnum, command.Paras)
		NewConns().GetConn(gatewayid).SendToGateway(pkg)
	case Report.Command_CMT_REQ_DEL_DEVICE:
		pkg := protocol.Parse_NSQ_Del_Device(serialnum, command.Paras)
		NewConns().GetConn(gatewayid).SendToGateway(pkg)
	case Report.Command_CMT_REQ_DEVICE_ATTR:
		pkg := protocol.Parse_NSQ_Query_Attr(serialnum, command.Paras)
		NewConns().GetConn(gatewayid).SendToGateway(pkg)
	case Report.Command_CMT_REQ_DEVICE_IDENTIFY:
		pkg := protocol.Parse_NSQ_Identify(serialnum, command.Paras)
		NewConns().GetConn(gatewayid).SendToGateway(pkg)
	case Report.Command_CMT_REQ_DEVICE_WARN:
		pkg := protocol.Parse_NSQ_Warn(serialnum, command.Paras)
		NewConns().GetConn(gatewayid).SendToGateway(pkg)
	case Report.Command_CMT_REQ_DEPLOYMENT:
		pkg := protocol.Parse_NSQ_Deployment(serialnum, command.Paras)
		NewConns().GetConn(gatewayid).SendToGateway(pkg)
	case Report.Command_CMT_REQ_DEVICE_LEVELCONTROL:
		pkg := protocol.Parse_NSQ_Level_Control(serialnum, command.Paras)
		NewConns().GetConn(gatewayid).SendToGateway(pkg)
	case Report.Command_CMT_REQ_ONOFF:
		pkg := protocol.Parse_NSQ_OnOff(serialnum, command.Paras)
		NewConns().GetConn(gatewayid).SendToGateway(pkg)
	case Report.Command_CMT_REQ_DEVICELIST:
		if NewConns().Check(gatewayid) {
			pkg := protocol.Parse_Device_List(protocol.GATEWAY_ONLINE, NewConns().GetConn(gatewayid).Gateway)
			GetServer().GetProducer().Send(NSQ_CONTROL_PUB_TOPIC, pkg.Serialize())
		} else {
			pkg := protocol.Parse_Device_List(protocol.GATEWAY_OFFLINE, nil)
			GetServer().GetProducer().Send(NSQ_CONTROL_PUB_TOPIC, pkg.Serialize())
		}
	}
}
