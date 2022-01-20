package config

import (
	"github.com/Baraulia/goLab/IndTask.git/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

var logger = logging.GetLogger()

type Config struct {
	IsDebug  bool   `yaml:"is_debug" env:"IS_DEBUG"  env-default:"true"`
	FilePath string `yaml:"file_path" env:"FILE_PATH" env-default:""`
	Listen   struct {
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
		SSLMode  string `yaml:"sslmode" env:"DB_SSL_MODE" env-default:"disable"`
	} `yaml:"db"`
	Mail struct {
		From     string `yaml:"from" env:"MAIL_FROM"`
		Password string `yaml:"password" env:"MAIL_PASSWORD"`
		SmtpHost string `yaml:"smtpHost" env:"MAIL_SMTP_HOST"`
		SmtpPort string `yaml:"smtpPort" env:"MAIL_SMTP_PORT"`
		Subject  string `yaml:"subject" env:"MAIL_SUBJECT"`
	} `yaml:"mail"`
	ProfitBook struct {
		Profitability   float32 `yaml:"profitability" env:"PROFIT_BOOK_PROFITABILITY"`
		MaxRentalNumber float32 `yaml:"max_rental_number" env:"PROFIT_BOOK_MAX_RENTAL_NUMBER"`
	} `yaml:"profit_book"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger.Info("read application configuration")
		instance = &Config{}

		err := cleanenv.ReadConfig("configs/config.yaml", instance)
		if err != nil {
			logger.Info(err)
			if err := cleanenv.ReadConfig(".env", instance); err != nil {
				help, _ := cleanenv.GetDescription(instance, nil)
				logger.Info(help)
				logger.Fatal(err)
			}
		}

	})
	return instance
}
