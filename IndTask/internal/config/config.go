package config

import (
	"github.com/Baraulia/goLab/IndTask.git/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

var logger = logging.GetLogger()

type Config struct {
	IsDebug bool `yaml:"is_debug" env:"IS_DEBUG"  env-default:"true"`
	Listen  struct {
		Type   string `yaml:"type" env:"LISTEN_TYPE" env-default:"port"`
		Port   string `yaml:"port"  env:"LISTEN_PORT" env-default:"8080"`
		BindIp string `yaml:"bind_ip"  env:"LISTEN_HOST" env-default:"127.0.0.1"`
	} `yaml:"listen"`
	DB struct {
		Host     string `yaml:"host" env:"DB_HOST" env-default:"127.0.0.1"`
		Port     string `yaml:"port" env:"DB_PORT" env-default:"5432"`
		Username string `yaml:"username" env:"DB_USERNAME" env-default:"postgres"`
		Password string `yaml:"password" env:"DB_PASSWORD" env-default:"postgres"`
		DBName   string `yaml:"dbname" env:"DB_NAME" env-default:"postgres"`
		SSLMode  string `yaml:"sslmode" env:"DB_SSLMODE" env-default:"disable"`
	} `yaml:"db"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger.Info("read application configuration")
		instance = &Config{}

		if err := cleanenv.ReadConfig("configs/config.yaml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
		}
		if err := cleanenv.ReadConfig(".env", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}
