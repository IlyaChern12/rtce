package config

import (
    "fmt"
    "os"
)

type Config struct {
    Port      string
    DB_DSN    string
    RedisAddr string
    JWTSecret string
}

func LoadConfig() *Config {
    port := os.Getenv("PORT")
    if port == "" { port = "8080" }

    dbDSN := os.Getenv("DB_DSN")
    if dbDSN == "" {
        dbDSN = "postgres://rtce:rtcepass@db:5432/rtce_dev?sslmode=disable"
    }

    redisAddr := os.Getenv("REDIS_ADDR")
    if redisAddr == "" {
        redisAddr = "redis:6379"
    }

     jwtSecret := os.Getenv("JWT_SECRET")
    if jwtSecret == "" {
        jwtSecret = "secret"
    }

    return &Config{
        Port:      port,
        DB_DSN:    dbDSN,
        RedisAddr: redisAddr,
        JWTSecret: jwtSecret,
    }
}

func (c *Config) String() string {
    return fmt.Sprintf("Port=%s, DB_DSN=%s, RedisAddr=%s", c.Port, c.DB_DSN, c.RedisAddr)
}