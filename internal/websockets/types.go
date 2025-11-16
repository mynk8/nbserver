package websockets

import (
	"os"
	"os/exec"
	"time"

	"github.com/gorilla/websocket"
)

var WebsocketType = map[int]string{
	websocket.BinaryMessage: "binary",
	websocket.TextMessage: "text",
	websocket.CloseMessage: "close",
	websocket.PingMessage: "ping",
	websocket.PongMessage: "pong",
}

type Hub struct {
	clients map[*Client]bool
	broadcast chan []byte
	register chan *Client
	unregister chan *Client
}

type Client struct {
	hub *Hub
	conn *websocket.Conn

	send chan []byte
}

type TTYSize struct {
	Cols uint16 `json:"cols"`
	Rows uint16 `json:"rows"`
	X uint16 `json:"x"`
	Y uint16 `json:"y"`
}

type PTY struct {
	TTY *os.File
	cmd *exec.Cmd
}

const (
	MaxBufferSizeBytes = 256
	KeepAlivePingTimeOut = 20 * time.Second
	DefaultConnectionErrorLimit = 12
)
