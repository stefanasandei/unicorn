package application

import (
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type Config struct {
	ID              uuid.UUID
	RedisAddress    string
	RabbitMQAddress string
	RuntimesDir     string
}

func LoadConfig() Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(".env file not found")
	}

	cfg := Config{
		ID:              uuid.New(),
		RedisAddress:    "localhost:6379",
		RabbitMQAddress: "amqp://guest:guest@localhost:5672/",
		RuntimesDir:     "./runtimes",
	}

	if rabbitMQAddress, exists := os.LookupEnv("RABBITMQ_ADDR"); exists {
		cfg.RabbitMQAddress = rabbitMQAddress
	}

	if redisAddr, exists := os.LookupEnv("REDIS_ADDR"); exists {
		cfg.RedisAddress = redisAddr
	}

	if runtimesDir, exists := os.LookupEnv("RUNTIMES_DIR"); exists {
		cfg.RuntimesDir = runtimesDir
	}

	if currentEnv, exists := os.LookupEnv("ENV"); exists {
		if currentEnv == "DEBUG" {
			log.Printf("Running in debug mode.")
			cfg.ID = uuid.UUID{}
		}
	}

	return cfg
}
