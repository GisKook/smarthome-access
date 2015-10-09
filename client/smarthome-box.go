package shb

import (
	"encoding/binary"
	"log"
	"net"
	"sync"
	"time"
)

func CheckSum(cmd []byte, cmdlen uint16) byte {
	temp := cmd[0]
	for i := uint16(1); i < cmdlen; i++ {
		temp ^= cmd[i]
	}

	return temp
}

type Device struct {
	DeviceID   uint64
	DeviceType uint8
	Company    uint16
	Name       string
	Status     uint8
}

type Smarthomebox struct {
	GatewayID   uint64
	Name        string
	DeviceCount uint16
	DeviceList  []*Device
	Wg          *sync.WaitGroup
	ExitChan    chan struct{}
}

func NewSmarthomebox(gatewayid uint64, name string) *Smarthomebox {
	return &Smarthomebox{
		GatewayID:   gatewayid,
		Name:        name,
		DeviceCount: 0,
		DeviceList:  nil,
		Wg:          &sync.WaitGroup{},
		ExitChan:    make(chan struct{}),
	}
}

func (b *Smarthomebox) Close() {
	close(b.ExitChan)
}

func (b *Smarthomebox) Check(deviceid uint64) bool {
	for i := uint16(0); i < b.DeviceCount; i++ {
		if b.DeviceList[i].DeviceID == deviceid {
			return true
		}
	}

	return false
}

func (b *Smarthomebox) Del(deviceid uint64) {
	for i := uint16(0); i < b.DeviceCount; i++ {
		if b.DeviceList[i].DeviceID == deviceid {
			b.DeviceCount--
			b.DeviceList[i] = nil
			b.DeviceList = append(b.DeviceList[:i], b.DeviceList[i+1:]...)
			return
		}
	}
}

func (b *Smarthomebox) Add(deviceid uint64, devicetype uint8, company uint16, name string, status uint8) {
	device := &Device{
		DeviceID:   deviceid,
		DeviceType: devicetype,
		Company:    company,
		Name:       name,
		Status:     status,
	}
	if b.Check(deviceid) {
		b.Del(deviceid)
		b.DeviceCount--
	}

	b.DeviceList = append(b.DeviceList, device)
	b.DeviceCount++
}

func (b *Smarthomebox) login(conn *net.TCPConn) uint8 {
	logincmd := []byte{0xCE, 0x00, 0x00, 0x00, 0x01}
	gatewayid_byte := make([]byte, 8)
	binary.BigEndian.PutUint64(gatewayid_byte, b.GatewayID)
	logincmd = append(logincmd, gatewayid_byte[2:]...)
	logincmd = append(logincmd, byte(len(b.Name)))
	logincmd = append(logincmd, []byte(b.Name)...)
	logincmd = append(logincmd, byte(1))
	logincmd = append(logincmd, byte(1))
	devicecount_byte := make([]byte, 2)
	binary.BigEndian.PutUint16(devicecount_byte, b.DeviceCount)
	logincmd = append(logincmd, devicecount_byte...)
	for i := uint16(0); i < b.DeviceCount; i++ {
		logincmd = append(logincmd, byte(b.DeviceList[i].DeviceType))
		logincmd = append(logincmd, byte(6))
		deviceid_byte := make([]byte, 8)
		binary.BigEndian.PutUint64(deviceid_byte, b.DeviceList[i].DeviceID)
		logincmd = append(logincmd, deviceid_byte[2:]...)
		logincmd = append(logincmd, 0x00)
		logincmd = append(logincmd, 0x01)
		logincmd = append(logincmd, 0x01)
		logincmd = append(logincmd, byte(len(b.DeviceList[i].Name)))
		logincmd = append(logincmd, []byte(b.DeviceList[i].Name)...)
	}
	cmdlen := len(logincmd) + 2 // 2 for checksum and end flag
	binary.BigEndian.PutUint16(logincmd[1:3], uint16(cmdlen))
	logincmd = append(logincmd, CheckSum(logincmd, uint16(cmdlen-2)))
	logincmd = append(logincmd, 0xCE)

	log.Printf("%X\n", logincmd)
	_, err := conn.Write(logincmd)
	if err != nil {
		log.Println(err.Error())

		return 0
	}

	buffer := make([]byte, 1024)
	conn.Read(buffer)
	if buffer[11] == 0x01 {
		log.Println("Login success")
		return 1
	} else {
		log.Println("Login fail")
		return 0
	}
}

