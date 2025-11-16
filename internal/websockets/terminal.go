package websockets

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/creack/pty"
)

type PTYTerminal struct{}

func NewPTYTerminal() *PTYTerminal {
	return &PTYTerminal{}
}

func (t *PTYTerminal) Start(command string) (*os.File, *exec.Cmd, error) {
	cmd := exec.Command(command)
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start pty: %w", err)
	}
	return ptmx, cmd, nil
}

func (t *PTYTerminal) Resize(tty *os.File, cols, rows uint16) error {
	if err := pty.Setsize(tty, &pty.Winsize{
		Rows: rows,
		Cols: cols,
	}); err != nil {
		return fmt.Errorf("failed to resize terminal: %w", err)
	}
	return nil
}

func (t *PTYTerminal) Write(tty *os.File, data []byte) error {
	_, err := tty.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write %d bytes to tty: %w", len(data), err)
	}
	return nil
}

func (t *PTYTerminal) Read(tty *os.File, buffer []byte) (int, error) {
	return tty.Read(buffer)
}
