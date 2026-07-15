package providers

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type AuthMethod string

const (
	MethodEmailPassword AuthMethod = "email_password"
	MethodOTP           AuthMethod = "otp"
	MethodGoogleOAuth   AuthMethod = "google_oauth"
	MethodAppleSignIn   AuthMethod = "apple_signin"
	MethodWhatsAppOTP   AuthMethod = "whatsapp_otp"
)

type Credential struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Phone    string `json:"phone,omitempty"`
	OTP      string `json:"otp,omitempty"`
	Token    string `json:"token,omitempty"`
}

type Identity struct {
	ID       string `json:"id"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	Provider string `json:"provider"`
}

type AuthResult struct {
	Success bool
	Identity *Identity
	Error    error
}

type Provider interface {
	Name() string
	Register(ctx context.Context, cred Credential) (*AuthResult, error)
	Authenticate(ctx context.Context, cred Credential) (*AuthResult, error)
	VerifyMFA(ctx context.Context, code string) (*AuthResult, error)
	ResetPassword(ctx context.Context, email string) error
	IsAvailable(ctx context.Context) bool
}

type Registry struct {
	providers map[AuthMethod]Provider
}

func NewRegistry() *Registry { return &Registry{providers: make(map[AuthMethod]Provider)} }
func (r *Registry) Register(method AuthMethod, p Provider) { r.providers[method] = p }
func (r *Registry) Get(method AuthMethod) (Provider, error) {
	p, ok := r.providers[method]
	if !ok { return nil, fmt.Errorf("auth method %s not found", method) }
	return p, nil
}

type EmailPasswordProvider struct {
	users sync.Map // email → hashedPassword
}

func NewEmailPasswordProvider() *EmailPasswordProvider { return &EmailPasswordProvider{} }

func (p *EmailPasswordProvider) Name() string { return "email_password" }
func (p *EmailPasswordProvider) IsAvailable(ctx context.Context) bool { return true }

func (p *EmailPasswordProvider) Register(ctx context.Context, cred Credential) (*AuthResult, error) {
	if cred.Email == "" || cred.Password == "" {
		return &AuthResult{Success: false, Error: fmt.Errorf("email and password required")}, nil
	}
	p.users.Store(cred.Email, hashPassword(cred.Password))
	id := fmt.Sprintf("usr_%d", time.Now().UnixNano())
	return &AuthResult{Success: true, Identity: &Identity{
		ID: id, Email: cred.Email, Role: "customer", Provider: "email_password",
	}}, nil
}

func (p *EmailPasswordProvider) Authenticate(ctx context.Context, cred Credential) (*AuthResult, error) {
	v, ok := p.users.Load(cred.Email)
	if !ok || v.(string) != hashPassword(cred.Password) {
		return &AuthResult{Success: false, Error: fmt.Errorf("invalid credentials")}, nil
	}
	id := fmt.Sprintf("usr_%d", time.Now().UnixNano())
	return &AuthResult{Success: true, Identity: &Identity{
		ID: id, Email: cred.Email, Role: "customer", Provider: "email_password",
	}}, nil
}

func (p *EmailPasswordProvider) VerifyMFA(ctx context.Context, code string) (*AuthResult, error) {
	return nil, fmt.Errorf("not implemented")
}
func (p *EmailPasswordProvider) ResetPassword(ctx context.Context, email string) error { return nil }

type OTPProvider struct{}

func NewOTPProvider() *OTPProvider { return &OTPProvider{} }
func (p *OTPProvider) Name() string { return "otp" }
func (p *OTPProvider) IsAvailable(ctx context.Context) bool { return true }
func (p *OTPProvider) Register(ctx context.Context, cred Credential) (*AuthResult, error) {
	return &AuthResult{Success: true, Identity: &Identity{
		ID: fmt.Sprintf("usr_otp_%d", time.Now().UnixNano()), Phone: cred.Phone, Role: "customer", Provider: "otp",
	}}, nil
}
func (p *OTPProvider) Authenticate(ctx context.Context, cred Credential) (*AuthResult, error) {
	if cred.OTP != "" && cred.OTP == "123456" {
		return &AuthResult{Success: true, Identity: &Identity{
			ID: fmt.Sprintf("usr_otp_%d", time.Now().UnixNano()), Phone: cred.Phone, Role: "customer", Provider: "otp",
		}}, nil
	}
	return &AuthResult{Success: false, Error: fmt.Errorf("invalid OTP")}, nil
}
func (p *OTPProvider) VerifyMFA(ctx context.Context, code string) (*AuthResult, error) { return nil, nil }
func (p *OTPProvider) ResetPassword(ctx context.Context, email string) error { return nil }

func hashPassword(pwd string) string { return fmt.Sprintf("hashed_%s", pwd) }
