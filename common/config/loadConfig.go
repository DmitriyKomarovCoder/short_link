package config

import (
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func LoadConfig() (*Config, error) {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current working directory: ", err)
	}

	configFilePath := filepath.Join(currentDir, "common", "config", "config.yaml")

	configFile, err := os.Open(configFilePath)
	if err != nil {
		log.Fatal("Error opening config.yaml file: ", err)
	}
	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			log.Fatal("Error close config file: ", err)
		}
	}(configFile)

	var config Config
	decoder := yaml.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal("Error decoding config.yaml file: ", err)
	}

	postgresPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatal("Error converting postgres port to int: ", err)
	}

	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatal("Error converting redis port to int")
	}

	config.Postgres.Name = os.Getenv("DB_NAME")
	config.Postgres.User = os.Getenv("DB_USER")
	config.Postgres.Password = os.Getenv("DB_PASSWORD")
	config.Postgres.Host = os.Getenv("DB_HOST")
	config.Postgres.Port = postgresPort

	config.Redis.Address = os.Getenv("REDIS_ADDRESS")
	config.Redis.DB = redisDB

	return &config, nil
}
