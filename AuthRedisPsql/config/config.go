package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBPort        string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	JWTSecret     string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	// dbPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))   // 将 DB_PORT 转换为 int，假设它是一个端口号
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB")) // 将 REDIS_DB 转换为

	return Config{
		DBHost:        os.Getenv("DB_HOST"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		DBPort:        os.Getenv("DB_PORT"),
		RedisAddr:     os.Getenv("REDIS_ADDR"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       redisDB,
		JWTSecret:     os.Getenv("JWT_SECRET"),
	}
}
