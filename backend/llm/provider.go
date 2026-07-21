package llm

import "context"

type ProviderKind string

const (
	ProviderRule     ProviderKind = "rule"
	ProviderQwen     ProviderKind = "qwen"
	ProviderDeepSeek ProviderKind = "deepseek"
	ProviderGPT      ProviderKind = "gpt"
	ProviderClaude   ProviderKind = "claude"
)

type CompletionRequest struct {
	Messages   []Message
	MaxTokens  int
	RawInput   string
	SessionID  string
	Context    map[string]string
}

type Message struct {
	Role    string
	Content string
}

type CompletionResponse struct {
	Content    string
	Confidence float64
	TokensUsed int
	Provider   ProviderKind
}

type Provider interface {
	Name() string
	Kind() ProviderKind
	Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error)
}

type ProviderRegistry struct {
	providers []Provider
	default_  Provider
}

func NewProviderRegistry(defaultProvider Provider, additional ...Provider) *ProviderRegistry {
	return &ProviderRegistry{
		providers: append([]Provider{defaultProvider}, additional...),
		default_:  defaultProvider,
	}
}

func (r *ProviderRegistry) Register(p Provider) {
	r.providers = append(r.providers, p)
}

func (r *ProviderRegistry) Default() Provider {
	return r.default_
}

func (r *ProviderRegistry) Select(kind ProviderKind) (Provider, error) {
	for _, p := range r.providers {
		if p.Kind() == kind {
			return p, nil
		}
	}
	return r.default_, nil
}
