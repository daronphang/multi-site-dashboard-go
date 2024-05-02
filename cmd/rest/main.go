package main

import (
	"context"
	"fmt"
	"multi-site-dashboard-go/internal/delivery/rest"
	"os"
	"os/signal"
	"time"
)

func main() {
	s, err := rest.NewServer()
	if err != nil {
		panic(err.Error())
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		fmt.Printf("starting ECHO server in port %v", s.Cfg.Port)
		if err := s.Echo.Start(fmt.Sprintf(":%v", s.Cfg.Port)); err != nil {
			s.Logger.Fatal(fmt.Sprintf("failed to start server: %v", err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s.GracefulShutdown(ctx)
}