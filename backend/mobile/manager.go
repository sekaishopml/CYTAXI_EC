package mobile

import (
	"fmt"
	"sync"
	"time"
)

type Device struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Type        string    `json:"type"` // ios, android
	Token       string    `json:"token"`
	Active      bool      `json:"active"`
	AppVersion  string    `json:"app_version"`
	OSVersion   string    `json:"os_version"`
	RegisteredAt time.Time `json:"registered_at"`
	LastSeenAt   time.Time `json:"last_seen_at"`
}

type OfflineAction struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Type        string    `json:"type"` // trip_request, accept_trip, reject_trip, complete_trip
	Payload     map[string]any `json:"payload"`
	Status      string    `json:"status"` // queued, syncing, completed, failed
	CreatedAt   time.Time `json:"created_at"`
	SyncedAt    *time.Time `json:"synced_at,omitempty"`
	Retries     int       `json:"retries"`
	MaxRetries  int       `json:"max_retries"`
}

type SyncResult struct {
	ActionID    string    `json:"action_id"`
	Success     bool      `json:"success"`
	ServerRef   string    `json:"server_ref,omitempty"`
	Error       string    `json:"error,omitempty"`
	SyncedAt    time.Time `json:"synced_at"`
	ConflictResolved bool  `json:"conflict_resolved"`
}

type ConflictResolver struct {
	mu sync.RWMutex
}

func NewConflictResolver() *ConflictResolver { return &ConflictResolver{} }

func (cr *ConflictResolver) Resolve(local, remote map[string]any) map[string]any {
	// Server wins by default (backend is source of truth)
	if remote == nil { return local }
	if localTimestamp, lok := local["updated_at"].(string); lok {
		if remoteTimestamp, rok := remote["updated_at"].(string); rok {
			if localTimestamp > remoteTimestamp {
				result := make(map[string]any)
				for k, v := range remote { result[k] = v }
				for k, v := range local { result[k] = v }
				return result
			}
		}
	}
	return remote
}

type PushNotification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Data      map[string]any `json:"data,omitempty"`
	Status    string    `json:"status"` // pending, delivered, failed
	SentAt    time.Time `json:"sent_at"`
	DeliveredAt *time.Time `json:"delivered_at,omitempty"`
}

type Session struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	DeviceID   string    `json:"device_id"`
	Token      string    `json:"token"`
	ExpiresAt  time.Time `json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`
	LastActivity time.Time `json:"last_activity"`
}

type Manager struct {
	mu            sync.RWMutex
	devices       map[string]*Device
	offlineQueue  map[string]chan OfflineAction
	sessions      map[string]*Session
	pushQueue     []PushNotification
	syncResults   []SyncResult
	conflicts     *ConflictResolver
}

func NewManager() *Manager {
	return &Manager{
		devices:      make(map[string]*Device),
		offlineQueue: make(map[string]chan OfflineAction),
		sessions:     make(map[string]*Session),
		conflicts:    NewConflictResolver(),
	}
}

func (m *Manager) RegisterDevice(deviceType, userID, token, appVersion, osVersion string) *Device {
	d := &Device{
		ID: fmt.Sprintf("dev_%d", time.Now().UnixNano()),
		UserID: userID, Type: deviceType, Token: token, Active: true,
		AppVersion: appVersion, OSVersion: osVersion,
		RegisteredAt: time.Now(), LastSeenAt: time.Now(),
	}
	m.mu.Lock()
	m.devices[d.ID] = d
	m.mu.Unlock()
	return d
}

func (m *Manager) GetDevices(userID string) []Device {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []Device
	for _, d := range m.devices {
		if userID == "" || d.UserID == userID {
			result = append(result, *d)
		}
	}
	return result
}

func (m *Manager) RevokeDevice(deviceID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	d, ok := m.devices[deviceID]
	if !ok { return fmt.Errorf("device not found") }
	d.Active = false
	return nil
}

func (m *Manager) EnqueueAction(action OfflineAction) *OfflineAction {
	action.ID = fmt.Sprintf("off_%d", time.Now().UnixNano())
	action.Status = "queued"
	action.CreatedAt = time.Now()
	if action.MaxRetries == 0 { action.MaxRetries = 3 }

	m.mu.Lock()
	if _, ok := m.offlineQueue[action.UserID]; !ok {
		m.offlineQueue[action.UserID] = make(chan OfflineAction, 100)
	}
	m.offlineQueue[action.UserID] <- action
	m.mu.Unlock()

	return &action
}

func (m *Manager) SyncActions(userID string) []SyncResult {
	m.mu.Lock()
	defer m.mu.Unlock()

	ch, ok := m.offlineQueue[userID]
	if !ok { return nil }

	var results []SyncResult
	for len(ch) > 0 {
		action := <-ch
		now := time.Now()
		result := SyncResult{
			ActionID: action.ID, Success: true,
			ServerRef: fmt.Sprintf("srv_%d", now.UnixNano()),
			SyncedAt: now,
		}
		action.Status = "completed"
		action.SyncedAt = &now
		m.syncResults = append(m.syncResults, result)
		results = append(results, result)
	}
	return results
}

func (m *Manager) SendPush(userID, title, body string, data map[string]any) *PushNotification {
	pn := &PushNotification{
		ID: fmt.Sprintf("pn_%d", time.Now().UnixNano()),
		UserID: userID, Title: title, Body: body,
		Data: data, Status: "pending", SentAt: time.Now(),
	}
	m.mu.Lock()
	m.pushQueue = append(m.pushQueue, *pn)
	m.mu.Unlock()
	return pn
}

func (m *Manager) GetPushQueue() []PushNotification {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.pushQueue
}

func (m *Manager) CreateSession(userID, deviceID string, ttl time.Duration) *Session {
	s := &Session{
		ID: fmt.Sprintf("sess_%d", time.Now().UnixNano()),
		UserID: userID, DeviceID: deviceID,
		ExpiresAt: time.Now().Add(ttl),
		CreatedAt: time.Now(), LastActivity: time.Now(),
	}
	m.mu.Lock()
	m.sessions[s.ID] = s
	m.mu.Unlock()
	return s
}

func (m *Manager) ValidateSession(sessionID string) (*Session, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.sessions[sessionID]
	if !ok { return nil, fmt.Errorf("session not found") }
	if time.Now().After(s.ExpiresAt) {
		return nil, fmt.Errorf("session expired")
	}
	s.LastActivity = time.Now()
	return s, nil
}

func (m *Manager) GetStatus() map[string]any {
	m.mu.RLock()
	defer m.mu.RUnlock()

	pendingSync := 0
	for _, ch := range m.offlineQueue {
		pendingSync += len(ch)
	}

	return map[string]any{
		"devices_registered": len(m.devices),
		"active_sessions":    len(m.sessions),
		"pending_sync":       pendingSync,
		"push_queue":         len(m.pushQueue),
		"sync_results":       len(m.syncResults),
	}
}
