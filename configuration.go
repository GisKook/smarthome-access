package sha

import (
	"encoding/json"
	"os"
)

type NsqConfiguration struct {
	Addr        string
	UpTopic     string
	DownTopic   string
	Downchannel string
}

type ServerConfiguration struct {
	BindPort          string
	ReadLimit         uint16
	WriteLimit        uint16
	ConnTimeout       uint16
	ConnCheckInterval uint16
	ServerStatistics  uint16
}

type Configuration struct {
	NsqConfig    *NsqConfiguration
	ServerConfig *ServerConfiguration
}

func ReadConfig(confpath string) (*Configuration, error) {
	file, _ := os.Open(confpath)
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)

	return &configuration, err
}

func (conf *Configuration) GetServerReadLimit() uint16 {
	return conf.ServerConfig.ReadLimit
}

func (conf *Configuration) GetServerWriteLimit() uint16 {
	return conf.ServerConfig.WriteLimit
}

func (conf *Configuration) GetServerConnCheckInterval() uint16 {
	return conf.ServerConfig.ConnCheckInterval
}

func (conf *Configuration) GetServerStatistics() uint16 {
	return conf.ServerConfig.ServerStatistics
}

var Config *Configuration

func SetConfiguration(config *Configuration) {
	Config = config
}

func GetConfiguration() *Configuration {
	return Config
}
