package config

import (
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
    IsDebug    bool    `env:"CLIENT1_IS_DEBUG" env-default:"false"`
    BindIp     string  `env:"CLIENT1_LISTENER_HOST" env-default:"0.0.0.0"`
    Port       int     `env:"CLIENT1_LISTENER_PORT" env-default:"8090"`
    ServerIp   string  `env-default:"subpub_server"`
    ServerPort int     `env:"SERVER_LISTENER_PORT" env-default:"8089"`
}

var instance *Config
var once     sync.Once

func GetConfig() *Config{
    once.Do(func() {
        instance = &Config{}
        cleanenv.ReadEnv(instance)
    })

    return instance
}