func (b *Smarthomebox) heart(conn *net.TCPConn) {
	b.Wg.Add(1)
	defer func() {
		b.Wg.Done()
	}()
	ticker := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-b.ExitChan:
			log.Println("heart exit")
			return
		case <-ticker.C:
			heartcmd := []byte{0xCE, 0x00, 0x0D, 0x00, 0x02}
			heart_byte := make([]byte, 8)
			binary.BigEndian.PutUint64(heart_byte, b.GatewayID)
			heartcmd = append(heartcmd, heart_byte[2:]...)
			heartcmd = append(heartcmd, CheckSum(heartcmd, uint16(len(heartcmd))))
			heartcmd = append(heartcmd, 0xCE)
			log.Printf("%X\n", heartcmd)
			_, err := conn.Write(heartcmd)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}

}

func (b *Smarthomebox) adddeldevice(conn *net.TCPConn) {
	b.Wg.Add(1)
	defer func() {
		b.Wg.Done()
	}()
	ticker := time.NewTicker(3 * time.Second)
	add := true
	devicename := "客厅的沙发"
	for {
		select {
		case <-b.ExitChan:
			log.Println("adddel exit")
			return
		case <-ticker.C:
			adddelcmd := []byte{0xCE, 0x00, 0x00, 0x00, 0x05}
			gatewayid_byte := make([]byte, 8)
			binary.BigEndian.PutUint64(gatewayid_byte, b.GatewayID)
			adddelcmd = append(adddelcmd, gatewayid_byte[2:]...)
			if add {
				adddelcmd = append(adddelcmd, byte(1))
				add = false
			} else {
				adddelcmd = append(adddelcmd, byte(2))
				add = true
			}
			adddelcmd = append(adddelcmd, byte(1))
			adddelcmd = append(adddelcmd, byte(6))
			adddelcmd = append(adddelcmd, []byte{0xFF, 0x00, 0x00, 0x00, 0x00, 0xEE, 0x00, 0x01, 0x01}...)
			adddelcmd = append(adddelcmd, byte(len(devicename)))
			adddelcmd = append(adddelcmd, []byte(devicename)...)
			cmdlen := len(adddelcmd) + 2 // 2 for checksum and end flag
			binary.BigEndian.PutUint16(adddelcmd[1:3], uint16(cmdlen))
			adddelcmd = append(adddelcmd, CheckSum(adddelcmd, uint16(cmdlen-2)))
			adddelcmd = append(adddelcmd, 0xCE)
			log.Printf("add or del %X\n", adddelcmd)
			_, err := conn.Write(adddelcmd)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}
}

func (b *Smarthomebox) recv(conn *net.TCPConn) {
	b.Wg.Add(1)
	defer func() {
		b.Wg.Done()
	}()
	for {
		select {
		case <-b.ExitChan:
			return
		default:
		}

		buffer := make([]byte, 1024)
		length, _ := conn.Read(buffer)
		log.Printf("recv %X\n", buffer[0:length-1])
	}
}

func (b *Smarthomebox) Do(srvaddr string) {
	b.Wg.Add(1)

	tcpaddr, _ := net.ResolveTCPAddr("tcp", srvaddr)

	conn, err := net.DialTCP("tcp", nil, tcpaddr)

	defer func() {
		b.Wg.Done()
		conn.Close()
	}()
	if err != nil {
		log.Println(err.Error())
		return
	}
	if b.login(conn) == 1 {
		go b.heart(conn)
		go b.recv(conn)
		go b.adddeldevice(conn)
	}
	b.Wg.Wait()
}
