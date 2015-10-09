package main

import (
	"github.com/giskook/smarthome-access/client"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	box := shb.NewSmarthomebox(189115999977674, "张凯家")
	box.Add(1, 1, 1, "厨房的灯", 1)
	go box.Do("192.168.8.90:8989")

	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Signal: ", <-chSig)
}
