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
	}

	switch command.Type {
	case Report.Command_CMT_REQ_SETNAME:
		if c != nil {
			pkg := protocol.Parse_NSQ_Set_Device_Name(gatewayid, serialnum, protocol.GATEWAY_ON_LINE, command.Paras)
			c.SendToGateway(pkg)
		} else {
			pkg := protocol.Parse_NSQ_Set_Device_Name(gatewayid, serialnum, protocol.GATEWAY_OFF_LINE, command.Paras)
			GetServer().GetProducer().Send(GetConfiguration().NsqConfig.UpTopic, pkg.Serialize())
		}
	case Report.Command_CMT_REQ_DEL_DEVICE:
		if c != nil {
			pkg := protocol.Parse_NSQ_Del_Device(gatewayid, serialnum, protocol.GATEWAY_ON_LINE, command.Paras)
			c.SendToGateway(pkg)
		} else {
			pkg := protocol.Parse_NSQ_Del_Device(gatewayid, serialnum, protocol.GATEWAY_OFF_LINE, command.Paras)
			GetServer().GetProducer().Send(GetConfiguration().NsqConfig.UpTopic, pkg.Serialize())
		}
	case Report.Command_CMT_REQ_DEVICE_ATTR:
		if c != nil {
			pkg := protocol.Parse_NSQ_Query_Attr(gatewayid, serialnum, protocol.GATEWAY_ON_LINE, command.Paras)
			c.SendToGateway(pkg)
		} else {
			pkg := protocol.Parse_NSQ_Query_Attr(gatewayid, serialnum, protocol.GATEWAY_OFF_LINE, command.Paras)
			GetServer().GetProducer().Send(GetConfiguration().NsqConfig.UpTopic, pkg.Serialize())
		}
	case Report.Command_CMT_REQ_DEVICE_IDENTIFY:
		pkg := protocol.Parse_NSQ_Identify(serialnum, command.Paras)
		c.SendToGateway(pkg)
	case Report.Command_CMT_REQ_DEVICE_WARN:
		pkg := protocol.Parse_NSQ_Warn(serialnum, command.Paras)
		c.SendToGateway(pkg)
	case Report.Command_CMT_REQ_DEPLOYMENT:
		if c != nil {
			pkg := protocol.Parse_NSQ_Deployment(gatewayid, serialnum, protocol.GATEWAY_ON_LINE, command.Paras)
			c.SendToGateway(pkg)
		} else {
			pkg := protocol.Parse_NSQ_Deployment(gatewayid, serialnum, protocol.GATEWAY_OFF_LINE, command.Paras)
			GetServer().GetProducer().Send(GetConfiguration().NsqConfig.UpTopic, pkg.Serialize())
		}
	case Report.Command_CMT_REQ_DEVICE_LEVELCONTROL:
		pkg := protocol.Parse_NSQ_Level_Control(serialnum, command.Paras)
		c.SendToGateway(pkg)
	case Report.Command_CMT_REQ_ONOFF:
		if c != nil {
			pkg := protocol.Parse_NSQ_OnOff(gatewayid, serialnum, protocol.GATEWAY_ON_LINE, command.Paras)
			c.SendToGateway(pkg)
		} else {
			pkg := protocol.Parse_NSQ_OnOff(gatewayid, serialnum, protocol.GATEWAY_OFF_LINE, command.Paras)
			GetServer().GetProducer().Send(GetConfiguration().NsqConfig.UpTopic, pkg.Serialize())
		}
	case Report.Command_CMT_REQ_DEVICELIST:
		if c != nil {
			pkg := protocol.Parse_Device_List(gatewayid, serialnum, protocol.GATEWAY_ON_LINE, NewConns().GetConn(gatewayid).Gateway)
			GetServer().GetProducer().Send(GetConfiguration().NsqConfig.UpTopic, pkg.Serialize())
		} else {
			pkg := protocol.Parse_Device_List(gatewayid, serialnum, protocol.GATEWAY_OFF_LINE, nil)
			GetServer().GetProducer().Send(GetConfiguration().NsqConfig.UpTopic, pkg.Serialize())
		}
	}
}
