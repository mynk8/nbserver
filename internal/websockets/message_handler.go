package websockets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

type TerminalMessageHandler struct {
	terminal Terminal
}

func NewTerminalMessageHandler(terminal Terminal) *TerminalMessageHandler {
	return &TerminalMessageHandler{
		terminal: terminal,
	}
}

func (h *TerminalMessageHandler) HandleResize(data []byte, tty *os.File) error {
	ttySize := &TTYSize{}
	resizeMessage := bytes.Trim(data[1:], " \n\r\t\x00\x01")
	if err := json.Unmarshal(resizeMessage, ttySize); err != nil {
		return fmt.Errorf("failed to unmarshal resize message: %w", err)
	}

	if err := h.terminal.Resize(tty, ttySize.Cols, ttySize.Rows); err != nil {
		return fmt.Errorf("failed to handle resize: %w", err)
	}
	return nil
}

func (h *TerminalMessageHandler) HandleInput(data []byte, tty *os.File) error {
	dataBuffer := bytes.Trim(data, "\x00")
	if err := h.terminal.Write(tty, dataBuffer); err != nil {
		return fmt.Errorf("failed to handle input: %w", err)
	}
	return nil
}
