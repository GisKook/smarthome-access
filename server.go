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
	config           *ServerConfig
	srv              *gotcp.Server
	nsqproducer      *NsqProducer
	nsqconsumer      *NsqConsumer
	database         *ExecDatabase
	checkconnsticker *time.Ticker
}

var Gserver *Server

func SetServer(server *Server) {
	Gserver = server
}

func GetServer() *Server {
	return Gserver
}

func NewServer(srv *gotcp.Server, nsqproducer *NsqProducer, nsqconsumer *NsqConsumer, config *ServerConfig, db *ExecDatabase) *Server {
	serverstatistics := GetConfiguration().GetServerStatistics()
	return &Server{
		config:           config,
		srv:              srv,
		nsqproducer:      nsqproducer,
		nsqconsumer:      nsqconsumer,
		database:         db,
		checkconnsticker: time.NewTicker(time.Duration(serverstatistics) * time.Second),
	}
}

func (s *Server) GetProducer() *NsqProducer {
	return s.nsqproducer
}

func (s *Server) GetConsumer() *NsqConsumer {
	return s.nsqconsumer
}

func (s *Server) GetDatabase() *ExecDatabase {
	return s.database
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

func (s *Server) checkStatistics() {
	for {
		<-c.checkconnsticker.C

	}
}
