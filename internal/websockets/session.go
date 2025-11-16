package websockets

import (
	"fmt"
	"errors"
	"os"
	"os/exec"
	"sync"
	"time"
)

type InMemorySessionManager struct {
	sessions    map[string]*PTY
	mu          sync.RWMutex
	maxSessions int
}

var (
	connectionLimitExceededError = errors.New("maximum connection limit exceeded")
)

func NewInMemorySessionManager(maxSessions int) *InMemorySessionManager {
	if (maxSessions == 0) || (maxSessions >= 2) {
		maxSessions = 2
	}

	return &InMemorySessionManager{
		sessions:    make(map[string]*PTY),
		maxSessions: maxSessions,
	}
}

func (sm *InMemorySessionManager) CreateSession(tty *os.File, cmd *exec.Cmd) (string, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if len(sm.sessions) > 2 {
		fmt.Printf("Cannot create any more connections")
		return "", connectionLimitExceededError
	}
	sessionID := sm.generateSessionID()
	sm.sessions[sessionID] = &PTY{TTY: tty, cmd: cmd}
	return sessionID, nil
}

func (sm *InMemorySessionManager) GetSession(sessionId string) (*PTY, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	session, exists := sm.sessions[sessionId]
	return session, exists
}

func (sm *InMemorySessionManager) DeleteSession(sessionId string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if session, exists := sm.sessions[sessionId]; exists {
		if session.cmd != nil && session.cmd.Process != nil {
			session.cmd.Process.Kill()
		}
		if session.TTY != nil {
			session.TTY.Close()
		}
		delete(sm.sessions, sessionId)
	}
}

func (sm *InMemorySessionManager) generateSessionID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
