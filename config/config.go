package config

import (
	"bytes"
	"log"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all settings
var defaultConfig = []byte(`
environment: D
http_address: 9000

postgres:
  database: axie
  host: 127.0.0.1
  password: axie
  port: 5402
  sslmode: ~
  timeout: 15
  username: axie
  max_connections: 30
  idle_connections: 10
  max_time_live: 60
`)

type Config struct {
	Base     `mapstructure:",squash"`
	Postgres *Postgres `yaml:"postgres" mapstructure:"postgres"`
}

type Base struct {
	HTTPAddress int `yaml:"http_address" mapstructure:"http_address"`
	GRPCAddress int `mapstructure:"grpc_address"`
	Environment string
}

type Postgres struct {
	Host            string `yaml:"host" mapstructure:"host"`
	Port            int    `yaml:"port" mapstructure:"port"`
	Username        string `yaml:"username" mapstructure:"username"`
	Password        string `yaml:"password" mapstructure:"password"`
	Database        string `yaml:"database" mapstructure:"database"`
	SSLMode         string `yaml:"sslmode" mapstructure:"ssl_mode"`
	Timeout         int    `yaml:"timeout" mapstructure:"timeout"`
	MaxConnections  uint32 `yaml:"max_connections" mapstructure:"max_connections"`
	IdleConnections uint32 `yaml:"idle_connections" mapstructure:"idle_connections"`
	MaxTimeLive     uint32 `yaml:"max_time_live" mapstructure:"max_time_live"`
}

func Load() *Config {
	var cfg = &Config{}

	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer(defaultConfig))
	if err != nil {
		log.Fatal("Failed to read viper config: ", err)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	viper.AutomaticEnv()

	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatal("Failed to unmarshal config: ", err)
	}

	return cfg
}
