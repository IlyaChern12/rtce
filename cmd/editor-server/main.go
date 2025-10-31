package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/IlyaChern12/rtce/internal/config"
	"github.com/IlyaChern12/rtce/internal/db"
	"github.com/IlyaChern12/rtce/internal/redisdb"
	"github.com/redis/go-redis/v9"
)

func main() {
	// подгружаем конфиги
	cfg := config.LoadConfig()
	log.Println("Loaded config:", cfg)

	// коннект к бд
	dbConn, err := db.PostgresConnect(cfg.DB_DSN)
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}
	defer func() {
    if err := dbConn.Close(); err != nil {
        log.Printf("Failed to close Postgres connection: %v", err)
    }
	}()

	// коннект к редису
	rdb, err := redisdb.RedisConnect(cfg.RedisAddr)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer func() {
    if err := rdb.Close(); err != nil {
        log.Printf("Failed to close Redis connection: %v", err)
    }
	}()

	app := &App{
		DB:    dbConn,
		Redis: rdb,
	}

	// основные хэндлеры

	// хэндлер на проверку на то жив ли сервис
	http.HandleFunc("/health", app.health)
	// хэндлер на готовность
	http.HandleFunc("/ready", app.ready)

	log.Println("Server started on port", cfg.Port)
	log.Fatal(http.ListenAndServe(":" + cfg.Port, nil))
}

type App struct {
	DB    *sql.DB
	Redis *redis.Client
}

func (a *App) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Service is ready"))
	if err != nil {
		log.Printf("Failed to write service status: %v", err)
	}
}

func (a *App) ready(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3 * time.Second)
	defer cancel()

	// подтянут ли postgres
	if err := a.DB.PingContext(ctx); err != nil {
		log.Printf("Postgres not ready: %v", err)
		http.Error(w, "Postgres is not ready", http.StatusServiceUnavailable)
		return
	}
	// подтянут ли redis
	if err := a.Redis.Ping(ctx).Err(); err != nil {
		log.Printf("Redis not ready: %v", err)
		http.Error(w, "Redis is not ready", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Service is ready"))
	if err != nil {
		log.Printf("Failed to write service status: %v", err)
	}
}