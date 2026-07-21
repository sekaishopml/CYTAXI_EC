package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sekaishopml/cytaxi/llm"
)

type HTTPAdapterConfig struct {
	Endpoint    string
	APIKey      string
	Model       string
	Provider    llm.ProviderKind
	MaxTokens   int
	Temperature float64
	Timeout     time.Duration
}

type HTTPAdapter struct {
	config HTTPAdapterConfig
	client *http.Client
}

func NewHTTPAdapter(config HTTPAdapterConfig) *HTTPAdapter {
	if config.MaxTokens == 0 {
		config.MaxTokens = 2048
	}
	if config.Temperature == 0 {
		config.Temperature = 0.1
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	return &HTTPAdapter{
		config: config,
		client: &http.Client{Timeout: config.Timeout},
	}
}

func (a *HTTPAdapter) Name() string {
	return fmt.Sprintf("http-%s", a.config.Model)
}

func (a *HTTPAdapter) Kind() llm.ProviderKind {
	return a.config.Provider
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIRequest struct {
	Model       string           `json:"model"`
	Messages    []openAIMessage  `json:"messages"`
	MaxTokens   int              `json:"max_tokens"`
	Temperature float64          `json:"temperature"`
	ResponseFormat *responseFormat `json:"response_format,omitempty"`
}

type responseFormat struct {
	Type string `json:"type"`
}

type openAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`
}

func (a *HTTPAdapter) Complete(ctx context.Context, req llm.CompletionRequest) (*llm.CompletionResponse, error) {
	messages := make([]openAIMessage, 0, len(req.Messages)+1)

	if req.RawInput != "" {
		messages = append(messages, openAIMessage{Role: "user", Content: req.RawInput})
	}
	for _, m := range req.Messages {
		messages = append(messages, openAIMessage{Role: m.Role, Content: m.Content})
	}

	apiReq := openAIRequest{
		Model:       a.config.Model,
		Messages:    messages,
		MaxTokens:   a.config.MaxTokens,
		Temperature: a.config.Temperature,
		ResponseFormat: &responseFormat{Type: "json_object"},
	}

	body, err := json.Marshal(apiReq)
	if err != nil {
		return nil, fmt.Errorf("llm: marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", a.config.Endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("llm: create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+a.config.APIKey)

	resp, err := a.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("llm: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("llm: read response: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("llm: API error %d: %s", resp.StatusCode, string(respBody))
	}

	var apiResp openAIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("llm: parse response: %w", err)
	}

	if len(apiResp.Choices) == 0 {
		return nil, fmt.Errorf("llm: no choices in response")
	}

	return &llm.CompletionResponse{
		Content:    apiResp.Choices[0].Message.Content,
		Confidence: 0.85,
		TokensUsed: apiResp.Usage.TotalTokens,
		Provider:   a.config.Provider,
	}, nil
}
