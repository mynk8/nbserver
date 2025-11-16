package websockets

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

type InMemorySessionManager struct {
	sessions map[string]*PTY
	mu       sync.RWMutex
}

func NewInMemorySessionManager() *InMemorySessionManager {
	return &InMemorySessionManager{
		sessions: make(map[string]*PTY),
	}
}

func (sm *InMemorySessionManager) CreateSession(tty *os.File, cmd *exec.Cmd) string {
	sessionID := sm.generateSessionID()
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.sessions[sessionID] = &PTY{TTY: tty, cmd: cmd}
	return sessionID
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
