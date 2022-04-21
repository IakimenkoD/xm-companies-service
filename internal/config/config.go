package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"path"
	"strings"
	"time"
)

type Config struct {
	Environment     string        `mapstructure:"environment"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`

	API      api    `mapstructure:"api"`
	DB       DB     `mapstructure:"db"`
	LogLevel string `mapstructure:"log_level"`

	MQ    MessageQueue `mapstructure:"mq"`
	IpApi ipApi        `mapstructure:"ip_api"`
}

type api struct {
	Address      string        `mapstructure:"address"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	JWTKey       []byte        `mapstructure:"jwt_key"`
}

type ipApi struct {
	Address string        `mapstructure:"address"`
	Timeout time.Duration `mapstructure:"timeout"`
}

type MessageQueue struct {
	Address string `mapstructure:"address"`
	Queue   string `mapstructure:"queue"`
}

type DB struct {
	URL          string `mapstructure:"url"`
	SchemaName   string `mapstructure:"schema_name"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

var defaults = map[string]interface{}{
	"environment":      "development",
	"shutdown_timeout": time.Second * 5,
	"version":          "dev",

	"db.url":            "postgres://root@localhost:5432/root?sslmode=disable",
	"db.schema_name":    "xm",
	"db.max_open_conns": 2,
	"db.max_idle_conns": 2,

	"api.address":       ":4000",
	"api.read_timeout":  time.Second * 5,
	"api.write_timeout": time.Second * 5,
	"api.jwt_key":       []byte("IGdvdCBhIHNlY3JldCBjYW4geW91IGtlZXAgaXQ="),

	"ip_api.address": "https://ipapi.co/",
	"ip_api.timeout": time.Second * 5,

	"mq.address": "amqp://guest:guest@localhost:5672/",
	"mq.queue":   "company_updated",

	"log_level": "debug",
}

func New(destination string, log *zap.Logger) (*Config, error) {
	directory, basename := path.Split(destination)
	name := strings.TrimSuffix(basename, path.Ext(basename))

	viper.AddConfigPath(directory)
	viper.SetConfigName(name)

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	for key, value := range defaults {
		viper.SetDefault(key, value)
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Info("can't read config using defaults")
	}

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}

	return &c, nil
}
