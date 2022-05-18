package config

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	PgMigrationsPath string `envconfig:"PG_MIGRATIONS_PATH"`
	PgPort           string `envconfig:"PG_PORT"`
	PgUser           string `envconfig:"PG_USER"`
	PgAddr           string `envconfig:"PG_ADDR"`
	PgPassword       string `envconfig:"PG_Password"`
	PgDB             string `envconfig:"PG_DB"`
}

var (
	config Config
	once   sync.Once
)

func Get() *Config {
	once.Do(func() {
		ReRead()
	})

	return &config
}

func ReRead() *Config {
	log.Println("reading app config")
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err)
	}
	_, err = json.MarshalIndent(config, "", "")
	configBytes, err := json.MarshalIndent(config, "", "")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Configuration", string(configBytes))
	return &config
}

func WriteConfig(stringWithConfig func() string) error {
	envMap, err := godotenv.Unmarshal(stringWithConfig())
	if err != nil {
		return err
	}
	return godotenv.Write(envMap, "./config/config.env")

}
