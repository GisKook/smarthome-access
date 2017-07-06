// Code generated by protoc-gen-go.
// source: Command.proto
// DO NOT EDIT!

/*
Package Report is a generated protocol buffer package.

It is generated from these files:
	Command.proto

It has these top-level messages:
	Command
	ControlReport
*/
package Report

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type Command_CommandType int32

const (
	Command_CMT_INVALID Command_CommandType = 0
	// gateway->web message
	Command_CMT_REQ_LOGIN                  Command_CommandType = 1
	Command_CMT_REP_ONLINE                 Command_CommandType = 2
	Command_CMT_REP_ADD_DEL_DEVICE         Command_CommandType = 5
	Command_CMT_REP_NOTIFICATION           Command_CommandType = 6
	Command_CMT_REP_SETNAME                Command_CommandType = 8
	Command_CMT_REP_DEL_DEVICE             Command_CommandType = 10
	Command_CMT_REP_DEVICE_ATTR            Command_CommandType = 11
	Command_CMT_REP_DEPLOYMENT             Command_CommandType = 15
	Command_CMT_REP_ONOFF                  Command_CommandType = 19
	Command_CMT_REP_DEVICE_ONLINE          Command_CommandType = 21
	Command_CMT_REP_NOTIFY_ONOFF           Command_CommandType = 23
	Command_CMT_REP_DEVICELIST             Command_CommandType = 2049
	Command_CMT_REP_NOTIFY_LEVEL           Command_CommandType = 25
	Command_CMT_REP_FEEDBACK_UPGRADE       Command_CommandType = 34
	Command_CMT_REP_NOTIFY_GATEWAY         Command_CommandType = 257
	Command_CMT_REP_NOTIFY_TEMPERATURE     Command_CommandType = 36
	Command_CMT_REP_NOTIFY_HUMIDITY        Command_CommandType = 37
	Command_CMT_REP_NOTIFY_SECURITY_AIDS   Command_CommandType = 48
	Command_CMT_REP_NOTIFY_ONLINE_STATUS   Command_CommandType = 32
	Command_CMT_REP_FEEDBACK_LEVEL         Command_CommandType = 17
	Command_CMT_REP_FEEDBACK_WARN          Command_CommandType = 13
	Command_CMT_REP_NOTIFY_GATEWAY_OFFLINE Command_CommandType = 258
	Command_CMT_REP_NOTIFY_UPGRADE_STATUS  Command_CommandType = 63
	// web->gateway
	Command_CMT_REP_LOGIN                  Command_CommandType = 32769
	Command_CMT_REQ_ONLINE                 Command_CommandType = 32770
	Command_CMT_REQ_SETNAME                Command_CommandType = 32776
	Command_CMT_REQ_DEL_DEVICE             Command_CommandType = 32778
	Command_CMT_REQ_DEVICE_ATTR            Command_CommandType = 32779
	Command_CMT_REQ_DEVICE_IDENTIFY        Command_CommandType = 32780
	Command_CMT_REQ_DEVICE_WARN            Command_CommandType = 32781
	Command_CMT_REQ_DEPLOYMENT             Command_CommandType = 32783
	Command_CMT_REQ_DEVICE_LEVELCONTROL    Command_CommandType = 32785
	Command_CMT_REQ_ONOFF                  Command_CommandType = 32787
	Command_CMT_REQ_DEVICE_ONLINE          Command_CommandType = 32789
	Command_CMT_REQ_ONOFF_STATUS           Command_CommandType = 32790
	Command_CMT_REQ_DEVICELIST             Command_CommandType = 34817
	Command_CMT_REQ_READ_DEPLOYMENT_STATUS Command_CommandType = 32799
	Command_CMT_REQ_UPGRADE                Command_CommandType = 32802
)

