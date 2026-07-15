package providers

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type MetaCloudAdapter struct {
	name       string
	apiToken   string
	phoneID    string
	apiURL     string
}

func NewMetaCloudAdapter(apiToken, phoneID string) *MetaCloudAdapter {
	return &MetaCloudAdapter{
		name:    "meta_cloud",
		apiToken: apiToken,
		phoneID:  phoneID,
		apiURL:   "https://graph.facebook.com/v18.0",
	}
}

func (a *MetaCloudAdapter) Name() string { return a.name }

func (a *MetaCloudAdapter) IsAvailable(ctx context.Context) bool {
	return a.apiToken != "" && a.phoneID != ""
}

func (a *MetaCloudAdapter) SendMessage(ctx context.Context, msg Message) (*SendResponse, error) {
	if a.apiToken == "" {
		return nil, fmt.Errorf("meta: API token not configured (set WHATSAPP_API_TOKEN)")
	}
	return nil, fmt.Errorf("meta: real HTTP call requires active token")
}

func (a *MetaCloudAdapter) SendTemplate(ctx context.Context, to string, template TemplateBody) (*SendResponse, error) {
	if a.apiToken == "" {
		return nil, fmt.Errorf("meta: API token not configured")
	}
	return nil, fmt.Errorf("meta: real HTTP call requires active token")
}

func (a *MetaCloudAdapter) MarkAsRead(ctx context.Context, messageID string) error {
	return fmt.Errorf("meta: API token not configured")
}

type MockProvider struct {
	name     string
	messages sync.Map
}

func NewMockProvider() *MockProvider { return &MockProvider{name: "mock"} }

func (m *MockProvider) Name() string { return m.name }
func (m *MockProvider) IsAvailable(ctx context.Context) bool { return true }

func (m *MockProvider) SendMessage(ctx context.Context, msg Message) (*SendResponse, error) {
	id := fmt.Sprintf("wa_%d", time.Now().UnixNano())
	m.messages.Store(id, msg)
	return &SendResponse{MessageID: id, Status: "sent"}, nil
}

func (m *MockProvider) SendTemplate(ctx context.Context, to string, template TemplateBody) (*SendResponse, error) {
	id := fmt.Sprintf("wa_tpl_%d", time.Now().UnixNano())
	msg := Message{ID: id, From: "bot", To: to, Type: MsgTemplate, Template: &template, Timestamp: time.Now()}
	m.messages.Store(id, msg)
	return &SendResponse{MessageID: id, Status: "sent"}, nil
}

func (m *MockProvider) MarkAsRead(ctx context.Context, messageID string) error { return nil }

func (m *MockProvider) GetMessage(id string) (Message, bool) {
	v, ok := m.messages.Load(id)
	if !ok { return Message{}, false }
	return v.(Message), true
}
