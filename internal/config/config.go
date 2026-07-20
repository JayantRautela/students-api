package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTpServer struct {
	Addr string
}

type Config struct {
	Env         string `yml:"env" env:"ENV" env-required:"true"` //env-default:"production"
	StoragePath string `yml:"storage_path" env-required:"true"`
	HTTpServer `yml:"http_server"`
}
// these are struct tags

func MustLoad() *Config {
	var configPath string 

	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "", "path to the configuration file")
		flag.Parse()

		configPath = *flags

		if configPath == "" {
			log.Fatal("Config Path is not set")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exists: %s", configPath)
	} 

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)

	if err != nil {
		log.Fatalf("Cannot read config file: %s", err.Error())
	}

	return &cfg
}