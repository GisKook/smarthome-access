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
		Topic:   "write_test",
		Channel: "ch",
	}
	nsqcserver := sha.NewNsqConsumer(nsqcconfig)

	shaserver := sha.NewServer(srv, nsqpserver, nsqcserver)
	shaserverconfig := &sha.ServerConfig{
		Listener:      listener,
		AcceptTimeout: time.Second,
	}
	shaserver.Start(shaserverconfig)

	// database

	dbconfig := &sha.DBConfig{
		Host:   "192.168.1.155",
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
	err = gatewayhub.Listen("gateway")
	if err != nil {
		panic(err)
	}

	gatewayhub.WaitForNotification()

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
