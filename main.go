package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
)

const (
	envPort      = "PORT"
	envRedisAddr = "REDIS_ADDR"
	envRedisPass = "REDIS_PASSWORD"
	envRedisDB   = "REDIS_DB"
)

func main() {
	logger := logrus.New()

	logger.SetLevel(logrus.TraceLevel)

	redisDB, err := strconv.Atoi(os.Getenv(envRedisDB))
	if err != nil {
		logger.WithError(err).Panic("REDIS_DB must be an integer")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv(envRedisAddr),
		Password: os.Getenv(envRedisPass),
		DB:       redisDB,
	})

	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		logger.WithError(err).Panic("Error connecting to redis")
	}

	logger.Infof("Connected with redis: %s", client)

	server := NewServer(logger, client)

	server.Listen()

	<-server.connClose
	logger.Info("Server shutdown complete")
}
