package config

import (
	"sync"
	"vk-test/pkg/logger"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
    IsDebug  bool    `env:"IS_DEBUG" env-default:"false"`
    BindIp   string  `env:"LISTENER_HOST" env-default:"0.0.0.0"`
    Port     int     `env:"LISTENER_PORT" env-default:"8089"`
}

var instance *Config
var once     sync.Once

func GetConfig(l logger.Logger) *Config{
    once.Do(func() {
        
        l.Info("Read application configuration")

        err := godotenv.Load(".env")

        if err != nil {
            l.Fatal(err)
        }

        instance = &Config{}

        if err := cleanenv.ReadEnv(instance); err != nil {
            help, _ := cleanenv.GetDescription(instance, nil)
            l.Info(help)
            l.Fatal(err)
        }
    })

    return instance
}

