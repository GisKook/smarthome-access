package main

import (
	"database/sql"
	//	"encoding/binary"
	"fmt"
	_ "github.com/lib/pq"
	//	"strconv"
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
}

func Char2Byte(c string) byte {
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

func Macaddr2Uint64(mac []uint8) uint64 {
	fmt.Println("%s\n", mac)
	temp := Char2Byte(string(mac[0]))

	fmt.Printf("ddddddddddddddddd%d\n", temp)
	return 0
	//return binary.BigEndian.Uint64(buffer)
}

func (g *GatewayHub) LoadAll() error {
	st, err := g.Db.Prepare("select * from gateway")
	if err != nil {
		return err
	}

	r, er := st.Query()
	if er != nil {
		return er
	}

	var mac []uint8
	for r.Next() {
		err = r.Scan(&mac)
		if err != nil {
			return err
		}
		fmt.Printf("%d\n", Macaddr2Uint64(mac))

	}

	return nil

}

func NewGatewayHub(conn *DBConfig) (*GatewayHub, error) {
	connstring := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", conn.User, conn.Passwd, conn.Host, conn.Port, conn.Dbname)
	db, err := sql.Open("postgres", connstring)
	if err != nil {
		return nil, err
	}

	return &GatewayHub{
		Db:      db,
		Gateway: make(map[uint64]*GatewayProperty),
	}, nil
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
	if err != nil {
		fmt.Println(err.Error())
	}

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