var Command_CommandType_name = map[int32]string{
	0:     "CMT_INVALID",
	1:     "CMT_REQ_LOGIN",
	2:     "CMT_REP_ONLINE",
	5:     "CMT_REP_ADD_DEL_DEVICE",
	6:     "CMT_REP_NOTIFICATION",
	8:     "CMT_REP_SETNAME",
	10:    "CMT_REP_DEL_DEVICE",
	11:    "CMT_REP_DEVICE_ATTR",
	15:    "CMT_REP_DEPLOYMENT",
	19:    "CMT_REP_ONOFF",
	21:    "CMT_REP_DEVICE_ONLINE",
	23:    "CMT_REP_NOTIFY_ONOFF",
	2049:  "CMT_REP_DEVICELIST",
	25:    "CMT_REP_NOTIFY_LEVEL",
	34:    "CMT_REP_FEEDBACK_UPGRADE",
	257:   "CMT_REP_NOTIFY_GATEWAY",
	36:    "CMT_REP_NOTIFY_TEMPERATURE",
	37:    "CMT_REP_NOTIFY_HUMIDITY",
	48:    "CMT_REP_NOTIFY_SECURITY_AIDS",
	32:    "CMT_REP_NOTIFY_ONLINE_STATUS",
	17:    "CMT_REP_FEEDBACK_LEVEL",
	13:    "CMT_REP_FEEDBACK_WARN",
	258:   "CMT_REP_NOTIFY_GATEWAY_OFFLINE",
	63:    "CMT_REP_NOTIFY_UPGRADE_STATUS",
	32769: "CMT_REP_LOGIN",
	32770: "CMT_REQ_ONLINE",
	32776: "CMT_REQ_SETNAME",
	32778: "CMT_REQ_DEL_DEVICE",
	32779: "CMT_REQ_DEVICE_ATTR",
	32780: "CMT_REQ_DEVICE_IDENTIFY",
	32781: "CMT_REQ_DEVICE_WARN",
	32783: "CMT_REQ_DEPLOYMENT",
	32785: "CMT_REQ_DEVICE_LEVELCONTROL",
	32787: "CMT_REQ_ONOFF",
	32789: "CMT_REQ_DEVICE_ONLINE",
	32790: "CMT_REQ_ONOFF_STATUS",
	34817: "CMT_REQ_DEVICELIST",
	32799: "CMT_REQ_READ_DEPLOYMENT_STATUS",
	32802: "CMT_REQ_UPGRADE",
}
var Command_CommandType_value = map[string]int32{
	"CMT_INVALID":                    0,
	"CMT_REQ_LOGIN":                  1,
	"CMT_REP_ONLINE":                 2,
	"CMT_REP_ADD_DEL_DEVICE":         5,
	"CMT_REP_NOTIFICATION":           6,
	"CMT_REP_SETNAME":                8,
	"CMT_REP_DEL_DEVICE":             10,
	"CMT_REP_DEVICE_ATTR":            11,
	"CMT_REP_DEPLOYMENT":             15,
	"CMT_REP_ONOFF":                  19,
	"CMT_REP_DEVICE_ONLINE":          21,
	"CMT_REP_NOTIFY_ONOFF":           23,
	"CMT_REP_DEVICELIST":             2049,
	"CMT_REP_NOTIFY_LEVEL":           25,
	"CMT_REP_FEEDBACK_UPGRADE":       34,
	"CMT_REP_NOTIFY_GATEWAY":         257,
	"CMT_REP_NOTIFY_TEMPERATURE":     36,
	"CMT_REP_NOTIFY_HUMIDITY":        37,
	"CMT_REP_NOTIFY_SECURITY_AIDS":   48,
	"CMT_REP_NOTIFY_ONLINE_STATUS":   32,
	"CMT_REP_FEEDBACK_LEVEL":         17,
	"CMT_REP_FEEDBACK_WARN":          13,
	"CMT_REP_NOTIFY_GATEWAY_OFFLINE": 258,
	"CMT_REP_NOTIFY_UPGRADE_STATUS":  63,
	"CMT_REP_LOGIN":                  32769,
	"CMT_REQ_ONLINE":                 32770,
	"CMT_REQ_SETNAME":                32776,
	"CMT_REQ_DEL_DEVICE":             32778,
	"CMT_REQ_DEVICE_ATTR":            32779,
	"CMT_REQ_DEVICE_IDENTIFY":        32780,
	"CMT_REQ_DEVICE_WARN":            32781,
	"CMT_REQ_DEPLOYMENT":             32783,
	"CMT_REQ_DEVICE_LEVELCONTROL":    32785,
	"CMT_REQ_ONOFF":                  32787,
	"CMT_REQ_DEVICE_ONLINE":          32789,
	"CMT_REQ_ONOFF_STATUS":           32790,
	"CMT_REQ_DEVICELIST":             34817,
	"CMT_REQ_READ_DEPLOYMENT_STATUS": 32799,
	"CMT_REQ_UPGRADE":                32802,
}

