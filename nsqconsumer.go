package sha

import (
	"log"
	"sync"

	"github.com/bitly/go-nsq"
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

type NsqConsumerConfig struct {
	Addr    string
	Topic   string
	Channel string
}

type NsqConsumer struct {
	config    *NsqConsumerConfig
	waitGroup *sync.WaitGroup

	consumer *nsq.Consumer
	producer *NsqProducer
}

func NewNsqConsumer(config *NsqConsumerConfig, producer *NsqProducer) *NsqConsumer {
	return &NsqConsumer{
		config:    config,
		waitGroup: &sync.WaitGroup{},
		producer:  producer,
	}
}

func (s *NsqConsumer) recvNsq() {
	s.consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		data := message.Body
		command := &Report.ControlReport{}
		err := proto.Unmarshal(data, command)
		if err != nil {
			log.Println("unmarshal error")
		}
		gatewayid := command.Tid
		serialnum := command.SerialNumber
		switch command.GetCommand().Type {
		case Report.Command_CMT_REQLOGIN:
			log.Println("login")
			log.Println(gatewayid)
			log.Println(serialnum)
			replogin := &Report.ControlReport{
				Tid:          gatewayid,
				SerialNumber: serialnum,
				Command: &Report.Command{
					Type: Report.Command_CMT_REPLOGIN,
					Paras: []*Report.Command_Param{
						&Report.Command_Param{
							Type:  Report.Command_Param_UINT8,
							Npara: 1,
						},
					},
				},
			}
			repdata, err := proto.Marshal(replogin)
			if err != nil {
				log.Fatal("marshaling error: ", err)
			}
			log.Println("send topic")
			s.producer.Send("topic", repdata)

		}
		//		controlpacket := NewControlPacket(10000, 20000, 1)
		//		NewConns().GetConn(10000).SendToGateway(controlpacket)

		return nil
	}))
}

func (s *NsqConsumer) Start() {
	s.waitGroup.Add(1)
	defer func() {
		s.waitGroup.Done()
		err := recover()
		if err != nil {
			log.Println("err found")
			s.Stop()
		}

	}()

	config := nsq.NewConfig()

	var errmsg error
	s.consumer, errmsg = nsq.NewConsumer(s.config.Topic, s.config.Channel, config)

	if errmsg != nil {
		//	panic("create consumer error -> " + errmsg.Error())
		log.Println("create consumer error -> " + errmsg.Error())
	}
	s.recvNsq()

	err := s.consumer.ConnectToNSQD(s.config.Addr)
	if err != nil {
		panic("Counld not connect to nsq -> " + err.Error())
	}
}

func (s *NsqConsumer) Stop() {
	s.waitGroup.Wait()

	errmsg := s.consumer.DisconnectFromNSQD(s.config.Addr)

	if errmsg != nil {
		log.Printf("stop consumer error ", errmsg.Error())
	}

	s.consumer.Stop()
}
