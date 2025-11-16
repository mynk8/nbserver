package websockets

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WSConnection struct {
	conn *websocket.Conn
}

func NewWSConnection(conn *websocket.Conn) *WSConnection {
	return &WSConnection{conn: conn}
}

func (w *WSConnection) ReadMessage() (messageType int, data []byte, err error) {
	return w.conn.ReadMessage()
}

func (w *WSConnection) WriteMessage(messageType int, data []byte) error {
	return w.conn.WriteMessage(messageType, data)
}

func (w *WSConnection) SetPongHandler(handler func(string) error) {
	w.conn.SetPongHandler(handler)
}

func (w *WSConnection) Close() error {
	return w.conn.Close()
}

// KeepAliveManager manages keep-alive functionality for websocket connections
type KeepAliveManager struct {
	pingInterval time.Duration
	timeout      time.Duration
}

func NewKeepAliveManager(pingInterval, timeout time.Duration) *KeepAliveManager {
	return &KeepAliveManager{
		pingInterval: pingInterval,
		timeout:      timeout,
	}
}

func (k *KeepAliveManager) Start(conn ConnectionManager, waiter *sync.WaitGroup) {
	defer waiter.Done()

	lastPongTime := time.Now()
	conn.SetPongHandler(func(msg string) error {
		lastPongTime = time.Now()
		return nil
	})

	ticker := time.NewTicker(k.pingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := conn.WriteMessage(websocket.PingMessage, []byte("keepalive")); err != nil {
				return
			}

			if time.Since(lastPongTime) > k.timeout {
				return
			}
		}
	}
}
