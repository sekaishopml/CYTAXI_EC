package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type claims struct {
	Sub   string   `json:"sub"`
	Role  string   `json:"role"`
	Scope []string `json:"scope"`
	Exp   int64    `json:"exp"`
	Iat   int64    `json:"iat"`
}

type jwtAuthenticator struct {
	secret []byte
}

func NewAuthenticator(secret string) Authenticator {
	return &jwtAuthenticator{secret: []byte(secret)}
}

func (a *jwtAuthenticator) Authenticate(token string) (*Principal, error) {
	if token == "" {
		return nil, fmt.Errorf("auth: empty token")
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("auth: invalid token format")
	}

	mac := hmac.New(sha256.New, a.secret)
	mac.Write([]byte(parts[0] + "." + parts[1]))
	expected := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(parts[2]), []byte(expected)) {
		return nil, fmt.Errorf("auth: invalid signature")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("auth: invalid payload: %w", err)
	}

	var c claims
	if err := json.Unmarshal(payload, &c); err != nil {
		return nil, fmt.Errorf("auth: invalid claims: %w", err)
	}

	if time.Now().Unix() > c.Exp {
		return nil, fmt.Errorf("auth: token expired")
	}

	return &Principal{
		ID:    c.Sub,
		Role:  c.Role,
		Scope: c.Scope,
	}, nil
}

func GenerateToken(secret string, principal *Principal, expiry time.Duration) (string, error) {
	h := header{Alg: "HS256", Typ: "JWT"}
	headerBytes, _ := json.Marshal(h)
	headerStr := base64.RawURLEncoding.EncodeToString(headerBytes)

	now := time.Now()
	c := claims{
		Sub:   principal.ID,
		Role:  principal.Role,
		Scope: principal.Scope,
		Iat:   now.Unix(),
		Exp:   now.Add(expiry).Unix(),
	}
	claimsBytes, _ := json.Marshal(c)
	claimsStr := base64.RawURLEncoding.EncodeToString(claimsBytes)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(headerStr + "." + claimsStr))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return headerStr + "." + claimsStr + "." + sig, nil
}

type roleAuthorizer struct {
	rules map[string][]string
}

func NewAuthorizer(rules map[string][]string) Authorizer {
	return &roleAuthorizer{rules: rules}
}

func (a *roleAuthorizer) Authorize(principal *Principal, resource, action string) bool {
	allowed, ok := a.rules[principal.Role]
	if !ok {
		return false
	}
	for _, a := range allowed {
		if a == resource+":"+action || a == "*" {
			return true
		}
	}
	return false
}
