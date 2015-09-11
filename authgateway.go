package sha

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

type Device struct {
	Oid     []byte
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
	Uid         []byte
	Devicecount uint16
	Devicelist  []Device
}

type GatewayHub struct {
	Db *sql.DB
	//	Gateway map[string]*GatewayProperty
}

func (g *GatewayHub) LoadAll() bool {
	st, err := g.Db.Prepare("select * from gateway")
	if err != nil {
		fmt.Println("can not select ")

		return false
	}
	st.Query()
	//	r, err2 := st.Query()
	//	if err2 != nil {
	//		fmt.Println("can not select")
	//
	//		return false
	//	}

	return true

}

func NewGatewayHub(conn *DBConfig) (*GatewayHub, error) {
	connstring := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", conn.User, conn.Passwd, conn.Host, conn.Port, conn.Dbname)
	db, err := sql.Open("postgres", connstring) //	if err != nil {
		fmt.Println("conn to database fail " + err.Error())

		return nil, err
	}

	return nil, err
}

func main() {
	config := &DBConfig{
		Host:   "192.168.1.155",
		Port:   "5432",
		User:   "postgres",
		Passwd: "cetca",
		Dbname: "gateway",
	}

	NewGatewayHub(config)
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
