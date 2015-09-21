package sha

import (
	"github.com/giskook/gotcp"
	"net"
	"time"
)

type ServerConfig struct {
	Listener      *net.TCPListener
	AcceptTimeout time.Duration
	Uptopic       string
}

type Server struct {
	config      *ServerConfig
	srv         *gotcp.Server
	nsqproducer *NsqProducer
	nsqconsumer *NsqConsumer
}

var server *Server

func SetServer(server *Server) {
	server = server
}

func GetServer() *Server {
	return server
}

func NewServer(srv *gotcp.Server, nsqproducer *NsqProducer, nsqconsumer *NsqConsumer, config *ServerConfig) *Server {
	return &Server{
		config:      config,
		srv:         srv,
		nsqproducer: nsqproducer,
		nsqconsumer: nsqconsumer,
	}
}

func (s *Server) GetProducer() *NsqProducer {
	return s.nsqproducer
}

func (s *Server) GetConsumer() *NsqConsumer {
	return s.nsqconsumer
}

func (s *Server) GetTopic() string {
	return s.config.Uptopic
}

func (s *Server) Start() {
	go s.nsqproducer.Start()
	go s.nsqconsumer.Start()

	go s.srv.Start(s.config.Listener, s.config.AcceptTimeout)
}

func (s *Server) Stop() {
	s.srv.Stop()
	s.nsqproducer.Stop()
	s.nsqconsumer.Stop()
}
