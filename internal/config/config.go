package config

import "os"

type Config struct {
	// Connection strings
	PostgresURL string
	RabbitURL   string
	RedisURL    string
	// Rabbit exchange name
	Exchange string
	// Queue Names
	QueueDB    string
	QueueBack  string
	QueueCache string
	// Routing key names
	KeyDB    string
	KeyBack  string
	KeyCache string
}

func New() *Config {
	return &Config{
		PostgresURL: getEnv("POSTGRES_URL", "postgres://postgres:demopsw@localhost:5432/messenger"),
		RabbitURL:   getEnv("RABBIT_URL", "amqp://guest:guest@localhost:5672"),
		RedisURL:    getEnv("REDIS_URL", "localhost:6379"),
		Exchange:    getEnv("EXCHANGE", "main_exchange"),
		QueueDB:     getEnv("QUEUE_DB", "db_queue"),
		QueueBack:   getEnv("QUEUE_BACK", "backend_queue"),
		QueueCache:  getEnv("QUEUE_CACHE", "cache_queue"),
		KeyDB:       getEnv("KEY_DB", "db_key"),
		KeyBack:     getEnv("KEY_BACK", "backend_key"),
		KeyCache:    getEnv("KEY_CACHE", "cache_key"),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
