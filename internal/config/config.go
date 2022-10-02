package config

import (
	"authstore/pkg/logging"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IsDebug  *bool    `yaml:"is_debug" env-default:"true"`
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database" env-required:"true"`
}
type Server struct {
	WriteTimeout uint8 `yaml:"write_timeout" env-default:"15"`
	ReedTimeout  uint8 `yaml:"reed_timeout" env-default:"15"`
	Listen       struct {
		Type   string `yaml:"type" env-default:"port"`
		BindIP string `yaml:"bind_ip" env-default:"127.0.0.1"`
		Port   string `yaml:"port" env-default:"8080"`
	} `yaml:"listen"`
}
type Database struct {
	Mysql struct {
		Host            string `yaml:"host" env-required:"true"`
		Port            string `yaml:"port" env-required:"true"`
		DBname          string `yaml:"dbname" env-required:"true"`
		Username        string `yaml:"username"`
		Password        string `yaml:"password"`
		ConnMaxLifetime int    `yaml:"conn_max_life_time"`
		MaxOpenConns    int    `yaml:"max_open_conns"`
		MaxIdleConns    int    `yaml:"max_idle_conns"`
	} `yaml:"mysql" env-required:"true"`
}

var instance *Config

var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("Read applicateion config")
		instance = &Config{}
		if err := cleanenv.ReadConfig("configs/config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})

	return instance
}
