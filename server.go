package sha

import (
	"github.com/giskook/gotcp"
	"time"
)

type Config struct {
	HeartBeat  time.Duration
	ReadLimit  int64
	WriteLimit int64
}

type Server struct {
	srv    *gotcp.Server
	nsqsrv *NsqServer
}
