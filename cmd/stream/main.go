package main

import (
	"context"
	"fmt"
	"multi-site-dashboard-go/internal"
	"multi-site-dashboard-go/internal/config"
	"multi-site-dashboard-go/internal/delivery/stream"
	"multi-site-dashboard-go/internal/delivery/stream/handler"
	"multi-site-dashboard-go/internal/repository"
	"os"
	"os/signal"
	"time"

	uc "multi-site-dashboard-go/internal/usecase"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	cfg, err := config.ProvideConfig()
	if err != nil {
		panic(fmt.Sprintf("error reading config file: %v", err))
    }

	if err := stream.CreateKafkaTopics(cfg); err != nil {
		panic(fmt.Sprintf("error creating Kafka topics: %v", err))
	}

	// Create UseCase.
	db, err := internal.WirePgConnPool(ctx)
	if err != nil {
		panic(err.Error())
	}
	repo := repository.New(db)
	sp := stream.New(cfg)

	uc := uc.NewUseCaseService(repo, sp)

	// Create handler.
	h := handler.NewHandler(uc)

	// Consume from topics.
	fmt.Println("consuming messages from Kafka...")
	go stream.ConsumeMsgsFromMachineResourceUsageTopic(ctx, cfg, h)

	// Wait for interrupt signal to gracefully shutdown the server with a timeout.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream.GracefulShutdown(ctx, sp)
}