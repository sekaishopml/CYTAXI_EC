package providers

import (
	"context"
	"fmt"
	"time"
)

type PaymentIntentRequest struct {
	Amount      float64
	Currency    string
	Description string
	Metadata    map[string]string
}

type PaymentIntentResponse struct {
	ID            string
	ClientSecret  string
	Status        string
	ProviderRef   string
	CreatedAt     time.Time
}

type AuthorizationRequest struct {
	IntentID   string
	Amount     float64
	Currency   string
	PaymentMethod string // card, wallet, transfer
}

type AuthorizationResponse struct {
	ID           string
	Status       string
	GatewayRef   string
	AuthCode     string
	Message      string
}

type CaptureRequest struct {
	AuthorizationID string
	Amount          float64
}

type CaptureResponse struct {
	ID         string
	Status     string
	GatewayRef string
}

type RefundRequest struct {
	PaymentID string
	Amount    float64
	Reason    string
}

type RefundResponse struct {
	ID         string
	Status     string
	GatewayRef string
}

type Provider interface {
	Name() string
	CreatePaymentIntent(ctx context.Context, req PaymentIntentRequest) (*PaymentIntentResponse, error)
	Authorize(ctx context.Context, req AuthorizationRequest) (*AuthorizationResponse, error)
	Capture(ctx context.Context, req CaptureRequest) (*CaptureResponse, error)
	Refund(ctx context.Context, req RefundRequest) (*RefundResponse, error)
	IsAvailable(ctx context.Context) bool
}

type ProviderRegistry struct {
	providers map[string]Provider
}

func NewRegistry() *ProviderRegistry {
	return &ProviderRegistry{providers: make(map[string]Provider)}
}

func (r *ProviderRegistry) Register(p Provider) {
	r.providers[p.Name()] = p
}

func (r *ProviderRegistry) Get(name string) (Provider, error) {
	p, ok := r.providers[name]
	if !ok {
		return nil, fmt.Errorf("provider %s not found", name)
	}
	return p, nil
}

func (r *ProviderRegistry) List() []string {
	names := make([]string, 0, len(r.providers))
	for name := range r.providers {
		names = append(names, name)
	}
	return names
}

type StripeAdapter struct{ name string }
func NewStripeAdapter() Provider     { return &StripeAdapter{name: "stripe"} }
func (a *StripeAdapter) Name() string { return a.name }
func (a *StripeAdapter) IsAvailable(ctx context.Context) bool { return false }
func (a *StripeAdapter) CreatePaymentIntent(ctx context.Context, req PaymentIntentRequest) (*PaymentIntentResponse, error) {
	return nil, fmt.Errorf("stripe: not configured (set STRIPE_API_KEY)")
}
func (a *StripeAdapter) Authorize(ctx context.Context, req AuthorizationRequest) (*AuthorizationResponse, error) {
	return nil, fmt.Errorf("stripe: not configured")
}
func (a *StripeAdapter) Capture(ctx context.Context, req CaptureRequest) (*CaptureResponse, error) {
	return nil, fmt.Errorf("stripe: not configured")
}
func (a *StripeAdapter) Refund(ctx context.Context, req RefundRequest) (*RefundResponse, error) {
	return nil, fmt.Errorf("stripe: not configured")
}

type KushkiAdapter struct{ name string }
func NewKushkiAdapter() Provider     { return &KushkiAdapter{name: "kushki"} }
func (a *KushkiAdapter) Name() string { return a.name }
func (a *KushkiAdapter) IsAvailable(ctx context.Context) bool { return false }
func (a *KushkiAdapter) CreatePaymentIntent(ctx context.Context, req PaymentIntentRequest) (*PaymentIntentResponse, error) {
	return nil, fmt.Errorf("kushki: not configured (set KUSHKI_API_KEY)")
}
func (a *KushkiAdapter) Authorize(ctx context.Context, req AuthorizationRequest) (*AuthorizationResponse, error) {
	return nil, fmt.Errorf("kushki: not configured")
}
func (a *KushkiAdapter) Capture(ctx context.Context, req CaptureRequest) (*CaptureResponse, error) {
	return nil, fmt.Errorf("kushki: not configured")
}
func (a *KushkiAdapter) Refund(ctx context.Context, req RefundRequest) (*RefundResponse, error) {
	return nil, fmt.Errorf("kushki: not configured")
}

type PayPhoneAdapter struct{ name string }
func NewPayPhoneAdapter() Provider     { return &PayPhoneAdapter{name: "payphone"} }
func (a *PayPhoneAdapter) Name() string { return a.name }
func (a *PayPhoneAdapter) IsAvailable(ctx context.Context) bool { return false }
func (a *PayPhoneAdapter) CreatePaymentIntent(ctx context.Context, req PaymentIntentRequest) (*PaymentIntentResponse, error) {
	return nil, fmt.Errorf("payphone: not configured (set PAYPHONE_API_KEY)")
}
func (a *PayPhoneAdapter) Authorize(ctx context.Context, req AuthorizationRequest) (*AuthorizationResponse, error) {
	return nil, fmt.Errorf("payphone: not configured")
}
func (a *PayPhoneAdapter) Capture(ctx context.Context, req CaptureRequest) (*CaptureResponse, error) {
	return nil, fmt.Errorf("payphone: not configured")
}
func (a *PayPhoneAdapter) Refund(ctx context.Context, req RefundRequest) (*RefundResponse, error) {
	return nil, fmt.Errorf("payphone: not configured")
}

type MockProvider struct {
	name string
	failOn string
}

func NewMockProvider() Provider {
	return &MockProvider{name: "mock"}
}

func (a *MockProvider) Name() string { return a.name }
func (a *MockProvider) IsAvailable(ctx context.Context) bool { return true }

func (a *MockProvider) CreatePaymentIntent(ctx context.Context, req PaymentIntentRequest) (*PaymentIntentResponse, error) {
	id := fmt.Sprintf("pi_%d", time.Now().UnixNano())
	return &PaymentIntentResponse{
		ID: id, ClientSecret: id + "_secret",
		Status: "requires_payment_method", ProviderRef: "mock_" + id,
		CreatedAt: time.Now(),
	}, nil
}

func (a *MockProvider) Authorize(ctx context.Context, req AuthorizationRequest) (*AuthorizationResponse, error) {
	id := fmt.Sprintf("auth_%d", time.Now().UnixNano())
	return &AuthorizationResponse{
		ID: id, Status: "authorized", GatewayRef: "gw_" + id, AuthCode: "123456",
		Message: "Authorization successful",
	}, nil
}

func (a *MockProvider) Capture(ctx context.Context, req CaptureRequest) (*CaptureResponse, error) {
	id := fmt.Sprintf("cap_%d", time.Now().UnixNano())
	return &CaptureResponse{ID: id, Status: "captured", GatewayRef: "gw_" + id}, nil
}

func (a *MockProvider) Refund(ctx context.Context, req RefundRequest) (*RefundResponse, error) {
	id := fmt.Sprintf("ref_%d", time.Now().UnixNano())
	return &RefundResponse{ID: id, Status: "refunded", GatewayRef: "gw_" + id}, nil
}
