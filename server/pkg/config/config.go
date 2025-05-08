package config

import (
	"sync"
	"vk-test/pkg/logger"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
    IsDebug  bool    `env:"SERVER_IS_DEBUG" env-default:"false"`
    BindIp   string  `env:"SERVER_LISTENER_HOST" env-default:"0.0.0.0"`
    Port     int     `env:"SERVER_LISTENER_PORT" env-default:"8089"`
}

var instance *Config
var once     sync.Once

func GetConfig(l logger.Logger) *Config{
    once.Do(func() {
        
        l.Info("Read application configuration")

        instance = &Config{}

        if err := cleanenv.ReadEnv(instance); err != nil {
            help, _ := cleanenv.GetDescription(instance, nil)
            l.Info(help)
            l.Fatal(err)
        }
    })

    return instance
}

