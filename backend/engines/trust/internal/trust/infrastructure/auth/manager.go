package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type Claims struct {
	Sub   string `json:"sub"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
	Name  string `json:"name,omitempty"`
	Role  string `json:"role"`
	Iat   int64  `json:"iat"`
	Exp   int64  `json:"exp"`
}

type Session struct {
	UserID     string    `json:"user_id"`
	DeviceID   string    `json:"device_id"`
	IP         string    `json:"ip"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiresAt  time.Time `json:"expires_at"`
	RefreshToken string  `json:"refresh_token"`
}

type Manager struct {
	secret       string
	accessTTL    time.Duration
	refreshTTL   time.Duration
	sessions     sync.Map // token → Session
	blacklist     sync.Map // jti → time.Time
}

func NewManager(secret string, accessTTL, refreshTTL time.Duration) *Manager {
	return &Manager{
		secret:     secret,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (m *Manager) GenerateTokens(claims Claims) (*TokenPair, error) {
	now := time.Now()

	accessClaims := claims
	accessClaims.Iat = now.Unix()
	accessClaims.Exp = now.Add(m.accessTTL).Unix()

	accessToken, err := m.sign(accessClaims)
	if err != nil {
		return nil, err
	}

	refreshClaims := claims
	refreshClaims.Iat = now.Unix()
	refreshClaims.Exp = now.Add(m.refreshTTL).Unix()

	refreshToken, err := m.sign(refreshClaims)
	if err != nil {
		return nil, err
	}

	m.sessions.Store(refreshToken, Session{
		UserID: claims.Sub, CreatedAt: now,
		ExpiresAt: now.Add(m.refreshTTL), RefreshToken: refreshToken,
	})

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(m.accessTTL.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

func (m *Manager) ValidateToken(tokenString string) (*Claims, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid payload: %w", err)
	}

	sig := parts[2]
	mac := hmac.New(sha256.New, []byte(m.secret))
	mac.Write([]byte(parts[0] + "." + parts[1]))
	expected := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(sig), []byte(expected)) {
		return nil, fmt.Errorf("invalid signature")
	}

	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, fmt.Errorf("invalid claims: %w", err)
	}

	if time.Now().Unix() > claims.Exp {
		return nil, fmt.Errorf("token expired")
	}

	return &claims, nil
}

func (m *Manager) RefreshToken(refreshToken string) (*TokenPair, error) {
	claims, err := m.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	_, ok := m.sessions.Load(refreshToken)
	if !ok {
		return nil, fmt.Errorf("session not found or revoked")
	}

	return m.GenerateTokens(*claims)
}

func (m *Manager) RevokeSession(refreshToken string) {
	m.sessions.Delete(refreshToken)
}

func (m *Manager) BlacklistToken(jti string) {
	m.blacklist.Store(jti, time.Now())
}

func (m *Manager) sign(claims Claims) (string, error) {
	header := map[string]string{"alg": "HS256", "typ": "JWT"}
	headerJSON, _ := json.Marshal(header)
	headerB64 := base64.RawURLEncoding.EncodeToString(headerJSON)

	claimsJSON, _ := json.Marshal(claims)
	claimsB64 := base64.RawURLEncoding.EncodeToString(claimsJSON)

	mac := hmac.New(sha256.New, []byte(m.secret))
	mac.Write([]byte(headerB64 + "." + claimsB64))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return headerB64 + "." + claimsB64 + "." + sig, nil
}

type RoleManager struct {
	roles map[string][]string
}

func NewRoleManager() *RoleManager {
	return &RoleManager{
		roles: map[string][]string{
			"customer": {"trip:create", "trip:read", "profile:read", "profile:write", "payment:read", "payment:create"},
			"driver":  {"trip:read", "trip:accept", "trip:reject", "trip:start", "trip:finish", "profile:read", "vehicle:manage", "availability:manage"},
			"operator": {"trip:read", "trip:cancel", "driver:read", "payment:read", "refund:create", "support:read"},
			"admin":   {"*"},
		},
	}
}

func (rm *RoleManager) Can(role, resource, action string) bool {
	perms, ok := rm.roles[role]
	if !ok { return false }
	needle := resource + ":" + action
	for _, p := range perms {
		if p == "*" || p == needle {
			return true
		}
	}
	return false
}
