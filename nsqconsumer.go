package sha

import (
	"log"
	"sync"

	"github.com/byte/go-nsq"
)

type NsqConsumerConfig struct {
	Addr    string
	Topic   string
	Channel string
}

type NsqConsumer struct {
	config    *NsqConsumerConfig
	waitGroup *sync.WaitGroup

	cosumer *nsq.Consumer
}

func NewNsqConsumer(config *NsqConsumerConfig) *NsqConsumer {
	return &NsqConsumer{
		config:    config,
		waitGroup: &sync.WaitGroup(),
	}
}

func (s *NsqConsumer) recvNsq() {
	s.consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		msg := message.Body
		log.Printf("recv nsq message " + msg)
	}))
}

func (s *NsqConsumer) Start() {
	s.waitGroup.Add(1)
	config := nsq.NewConfig

	var errmsg error
	s.consumer, errmsg = nsq.NewConsumer(s.config.Addr, s.config.Topic, config)

	if errmsg != nil {
		log.Printf("create consumer error " + errmsg.Error())
	}
}

func (s *NsqConsumer) Stop() {
	s.waitGroup.Done()
	s.waitGroup.Wait()

	errmsg := s.consumer.DisconnectFromNSQD(s.config.Addr)

	if errmsg != nil {
		log.Printf("stop consumer error ", errmsg.Error())
	}

	s.consumer.Stop()
}
