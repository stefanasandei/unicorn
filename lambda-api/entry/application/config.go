package application

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RedisAddress    string
	RabbitMQAddress string
	ServerAddr      string
}

func LoadConfig() Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(".env file not found")
	}

	cfg := Config{
		RedisAddress:    "localhost:6379",
		RabbitMQAddress: "amqp://guest:guest@localhost:5672/",
		ServerAddr:      "127.0.0.1:3000",
	}

	if redisAddr, exists := os.LookupEnv("REDIS_ADDR"); exists {
		cfg.RedisAddress = redisAddr
	}

	if rabbitMQAddress, exists := os.LookupEnv("RABBITMQ_ADDR"); exists {
		cfg.RabbitMQAddress = rabbitMQAddress
	}

	if serverAddress, exists := os.LookupEnv("SERVER_ADDR"); exists {
		cfg.ServerAddr = serverAddress
	}

	return cfg
}
