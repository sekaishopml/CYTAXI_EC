package entity

import "time"

type SessionID string

type Session struct {
	ID             SessionID
	ConversationID ConversationID
	Phone          string
	Status         SessionStatus
	CurrentState   ConversationState
	CreatedAt      time.Time
	UpdatedAt      time.Time
	ExpiresAt      time.Time
	Metadata       map[string]string
}

type SessionStatus string

const (
	SessionStatusActive  SessionStatus = "active"
	SessionStatusIdle    SessionStatus = "idle"
	SessionStatusExpired SessionStatus = "expired"
	SessionStatusClosed  SessionStatus = "closed"
)

func NewSession(phone string) *Session {
	now := time.Now()
	return &Session{
		ID:             SessionID(NewID()),
		Phone:          phone,
		Status:         SessionStatusActive,
		CurrentState:   StateNew,
		CreatedAt:      now,
		UpdatedAt:      now,
		ExpiresAt:      now.Add(30 * time.Minute),
		Metadata:       make(map[string]string),
	}
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

func (s *Session) Touch() {
	s.UpdatedAt = time.Now()
	s.ExpiresAt = time.Now().Add(30 * time.Minute)
}

func (s *Session) MarkIdle() {
	s.Status = SessionStatusIdle
	s.UpdatedAt = time.Now()
}

func (s *Session) MarkExpired() {
	s.Status = SessionStatusExpired
	s.UpdatedAt = time.Now()
}

func (s *Session) MarkClosed() {
	s.Status = SessionStatusClosed
	s.CurrentState = StateClosed
	s.UpdatedAt = time.Now()
}

func (s *Session) Reactivate() {
	s.Status = SessionStatusActive
	s.UpdatedAt = time.Now()
	s.ExpiresAt = time.Now().Add(30 * time.Minute)
}
