package websockets

import (
	"os"
	"os/exec"
)

type Terminal interface {
	Start(command string) (*os.File, *exec.Cmd, error)
	Resize(tty *os.File, cols, rows uint16) error
	Write(tty *os.File, data []byte) error
	Read(tty *os.File, buffer []byte) (int, error)
}

type SessionManager interface {
	CreateSession(tty *os.File, cmd *exec.Cmd) (string, error)
	GetSession(sessionId string) (*PTY, bool)
	DeleteSession(sessionId string)
}

type MessageHandler interface {
	HandleResize(data []byte, tty *os.File) error
	HandleInput(data []byte, tty *os.File) error
}

type ConnectionManager interface {
	ReadMessage() (messageType int, data []byte, err error)
	WriteMessage(messageType int, data []byte) error
	SetPongHandler(handler func(string) error)
	Close() error
}
