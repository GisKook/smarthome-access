package sha

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"strings"
	"sync"
	"time"
)

type UserPasswdHub struct {
	Db   *sql.DB
	User map[uint64]string

	Listener  *pq.Listener
	waitGroup *sync.WaitGroup
}

var UserPasswd *UserPasswdHub

func (u *UserPasswdHub) add(gatewayid uint64, passwd string) {
	u.User[gatewayid] = passwd
}

func (g *UserPasswdHub) LoadAll() error {
	st, err := g.Db.Prepare("select gatewayid, passwd from passwd")
	if err != nil {
		return err
	}

	r, er := st.Query()
	if er != nil {
		return er
	}
	defer st.Close()

	var gmac []uint8
	var passwd string
	for r.Next() {
		err = r.Scan(&gmac, &passwd)
		if err != nil {
			return err
		}
		gatewayid := Macaddr2uint64(gmac)
		g.add(gatewayid, passwd)
	}
	defer r.Close()

	return nil
}

func (g *UserPasswdHub) Listen(table string) error {
	return g.Listener.Listen(table)
}

func NewUserPasswdHub(conn *DBConfig) (*UserPasswdHub, error) {
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

	return &UserPasswdHub{
		Db:        db,
		User:      make(map[uint64]string),
		Listener:  pq.NewListener(connstring, 10*time.Second, time.Minute, reportProblem),
		waitGroup: &sync.WaitGroup{},
	}, nil
}

func (g *UserPasswdHub) parsepayload(payload string) (uint64, string) {
	values := strings.Split(payload, "^")
	gatewayid := Macaddr2uint64([]uint8(values[1]))
	passwd := values[2]

	return gatewayid, passwd
}

func (u *UserPasswdHub) insert(payload string) {
	gatewayid, passwd := u.parsepayload(payload)
	u.add(gatewayid, passwd)
	fmt.Println(u.User)
}

func (g *UserPasswdHub) del(payload string) {
	gatewayid, _ := g.parsepayload(payload)
	delete(g.User, gatewayid)
	fmt.Println(g.User)
}

func (g *UserPasswdHub) update(payload string) {
	gatewayid, passwd := g.parsepayload(payload)
	g.User[gatewayid] = passwd
	fmt.Println(g.User)
}

func (g *UserPasswdHub) WaitForNotification() {
	for {
		select {
		case notify := <-g.Listener.Notify:
			fmt.Println(notify.Extra)
			switch notify.Extra[0] {
			case 'U':
				g.update(notify.Extra)
			case 'I':
				g.insert(notify.Extra)
			case 'D':
				g.del(notify.Extra)
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

func (u *UserPasswdHub) Check(gatewayid uint64, passwd string) bool {
	password, _ := u.User[gatewayid]

	return passwd == password
}

func SetUserPasswdHub(uph *UserPasswdHub) {
	UserPasswd = uph
}

func GetUserPasswdHub() *UserPasswdHub {
	return UserPasswd
}

//func main() {
//	config := &DBConfig{
//		Host:   "192.168.1.155",
//		Port:   "5432",
//		User:   "postgres",
//		Passwd: "cetc",
//		Dbname: "gateway",
//	}
//
//	gatewayhub, err := NewGatewayHub(config)
//	if err != nil {
//		fmt.Println(err.Error())
//		return
//	}
//	err = gatewayhub.LoadAll()
//	err = gatewayhub.Listen("gateway")
//	if err != nil {
//		panic(err)
//	}
//
//	gatewayhub.WaitForNotification()
//}
