package main

import (
	"fmt"
	"github.com/giskook/gotcp"
	"github.com/giskook/smarthome-access"
	"os/signal"
	"runtime"
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
	srv := gotcp.NewServer(config, &Callback{}, &echo.EchoProtocol{})

	// creates a nsqproducer server
	nsqpconfig := &sha.NsqProducerConfig{
		Addr: "127.0.0.1:4150",
	}
	nsqpserver := sha.NewNsqProducer(nsqpconfig)

	// creates a nsqconsumer server
	nsqcconfig := &sha.NsqConsumerConfig{
		Addr:    "127.0.0.1:4150",
		Topic:   "topic",
		Channel: "channel",
	}
	nsqcserver := sha.NewNsqConsumer(nsqcconfig)

	shaserver := sha.NewServer(srv, nsqpserver, nsqcserver)
	shaserver.Start()

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
