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

type DatabaseConfiguration struct {
	Host         string
	Port         string
	User         string
	Passwd       string
	Dbname       string
	Monitortable string
}

type Configuration struct {
	NsqConfig *NsqConfiguration
	DbConfig  *DatabaseConfiguration
}

func ReadConfig(confpath string) (*Configuration, error) {
	file, _ := os.Open(confpath)
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)

	return &configuration, err
}
