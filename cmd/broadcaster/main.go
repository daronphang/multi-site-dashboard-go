package main

import (
	"context"
	"fmt"
	"multi-site-dashboard-go/internal"
	"multi-site-dashboard-go/internal/config"
	"multi-site-dashboard-go/internal/delivery/sse"
	"os"
	"os/signal"
	"time"
)

var (
	logger, _ = internal.WireLogger()
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	cfg, err := config.ProvideConfig()
	if err != nil {
		panic(fmt.Sprintf("error reading config file: %v", err))
    }

	s := sse.NewHTTPServer(ctx, cfg)
	fmt.Printf("starting SSE server on port %v...", cfg.SSEPort)
	go func() {
		if err := s.ListenAndServe(); err != nil {
			logger.Fatal(err.Error())
		}
	}()
	// Wait for interrupt signal to gracefully shutdown the server with a timeout.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		logger.Fatal(err.Error())
	}
}