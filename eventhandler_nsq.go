package sha

import (
	"github.com/giskook/smarthome-access/pb"
	"github.com/giskook/smarthome-access/protocol"
	"log"
)

func Nsq_EventHandler(gatewayid uint64, serialnum uint32, command *Report.Command) {
	c := NewConns().GetConn(gatewayid)
	if c == nil {
		log.Printf("<INFO>   can not find gateway %x \n", gatewayid)
		pkg := protocol.Feedback_General_Offline_Package{
			GatewayID: gatewayid,
			SerialNum: serialnum,
		}
		GetServer().GetProducer().Send(GetConfiguration().NsqConfig.UpTopic, pkg.Serialize())
		return
	}

	switch command.Type {
	case Report.Command_CMT_REQ_SETNAME:
		pkg := protocol.Parse_NSQ_Set_Device_Name(serialnum, command.Paras)
		c.SendToGateway(pkg)
	case Report.Command_CMT_REQ_DEL_DEVICE:
		pkg := protocol.Parse_NSQ_Del_Device(serialnum, command.Paras)
		c.SendToGateway(pkg)
	case Report.Command_CMT_REQ_DEVICE_ATTR:
		pkg := protocol.Parse_NSQ_Query_Attr(serialnum, command.Paras)
		c.SendToGateway(pkg)
	case Report.Command_CMT_REQ_DEVICE_IDENTIFY:
		pkg := protocol.Parse_NSQ_Identify(serialnum, command.Paras)
		c.SendToGateway(pkg)
	case Report.Command_CMT_REQ_DEVICE_WARN:
		pkg := protocol.Parse_NSQ_Warn(serialnum, command.Paras)
		c.SendToGateway(pkg)
	case Report.Command_CMT_REQ_DEPLOYMENT:
		pkg := protocol.Parse_NSQ_Deployment(serialnum, command.Paras)
		c.SendToGateway(pkg)
	case Report.Command_CMT_REQ_DEVICE_LEVELCONTROL:
		pkg := protocol.Parse_NSQ_Level_Control(serialnum, command.Paras)
		c.SendToGateway(pkg)
	case Report.Command_CMT_REQ_ONOFF:
		pkg := protocol.Parse_NSQ_OnOff(serialnum, command.Paras)
		c.SendToGateway(pkg)
	case Report.Command_CMT_REQ_DEVICELIST:
		pkg := protocol.Parse_Device_List(protocol.GATEWAY_ONLINE, NewConns().GetConn(gatewayid).Gateway)
		GetServer().GetProducer().Send(GetConfiguration().NsqConfig.UpTopic, pkg.Serialize())
	}
}
