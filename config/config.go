package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Server   ServerConfig
	Database DatabaseConfig
	AWS      AWSConfig
	RabbitMQ RabbitMQConfig
	JWT      JWTConfig
}

type AppConfig struct {
	Name    string
	Version string
	Env     string
}

type ServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	Uri          string
	DatabaseName string
}

type AWSConfig struct {
	S3Bucket        string
	Region          string
	AccessKeyId     string
	SecretAccessKey string
}

type RabbitMQConfig struct {
	Uri      string
	Services []string
	Exchange string
}

type JWTConfig struct {
	Secret string
}

func LoadConfig(dir string) {
	err := godotenv.Load(dir)
	if err != nil {
		log.Printf("No se pudo cargar el archivo .env, usando variables de entorno del sistema: %v", err)
	}
}

func NewConfig() *Config {
	LoadConfig(".env")
	config := &Config{
		App: AppConfig{
			Name:    os.Getenv("APP_NAME"),
			Version: os.Getenv("APP_VERSION"),
			Env:     os.Getenv("APP_ENV"),
		},
		Server: ServerConfig{
			Port:         getEnvAsInt("SERVER_PORT", 8080),
			ReadTimeout:  getEnvAsDuration("SERVER_READ_TIMEOUT", 10) * time.Second,
			WriteTimeout: getEnvAsDuration("SERVER_WRITE_TIMEOUT", 10) * time.Second,
		},
		JWT: JWTConfig{
			Secret: os.Getenv("JWT_SECRET"),
		},
		Database: DatabaseConfig{
			Uri:          os.Getenv("DATABASE_URI"),
			DatabaseName: os.Getenv("DATABASE_NAME"),
		},
		AWS: AWSConfig{
			S3Bucket:        os.Getenv("AWS_S3_BUCKET"),
			Region:          os.Getenv("AWS_REGION"),
			AccessKeyId:     os.Getenv("AWS_ACCESS_KEY_ID"),
			SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		},
		RabbitMQ: RabbitMQConfig{
			Uri:      os.Getenv("RABBITMQ_URI"),
			Exchange: os.Getenv("RABBITMQ_EXCHANGE"),
			Services: []string{
				os.Getenv("RABBITMQ_USER_SERVICE"),
				os.Getenv("RABBITMQ_AUTH_SERVICE"),
			},
		},
	}

	return config
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := os.Getenv(name)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

func getEnvAsDuration(name string, defaultVal int) time.Duration {
	valueStr := os.Getenv(name)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return time.Duration(value)
	}
	return time.Duration(defaultVal)
}
