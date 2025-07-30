package factory

import (
	"github.com/networkgcorefullcode/scp/util"
)

const (
	SCP_EXPECTED_CONFIG_VERSION = "1.0.0"
	SCP_HTTP_PORT               = 8080
)

type Config struct {
	Info          *Info          `yaml:"info"`
	Configuration *Configuration `yaml:"configuration"`
	Logger        *util.Logger   `yaml:"logger"`
	CfgLocation   string
	Rcvd          bool
}

type Info struct {
	Version     string `yaml:"version,omitempty"`
	Description string `yaml:"description,omitempty"`
}

type Configuration struct {
	ScpName   string     `yaml:"scpName,omitempty"`
	ScpDBName string     `yaml:"ScpDBName,omitempty"`
	PortHttp  int        `yaml:"portHttp,omitempty"`
	Mongodb   *Mongodb   `yaml:"mongodb,omitempty"`
	KafkaInfo *KafkaInfo `yaml:"kafkaInfo,omitempty"`
}

type Mongodb struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
}

type KafkaInfo struct {
	EnableKafka *bool  `yaml:"enableKafka,omitempty"`
	BrokerUri   string `yaml:"brokerUri,omitempty"`
	BrokerPort  int    `yaml:"brokerPort,omitempty"`
	Topic       string `yaml:"topicName,omitempty"`
}

func (c *Config) GetVersion() string {
	if c.Info != nil && c.Info.Version != "" {
		return c.Info.Version
	}
	return ""
}
