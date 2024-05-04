package sse

import (
	"context"
	"fmt"
	"io"
	"multi-site-dashboard-go/internal"
	"multi-site-dashboard-go/internal/config"
	"net/http"
	"time"
)

var (
	clients = make(map[chan string]struct{})
	logger, _ = internal.WireLogger()
)

type SSE struct {
	clients map[chan string]struct{}
}

func (sse SSE) Broadcast(ctx context.Context, data []byte) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	for client := range sse.clients {
		client <- string(data)
	}
	return nil
}

func New() SSE {
	return SSE{clients: clients}
}

func NewHTTPServer(ctx context.Context, cfg *config.Config) *http.Server {
	s := &http.Server{Addr: fmt.Sprintf(":%v", cfg.SSEPort)}
	http.HandleFunc("/", heartbeatHandler)
	http.HandleFunc("/sse", SSEHandler)
	return s
}

func heartbeatHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "SSE server is alive")
}

func SSEHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/event-stream")

	eventChan := make(chan string)
	clients[eventChan] = struct{}{}

	// Remove the client when they disconnect.
	defer func() {
		logger.Info(fmt.Sprintf("client %v disconnected and removing from list", r.RemoteAddr))
		delete(clients, eventChan) 
		close(eventChan)
	   }()
	
	// Establish connection to client.
	w.(http.Flusher).Flush()

	for {
		select {
		case <-r.Context().Done():
			return
		case data := <- eventChan:
			fmt.Fprintf(w, "data: %s\n\n", data)
			w.(http.Flusher).Flush()
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}