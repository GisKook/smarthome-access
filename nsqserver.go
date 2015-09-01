package sha

import (
	"log"
	"sync"

	"github.com/bitly/go-nsq"
)

type NsqServerConfig struct {
	Addr    string
	Topic   string
	Channel string
}

type NsqServer struct {
	config    *NsqServerConfig
	waitGroup *sync.WaitGroup

	producer *nsq.Producer
}

func NewNsqServer(config *NsqServerConfig) *NsqServer {
	return &NsqServer{
		config:    config,
		waitGroup: &sync.WaitGroup{},
	}
}

func (s *NsqServer) Start() {
	q.waitGroup.Add(1)
	config := nsq.NewConfig()

	var errmsg error
	s.producer, errmsg = nsq.NewProducer(s.config.Addr, config)

	if errmsg != nil {
		log.Printf("create producer error" + errmsg)
	}
}
