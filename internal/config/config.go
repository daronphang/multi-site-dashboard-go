package config

import (
	"os"
	"path"
	"runtime"
	"sync"

	"github.com/spf13/viper"
)

type Environment int32

const (
	Development Environment = iota
	Production
	Testing
)

func (e Environment) String() string {
	switch e {
	case Development:
		return "DEVELOPMENT"
	case Production:
		return "PRODUCTION"
	case Testing:
		return "TESTING"
	}
	return "unknown"
}

type KafkaConfig struct {
	BrokerAddresses string `yaml:"BrokerAddresses"` // localhost:9092,localhost:9093
	MachineResourceUsageGroupID string `yaml:"machineResourceUsageGroupId"`
}

type PostgresConfig struct {
	Server string `yaml:"server"`
	Port int `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName string `yaml:"dbName"`
}

type Config struct {
	Environment string `yaml:"environment"`
	Port int `yaml:"port"`
	LogDir string `yaml:"logDir"`
	Postgres PostgresConfig `yaml:"postgres"`
	Kafka KafkaConfig `yaml:"kafka"`
	WebsocketPort int `yaml:"websocketPort"`
	SSEPort int `yaml:"ssePort"`
}

var syncOnceConfig sync.Once
var config *Config

func ProvideConfig() (*Config, error) {
	var err error
	syncOnceConfig.Do(func() {
		err = readConfigFromFile()
	})
	if err != nil {
		return config, err
	}
	return config, nil
}

func readConfigFromFile() error {
	env := os.Getenv("GO_ENV")
	_, filename, _, _ := runtime.Caller(0)
	if env == "TESTING" {
		viper.SetConfigName("config.test")
	} else {
		viper.SetConfigName("config")
	}
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path.Dir(filename))
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	viper.SetDefault("environment", Development.String())
	viper.SetDefault("port", 8000)
	viper.SetDefault("logDir", path.Join(path.Dir(filename), "../../log"))

	err = viper.Unmarshal(&config)
	if err != nil {
		return err
	}
	return nil
}