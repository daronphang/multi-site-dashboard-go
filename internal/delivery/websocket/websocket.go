package websocket

import (
	"errors"
	"fmt"
	"io"
	"multi-site-dashboard-go/internal"
	"multi-site-dashboard-go/internal/config"
	"net/http"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var (
	conn *websocket.Conn
	logger, _ = internal.WireLogger()
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func ProvideWebsocketConn() (*websocket.Conn, error) {
	if conn == nil {
		return nil, errors.New("websocket connection has not been established")
	}
	return conn, nil
}

func StartHTTPServer() {
	cfg, err := config.ProvideConfig()
	if err != nil {
		logger.Fatal(err.Error())
	}
	http.HandleFunc("/", heartbeatHandler)
	http.HandleFunc("/ws", WebsocketHandler)

	if err := http.ListenAndServe(fmt.Sprintf(":%v", cfg.WebsocketPort), nil); err != nil {
		logger.Fatal(err.Error())
	}
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
	reader(conn)
}

func reader(conn *websocket.Conn) {
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			logger.Error("unable to establish websocket connection", zap.String("trace", err.Error()))
			return
		}
		fmt.Println(string(p))
	}
}

func WriteMsgToWebsocket(data []byte) error {
	conn, err := ProvideWebsocketConn()
	if err != nil {
		return err
	}
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return err
	}
	return nil
}