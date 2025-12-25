package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Address       string        `yaml:"address" env:"SEARCH_ADDRESS" env-default:"0.0.0.0:8082"`
	DBAddress     string        `yaml:"db_address" env:"DB_ADDRESS" env-required:"true"`
	WordsAddress  string        `yaml:"words_address" env:"WORDS_ADDRESS" env-required:"true"`
	LogLevel      string        `yaml:"log_level" env:"LOG_LEVEL" env-default:"INFO"`
	IndexTTL      time.Duration `yaml:"index_ttl" env:"INDEX_TTL" env-default:"20s"`
	BrokerAddress string        `yaml:"broker_address" env:"BROKER_ADDRESS" env-default:"nats://nats:4222"`
}

func MustLoad(configPath string) Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Printf("config file %s does not exist, reading from env", configPath)
		var cfg Config
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			log.Fatalf("cannot read config from env: %s", err)
		}
		return cfg
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	return cfg
}
