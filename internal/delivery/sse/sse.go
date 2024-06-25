package sse

import (
	"context"
	"fmt"
	"io"
	"multi-site-dashboard-go/internal"
	"multi-site-dashboard-go/internal/config"
	"net/http"
	"sync"
	"time"
)

var (
	lock sync.Mutex
	clients = make(map[*Client]struct{})
	logger, _ = internal.WireLogger()
)

type SSE struct {
	clients map[*Client]struct{}
}

type Client struct {
	eventChan chan string
}

func (sse SSE) Broadcast(ctx context.Context, data []byte) error {
	// For sender, closing of channel is not needed once it is done.
	// When SSE connection is closed, the client will be removed.
	if ctx.Err() != nil {
		return ctx.Err()
	}

	for client := range sse.clients {
		client.eventChan <- string(data)
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

	client := addClient()
	defer removeClient(client, r)
	
	// Establish connection to client.
	w.(http.Flusher).Flush()

	for {
		select {
		case <-r.Context().Done():
			return
		case data := <- client.eventChan:
			fmt.Fprintf(w, "data: %s\n\n", data)
			w.(http.Flusher).Flush()
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func addClient() *Client {
	lock.Lock()
	defer lock.Unlock()
	client := &Client{eventChan: make(chan string)}
	clients[client] = struct{}{}
	return client
}

func removeClient(client *Client, r *http.Request) {
	logger.Info(fmt.Sprintf("client %v disconnected from SSE and removing from list", r.RemoteAddr))
	lock.Lock()
	defer lock.Unlock()
	delete(clients, client) 
}