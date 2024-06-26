package main

import (
	"context"
	"fmt"
	"multi-site-dashboard-go/internal"
	"multi-site-dashboard-go/internal/config"
	"multi-site-dashboard-go/internal/database"
	"multi-site-dashboard-go/internal/delivery/kafka"
	kh "multi-site-dashboard-go/internal/delivery/kafka/handler"
	"multi-site-dashboard-go/internal/delivery/rest"
	"multi-site-dashboard-go/internal/delivery/sse"
	"multi-site-dashboard-go/internal/repository"
	uc "multi-site-dashboard-go/internal/usecase"
	"os"
	"os/signal"
	"time"

	_ "multi-site-dashboard-go/docs"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var logger *zap.Logger

func main() {
	ctx := context.Background()

	// Init config.
	cfg, err := config.ProvideConfig()
	if err != nil {
		panic(fmt.Sprintf("error reading config file: %v", err))
    }

	// Init logger.
	logger, err = internal.WireLogger()
	if err != nil {
		panic(fmt.Sprintf("error setting up logger: %v", err))
    }

	// Setup db.
	if err := database.SetupDatabase(cfg); err != nil {
		logger.Fatal("error setting up db", zap.String("trace", err.Error()))
	}

	// Migrate db.
	m, err := internal.WirePgMigrateInstance(true)
	if err != nil {
		logger.Fatal("error creating db migration instance", zap.String("trace", err.Error()))
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Fatal("error migrating db", zap.String("trace", err.Error()))
	}

	// Init UseCase with dependencies.
	db, err := internal.WirePgConnPool(ctx)
	if err != nil {
		logger.Fatal("error creating db pool", zap.String("trace", err.Error()))
	}
	repo := repository.New(db)
	kw, err := kafka.New(cfg)
	if err != nil {
		logger.Fatal("error creating kafka topics", zap.String("trace", err.Error()))
	}
	b := sse.New()
	uc := uc.NewUseCaseService(repo, kw, b)

	// Create Kafka topics.
	if err := kafka.CreateTopics(cfg); err != nil {
		logger.Fatal("error creating Kafka topics", zap.String("trace", err.Error()))
	}

	// Create server.
	s, err := rest.NewServer(ctx, logger, uc)
	if err != nil {
		logger.Fatal("error creating server", zap.String("trace", err.Error()))
	}

	// Server and dependencies have been successfully initialized from here.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Run server.
	go func() {
		fmt.Printf("starting ECHO server in port %v", cfg.Port)
		if err := s.Echo.Start(fmt.Sprintf(":%v", cfg.Port)); err != nil {
			logger.Fatal("server error", zap.String("trace", err.Error()))
		}
	}()

	// Consume from Kafka topics.
	// Number of goroutines should not be greater than number of partitions in topic.
	kh := kh.NewKafkaHandler(uc)
	go kafka.ConsumeMsgFromMachineResourceUsageTopic(ctx, cfg, kh)

	// Wait for interrupt signal to gracefully shutdown the server with a timeout.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	gracefulShutdown(ctx, s, db, kw)
}

func gracefulShutdown(ctx context.Context, s *rest.Server, db *pgxpool.Pool, kw kafka.KafkaWriter) {
	fmt.Printf("performing graceful shutdown with timeout of %v...", 10*time.Second)

	db.Close()

	if err := kw.Writer.Close(); err != nil {
		logger.Error("failed to close Kafka writer", zap.String("trace", err.Error()))
	}

	if err := s.Echo.Shutdown(ctx); err != nil {
		logger.Error(fmt.Sprintf("failed to shutdown echo server: %v", err.Error()))
	}
}