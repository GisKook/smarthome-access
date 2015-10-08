package main

import (
	"fmt"
	"github.com/giskook/gotcp"
	"github.com/giskook/smarthome-access"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// creates a tcp listener
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":8989")
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	// creates a tcp server
	config := &gotcp.Config{
		PacketSendChanLimit:    20,
		PacketReceiveChanLimit: 20,
	}
	srv := gotcp.NewServer(config, &sha.Callback{}, &sha.ShaProtocol{})

	// creates a nsqproducer server
	nsqpconfig := &sha.NsqProducerConfig{
		Addr:  "127.0.0.1:4150",
		Topic: "sha2app",
	}
	nsqpserver := sha.NewNsqProducer(nsqpconfig)

	// creates a nsqconsumer server
	nsqcconfig := &sha.NsqConsumerConfig{
		Addr:    "127.0.0.1:4150",
		Topic:   "app2sha",
		Channel: "ch",
	}
	nsqcserver := sha.NewNsqConsumer(nsqcconfig, nsqpserver)

	// database passwd_monitor
	dbconfig := &sha.DBConfig{
		Host:   "192.168.8.90",
		Port:   "5432",
		User:   "postgres",
		Passwd: "cetc",
		Dbname: "gateway",
	}

	database, dberr := sha.NewExecDatabase(dbconfig)
	if dberr != nil {
		log.Println(dberr.Error())
	} else {
		log.Println("conn to database success")
	}

	// create sha server
	shaserverconfig := &sha.ServerConfig{
		Listener:      listener,
		AcceptTimeout: time.Second,
		Uptopic:       "sha2app",
	}
	shaserver := sha.NewServer(srv, nsqpserver, nsqcserver, shaserverconfig, database)
	sha.SetServer(shaserver)
	shaserver.Start()

	userhub, err := sha.NewUserPasswdHub(dbconfig)
	sha.SetUserPasswdHub(userhub)
	err = userhub.LoadAll()
	if err != nil {
		fmt.Println("connect to db error " + err.Error())
		return
	}
	fmt.Println("user info has been loaded")
	err = userhub.Listen("passwd")
	if err != nil {
		panic(err)
	}
	go userhub.WaitForNotification()

	// starts service
	fmt.Println("listening:", listener.Addr())

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

	// stops service
	srv.Stop()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
