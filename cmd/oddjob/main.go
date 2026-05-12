package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/its-the-vibe/OddJob/internal/config"
	"github.com/its-the-vibe/OddJob/internal/dispatcher"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	configPath := flag.String("config", "config.json", "Path to JSON config file")
	envFile := flag.String("env-file", ".env", "Path to .env file")
	flag.Parse()

	if err := godotenv.Load(*envFile); err != nil && !os.IsNotExist(err) {
		log.Fatalf("load env file: %v", err)
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Username: cfg.Redis.Username,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	defer client.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("redis ping failed: %v", err)
	}

	registry := dispatcher.NewRegistry(dispatcher.NewHelloWorldTransformer())
	service := dispatcher.NewService(cfg, client, registry, log.Default())

	if err := service.Run(ctx); err != nil {
		log.Fatalf("service failed: %v", err)
	}
}
