package websockets

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

// TerminalHandler handles terminal websocket connections
type TerminalHandler struct {
	terminal       Terminal
	sessionManager SessionManager
	messageHandler *TerminalMessageHandler
	keepAliveMgr   *KeepAliveManager
	command        string
	errorLimit     int
	bufferSize     int
}

func NewTerminalHandler(
	terminal Terminal,
	sessionManager SessionManager,
	command string,
) *TerminalHandler {
	keepAliveMgr := NewKeepAliveManager(
		KeepAlivePingTimeOut/2,
		KeepAlivePingTimeOut,
	)

	return &TerminalHandler{
		terminal:       terminal,
		sessionManager: sessionManager,
		messageHandler: NewTerminalMessageHandler(terminal),
		keepAliveMgr:   keepAliveMgr,
		command:        command,
		errorLimit:     DefaultConnectionErrorLimit,
		bufferSize:     MaxBufferSizeBytes,
	}
}

// Connect handles a new websocket connection for terminal access
func (th *TerminalHandler) Connect(w http.ResponseWriter, r *http.Request, h http.Header) error {
	connection, err := upgrader.Upgrade(w, r, h)
	if err != nil {
		return fmt.Errorf("failed to upgrade connection: %w", err)
	}
	defer connection.Close()

	conn := NewWSConnection(connection)

	tty, cmd, err := th.terminal.Start(th.command)
	if err != nil {
		return fmt.Errorf("failed to start terminal: %w", err)
	}

	sessionId, err := th.sessionManager.CreateSession(tty, cmd)
	if err != nil {
		if tty != nil {
			tty.Close()
		}
		if cmd != nil && cmd.Process != nil {
			cmd.Process.Kill()
		}
		return fmt.Errorf("failed to create connection")
	}
	defer th.sessionManager.DeleteSession(sessionId)

	var waiter sync.WaitGroup
	waiter.Add(2)

	go th.keepAliveMgr.Start(conn, &waiter)

	go th.readFromTerminal(tty, conn, &waiter)

	th.writeToTerminal(tty, conn)

	waiter.Wait()
	return nil
}

func (th *TerminalHandler) writeToTerminal(tty *os.File, conn ConnectionManager) {
	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			return
		}

		if messageType == websocket.BinaryMessage && len(data) > 0 && data[0] == 1 {
			if err := th.messageHandler.HandleResize(data, tty); err != nil {
				fmt.Printf("failed to handle resize: %v\n", err)
			}
			continue
		}

		if err := th.messageHandler.HandleInput(data, tty); err != nil {
			fmt.Printf("failed to write data to terminal: %v\n", err)
		}
	}
}

func (th *TerminalHandler) readFromTerminal(tty *os.File, conn ConnectionManager, waiter *sync.WaitGroup) {
	defer waiter.Done()
	errorCounter := 0

	for {
		if errorCounter > th.errorLimit {
			break
		}

		buffer := make([]byte, th.bufferSize)
		readLength, err := th.terminal.Read(tty, buffer)
		if err != nil {
			break
		}

		if err := conn.WriteMessage(websocket.BinaryMessage, buffer[:readLength]); err != nil {
			errorCounter++
			continue
		}

		errorCounter = 0
	}
}

func getMessageType(messageType int) string {
	if dtype, ok := WebsocketType[messageType]; ok {
		return dtype
	}
	return "unknown"
}
