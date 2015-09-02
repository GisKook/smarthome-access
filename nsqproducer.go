package sha

import (
	"log"
	"sync"

	"github.com/bitly/go-nsq"
)

type NsqProducerConfig struct {
	Addr    string
	Topic   string
	Channel string
}

type NsqProducer struct {
	config    *NsqProducerConfig
	waitGroup *sync.WaitGroup

	producer *nsq.Producer
}

func NewNsqProducer(config *NsqProducerConfig) *NsqProducer {
	return &NsqProducer{
		config:    config,
		waitGroup: &sync.WaitGroup{},
	}
}

func (s *NsqProducer) Start() {
	q.waitGroup.Add(1)
	config := nsq.NewConfig()

	var errmsg error
	s.producer, errmsg = nsq.NewProducer(s.config.Addr, config)

	if errmsg != nil {
		log.Printf("create producer error" + errmsg.Error())
	}

}

func (s *NsqProducer) Stop() {
	s.waitGroup.Done()
	s.waitGroup.Wait()

	s.producer.Stop()
}

func (s *NsqProducer) Send(topic string, value []byte) error {
	err := s.producer.Publish(topic, value)

	return err
}
