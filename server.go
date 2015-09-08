package sha

import (
	"github.com/giskook/gotcp"
	"net"
	"time"
)

type ServerConfig struct {
	Listener      *net.TCPListener
	AcceptTimeout time.Duration
}

type Server struct {
	srv         *gotcp.Server
	nsqproducer *NsqProducer
	nsqconsumer *NsqConsumer
}

func NewServer(srv *gotcp.Server, nsqproducer *NsqProducer, nsqconsumer *NsqConsumer) *Server {
	return &Server{
		srv:         srv,
		nsqproducer: nsqproducer,
		nsqconsumer: nsqconsumer,
	}
}

func (s *Server) Start(config *ServerConfig) {
	go s.nsqproducer.Start()
	go s.nsqconsumer.Start()

	go s.srv.Start(config.Listener, config.AcceptTimeout)
}

func (s *Server) Stop() {
	s.srv.Stop()
	s.nsqproducer.Stop()
	s.nsqconsumer.Stop()
}
