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
		Addr: "127.0.0.1:4150",
	}
	nsqpserver := sha.NewNsqProducer(nsqpconfig)

	// creates a nsqconsumer server
	nsqcconfig := &sha.NsqConsumerConfig{
		Addr:    "127.0.0.1:4150",
		Topic:   "gateway",
		Channel: "ch",
	}
	nsqcserver := sha.NewNsqConsumer(nsqcconfig, nsqpserver)

	// create sha server
	shaserverconfig := &sha.ServerConfig{
		Listener:      listener,
		AcceptTimeout: time.Second,
		Uptopic:       "gateway",
	}
	shaserver := sha.NewServer(srv, nsqpserver, nsqcserver, shaserverconfig)
	sha.SetServer(shaserver)
	shaserver.Start()

	// database
	dbconfig := &sha.DBConfig{
		Host:   "192.168.8.90",
		Port:   "5432",
		User:   "postgres",
		Passwd: "cetc",
		Dbname: "gateway",
	}

	gatewayhub, err := sha.NewGatewayHub(dbconfig)
	sha.SetGatewayHub(gatewayhub)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = gatewayhub.LoadAll()
	if err != nil {
		fmt.Println("connect to db error " + err.Error())
		return
	}
	fmt.Println("gateways has been loaded")
	err = gatewayhub.Listen("gateway")
	if err != nil {
		panic(err)
	}

	go gatewayhub.WaitForNotification()

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
