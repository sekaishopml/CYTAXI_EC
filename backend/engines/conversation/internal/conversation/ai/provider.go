package ai

import (
	"context"
	"errors"
	"math/rand"
)

type ProviderKind string

const (
	ProviderQwen    ProviderKind = "qwen"
	ProviderDeepSeek ProviderKind = "deepseek"
	ProviderGPT     ProviderKind = "gpt"
	ProviderClaude  ProviderKind = "claude"
)

type CompletionRequest struct {
	Prompt     string
	MaxTokens  int
	SystemMsg  string
	Messages   []Message
}

type Message struct {
	Role    string
	Content string
}

type CompletionResponse struct {
	Content    string
	Confidence float64
	TokensUsed int
}

type Provider interface {
	Name() string
	Kind() ProviderKind
	Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error)
}

type ProviderRegistry struct {
	providers []Provider
}

func NewProviderRegistry(providers ...Provider) ProviderRegistry {
	return ProviderRegistry{providers: providers}
}

func (r *ProviderRegistry) Register(p Provider) {
	r.providers = append(r.providers, p)
}

func (r *ProviderRegistry) Select(capabilities ...string) (Provider, error) {
	if len(r.providers) == 0 {
		return nil, errors.New("ai: no providers registered")
	}
	return r.providers[0], nil
}

type mockProvider struct {
	name   string
	kind   ProviderKind
}

func NewMockProvider(kind ProviderKind) Provider {
	return &mockProvider{
		name:   string(kind),
		kind:   kind,
	}
}

func (p *mockProvider) Name() string { return p.name }
func (p *mockProvider) Kind() ProviderKind { return p.kind }

func (p *mockProvider) Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	return &CompletionResponse{
		Content:    "mock response for: " + req.Prompt[:min(len(req.Prompt), 50)],
		Confidence: rand.Float64()*0.3 + 0.7,
		TokensUsed: len(req.Prompt) / 4,
	}, nil
}
