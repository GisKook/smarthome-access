package main

import (
	"github.com/bitly/go-nsq"
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
	"log"
)

func main() {
	config := nsq.NewConfig()
	var w *nsq.Producer
	w, _ = nsq.NewProducer("127.0.0.1:4150", config)

	replogin := &Report.ControlReport{
		Tid:          1000,
		SerialNumber: 1,
		Command: &Report.Command{
			Type:  Report.Command_CMT_REQLOGIN,
			Paras: nil,
		},
	}
	reqdata, err := proto.Marshal(replogin)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	log.Println("send topic")

	err = w.Publish("write_test", reqdata)
	if err != nil {
		log.Panic("Could not connect")
	}

	w.Stop()
}
