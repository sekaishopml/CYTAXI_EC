package whatsapp

import (
	"context"
	"sync"
	"time"
)

type SessionManager struct {
	mu       sync.RWMutex
	sessions map[string]*Session
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*Session),
	}
}

func (sm *SessionManager) CreateSession(ctx context.Context, id string) (*Session, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session := &Session{
		ID:        id,
		Status:    StatusQRReady,
		CreatedAt: time.Now(),
	}
	sm.sessions[id] = session
	return session, nil
}

func (sm *SessionManager) GetSession(ctx context.Context, id string) (*Session, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	session, ok := sm.sessions[id]
	if !ok {
		return nil, nil
	}
	return session, nil
}

func (sm *SessionManager) SetQRCode(ctx context.Context, id string, qr *QRCode) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session, ok := sm.sessions[id]
	if !ok {
		return ErrNotConnected
	}

	session.QRCode = qr
	session.Status = StatusQRReady
	return nil
}

func (sm *SessionManager) UpdateStatus(ctx context.Context, id string, status ConnectionStatus) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session, ok := sm.sessions[id]
	if !ok {
		return ErrNotConnected
	}

	session.Status = status
	if status == StatusConnected {
		session.QRCode = nil
	}
	return nil
}

func (sm *SessionManager) DeleteSession(ctx context.Context, id string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.sessions, id)
	return nil
}