func (x Command_CommandType) String() string {
	return proto.EnumName(Command_CommandType_name, int32(x))
}

type Command_Param_ParaType int32

const (
	Command_Param_Null   Command_Param_ParaType = 0
	Command_Param_UINT8  Command_Param_ParaType = 1
	Command_Param_UINT16 Command_Param_ParaType = 2
	Command_Param_UINT32 Command_Param_ParaType = 3
	Command_Param_UINT64 Command_Param_ParaType = 4
	Command_Param_FLOAT  Command_Param_ParaType = 16
	Command_Param_DOUBLE Command_Param_ParaType = 17
	Command_Param_STRING Command_Param_ParaType = 32
	Command_Param_BYTES  Command_Param_ParaType = 33
)

var Command_Param_ParaType_name = map[int32]string{
	0:  "Null",
	1:  "UINT8",
	2:  "UINT16",
	3:  "UINT32",
	4:  "UINT64",
	16: "FLOAT",
	17: "DOUBLE",
	32: "STRING",
	33: "BYTES",
}
var Command_Param_ParaType_value = map[string]int32{
	"Null":   0,
	"UINT8":  1,
	"UINT16": 2,
	"UINT32": 3,
	"UINT64": 4,
	"FLOAT":  16,
	"DOUBLE": 17,
	"STRING": 32,
	"BYTES":  33,
}

func (x Command_Param_ParaType) String() string {
	return proto.EnumName(Command_Param_ParaType_name, int32(x))
}

type Command struct {
	Type  Command_CommandType `protobuf:"varint,1,opt,name=type,enum=Report.Command_CommandType" json:"type,omitempty"`
	Paras []*Command_Param    `protobuf:"bytes,2,rep,name=paras" json:"paras,omitempty"`
}

func (m *Command) Reset()         { *m = Command{} }
func (m *Command) String() string { return proto.CompactTextString(m) }
func (*Command) ProtoMessage()    {}

func (m *Command) GetParas() []*Command_Param {
	if m != nil {
		return m.Paras
	}
	return nil
}

type Command_Param struct {
	Type    Command_Param_ParaType `protobuf:"varint,1,opt,name=type,enum=Report.Command_Param_ParaType" json:"type,omitempty"`
	Npara   uint64                 `protobuf:"varint,2,opt,name=npara" json:"npara,omitempty"`
	Dpara   float64                `protobuf:"fixed64,3,opt,name=dpara" json:"dpara,omitempty"`
	Strpara string                 `protobuf:"bytes,4,opt,name=strpara" json:"strpara,omitempty"`
	Bpara   []byte                 `protobuf:"bytes,5,opt,name=bpara,proto3" json:"bpara,omitempty"`
}

func (m *Command_Param) Reset()         { *m = Command_Param{} }
func (m *Command_Param) String() string { return proto.CompactTextString(m) }
func (*Command_Param) ProtoMessage()    {}

type ControlReport struct {
	Tid          uint64   `protobuf:"varint,1,opt,name=tid" json:"tid,omitempty"`
	SerialNumber uint32   `protobuf:"varint,2,opt,name=serial_number" json:"serial_number,omitempty"`
	Command      *Command `protobuf:"bytes,3,opt,name=command" json:"command,omitempty"`
}

func (m *ControlReport) Reset()         { *m = ControlReport{} }
func (m *ControlReport) String() string { return proto.CompactTextString(m) }
func (*ControlReport) ProtoMessage()    {}

func (m *ControlReport) GetCommand() *Command {
	if m != nil {
		return m.Command
	}
	return nil
}

func init() {
	proto.RegisterEnum("Report.Command_CommandType", Command_CommandType_name, Command_CommandType_value)
	proto.RegisterEnum("Report.Command_Param_ParaType", Command_Param_ParaType_name, Command_Param_ParaType_value)
}
