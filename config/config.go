package config

import (
	"os"
	"path"
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
	BrokerURL string `yaml:"BrokerURL"`
	BrokerPort int  `yaml:"BrokerPort"`
	Topic string `yaml:"Topic"`
}

type TimeScaleConfig struct {
	Server string `yaml:"server"`
	Port int `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Config struct {
	Environment string `yaml:"environment"`
	Port int `yaml:"port"`
	LogDir string `yaml:"logDir"`
	TimeScale TimeScaleConfig `yaml:"timeScale"`
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
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(cwd)
	viper.AddConfigPath(path.Join(cwd, "config")) 
	err = viper.ReadInConfig()
	if err != nil {
		return err
	}

	viper.SetDefault("environment", Development.String())
	viper.SetDefault("port", 8000)
	viper.SetDefault("logDir", path.Join(cwd, "log"))

	err = viper.Unmarshal(&config)
	if err != nil {
		return err
	}
	return nil
}