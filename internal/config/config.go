package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env    string `yaml:"env" env:"ENV" env-default:"local" env-required:"true"`
	Server `yaml:"server"`
	Client `yaml:"client"`
}

type Server struct {
	Host string `yaml:"host" env-default:"localhost" env-required:"true"`
	Port string `yaml:"port" env-default:"8080" env-required:"true"`
}

type Client struct {
	Address string `yaml:"address" env-default:"server:8080" env-required:"true"`
}

func MustLoad() *Config {
	cp := os.Getenv("CONFIG_PATH")
	if cp == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(cp); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s, err: %s", cp, err.Error())
	}

	var cfg Config

	if err := cleanenv.ReadConfig(cp, &cfg); err != nil {
		log.Fatalf("cannot read config: %s, err: %s", cp, err.Error())
	}

	return &cfg
}
