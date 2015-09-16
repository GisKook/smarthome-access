package main

import (
	"database/sql"
	"encoding/binary"
	"fmt"
	"github.com/lib/pq"
	"strconv"
	"sync"
	"time"
)

type Device struct {
	Oid     uint64
	Type    uint8
	Company uint16
}

type DBConfig struct {
	Host   string
	Port   string
	User   string
	Passwd string
	Dbname string
}

type GatewayProperty struct {
	Uid         uint64
	Devicecount uint16
	Devicelist  []Device
}

type GatewayHub struct {
	Db      *sql.DB
	Gateway map[uint64]*GatewayProperty

	Listener  *pq.Listener
	waitGroup *sync.WaitGroup
}

func char2byte(c string) byte {
	switch c {
	case "0":
		return 0
	case "1":
		return 1
	case "2":
		return 2
	case "3":
		return 3
	case "4":
		return 4
	case "5":
		return 5
	case "6":
		return 6
	case "7":
		return 7
	case "8":
		return 8
	case "9":
		return 9
	case "a":
		return 10
	case "b":
		return 11
	case "c":
		return 12
	case "d":
		return 13
	case "e":
		return 14
	case "f":
		return 15
	}
	return 0
}

func macaddr2uint64(mac []uint8) uint64 {
	var buffer []byte
	buffer = append(buffer, 0)
	buffer = append(buffer, 0)
	value := char2byte(string(mac[0]))*16 + char2byte(string(mac[1]))
	buffer = append(buffer, value)
	value = char2byte(string(mac[3]))*16 + char2byte(string(mac[4]))
	buffer = append(buffer, value)
	value = char2byte(string(mac[6]))*16 + char2byte(string(mac[7]))
	buffer = append(buffer, value)
	value = char2byte(string(mac[9]))*16 + char2byte(string(mac[10]))
	buffer = append(buffer, value)
	value = char2byte(string(mac[12]))*16 + char2byte(string(mac[13]))
	buffer = append(buffer, value)
	value = char2byte(string(mac[15]))*16 + char2byte(string(mac[16]))
	buffer = append(buffer, value)

	return binary.BigEndian.Uint64(buffer)
}

func (g *GatewayHub) add(gatewayid uint64, deviceid uint64, devicetype uint8, company uint16) {
	_, ok := g.Gateway[gatewayid]
	if !ok {
		g.Gateway[gatewayid] = &GatewayProperty{
			Uid:         gatewayid,
			Devicecount: 1,
			Devicelist: []Device{
				{
					Oid:     deviceid,
					Type:    devicetype,
					Company: company,
				},
			},
		}
	} else {
		g.Gateway[gatewayid].Devicelist = append(g.Gateway[gatewayid].Devicelist, Device{
			Oid:     deviceid,
			Type:    devicetype,
			Company: company,
		})
		g.Gateway[gatewayid].Devicecount++
	}
}

func (g *GatewayHub) LoadAll() error {
	st, err := g.Db.Prepare("select deviceid, gatewayid, devicetype, company from gateway")
	if err != nil {
		return err
	}

	r, er := st.Query()
	if er != nil {
		return er
	}
	defer st.Close()

	var dmac []uint8
	var gmac []uint8
	var devicetype uint8
	var company uint16
	for r.Next() {
		err = r.Scan(&dmac, &gmac, &devicetype, &company)
		if err != nil {
			return err
		}
		gatewayid := macaddr2uint64(gmac)
		deviceid := macaddr2uint64(dmac)
		add(gatewayid, deviceid, devicetype, company)
	}
	defer r.Close()

	return nil
}

func (g *GatewayHub) Listen(table string) error {
	return g.Listener.Listen(table)
}

func NewGatewayHub(conn *DBConfig) (*GatewayHub, error) {
	connstring := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", conn.User, conn.Passwd, conn.Host, conn.Port, conn.Dbname)
	db, err := sql.Open("postgres", connstring)
	if err != nil {
		return nil, err
	}

	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	return &GatewayHub{
		Db:        db,
		Gateway:   make(map[uint64]*GatewayProperty),
		Listener:  pq.NewListener(connstring, 10*time.Second, time.Minute, reportProblem),
		waitGroup: &sync.WaitGroup{},
	}, nil
}

func (g *GatewayHub) parsepayload(payload string) (uint64, uint64, uint8, uint16) {
	values := strings.Split(payload, '^')
	deviceid := macaddr2uint64(values[1])
	gatewayid := macaddr2uint64(values[2])
	devicetype := uint8(strconv.Atoi(values[3]))
	company := uint16(strconv.Atoi(values[4]))

	return deviceid, gatewayid, devicetype, company
}

func (g *GatewayHub) insert(payload string) {
	deviceid, gatewayid, devicetype, company := parsepayload(payload)
	add(gatewayid, deviceid, devicetype, company)
}

func (g *GatewayHub) waitForNotification() {
	for {
		select {
		case notify := <-g.Listener.Notify:
			fmt.Println(notify.Extra)
			switch notify.Extra[0] {
			case 'U':
			case 'I':
			case 'D':
			}
			break
		case <-time.After(90 * time.Second):
			go func() {
				g.Listener.Ping()
			}()
			// Check if there's more work available, just in case it takes
			// a while for the Listener to notice connection loss and
			// reconnect.
			fmt.Println("received no work for 90 seconds, checking for new work")
			break
		}
	}
}

func main() {
	config := &DBConfig{
		Host:   "192.168.1.155",
		Port:   "5432",
		User:   "postgres",
		Passwd: "cetc",
		Dbname: "gateway",
	}

	gatewayhub, err := NewGatewayHub(config)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = gatewayhub.LoadAll()
	err = gatewayhub.Listen("gateway")
	if err != nil {
		panic(err)
	}

	gatewayhub.waitForNotification()

}

//func main() {
//	db, err := sql.Open("postgres", "user=postgres password=cetc dbname=gateway host=192.168.1.155 port=5432 sslmode=disable")
//	if err != nil {
//		log.Println("dddddddddddddddd")
//		log.Fatal(err)
//	}
//	rows, err2 := db.Query("select * from gateway")
//	if err2 != nil {
//		log.Println("aaaaaaaaaadddddd")
//		log.Fatal(err2)
//	}
//	log.Println(rows)
//
//}
