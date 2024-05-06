package ws

import (
	"context"
	"fmt"
	"io"
	"multi-site-dashboard-go/internal"
	"multi-site-dashboard-go/internal/config"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	writeWait = 10 * time.Second
)

var (
	wg sync.WaitGroup
	lock sync.Mutex
	clients = make(map[*Client]struct{})
	logger, _ = internal.WireLogger()
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

type Client struct {
	conn *websocket.Conn
	eventChan chan []byte
}

type Websocket struct {
	clients map[*Client]struct{}
}

func New() Websocket {
	return Websocket{clients: clients}
}


func (ws Websocket) Broadcast(ctx context.Context, data []byte) error {
	// For sender, closing of channel is not needed once it is done.
	// When websocket connection is closed, the client will be removed.
	if ctx.Err() != nil {
		return ctx.Err()
	}

	for client := range ws.clients {
		client.eventChan <- data
	}
	return nil
}

func NewHTTPServer(ctx context.Context, cfg *config.Config) *http.Server {
	s := &http.Server{Addr: fmt.Sprintf(":%v", cfg.WebsocketPort)}
	http.HandleFunc("/", heartbeatHandler)
	http.HandleFunc("/ws", WebsocketHandler)
	return s
}

func heartbeatHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Websocket server is alive")
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        logger.Error("unable to establish websocket connection", zap.String("trace", err.Error()))
        return
    }

	client := addClient(conn)
	defer removeClient(client, r)

	wg.Add(2)
	go reader(client)
	go writer(client)
	wg.Wait()
}

func addClient(conn *websocket.Conn) *Client {
	lock.Lock()
	defer lock.Unlock()
	client := &Client{conn: conn, eventChan: make(chan []byte)}
	clients[client] = struct{}{}
	return client
}

func removeClient(client *Client, r *http.Request) {
	logger.Info(fmt.Sprintf("client %v disconnected from websocket and removing from list", r.RemoteAddr))
	lock.Lock()
	defer lock.Unlock()
	delete(clients, client) 
}

func reader(client *Client) {
	defer wg.Done()
	for {
		_, p, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure){
				logger.Error("unexpected websocket close error", zap.String("trace", err.Error()))
			}
			break
		} 
		fmt.Println(string(p))
	}
}

func writer(client *Client) {
	defer func() {
		wg.Done()
		if err := client.conn.Close(); err != nil {
			logger.Error("unable to close websocket connection", zap.String("trace", err.Error()))
		}
	}()
	for {
		select {
		case data, ok := <- client.eventChan:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Channel is closed.
				return 
			}
			_ = client.conn.WriteMessage(websocket.TextMessage, data)
		default:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				// Connection is closed.
				return
			}
		}
	}
}