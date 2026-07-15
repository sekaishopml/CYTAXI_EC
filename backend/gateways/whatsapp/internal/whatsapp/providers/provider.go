package providers

import (
	"context"
	"fmt"
	"time"
)

type MessageType string

const (
	MsgText        MessageType = "text"
	MsgImage       MessageType = "image"
	MsgAudio       MessageType = "audio"
	MsgVideo       MessageType = "video"
	MsgDocument    MessageType = "document"
	MsgLocation    MessageType = "location"
	MsgInteractive MessageType = "interactive"
	MsgTemplate    MessageType = "template"
	MsgButton      MessageType = "button"
	MsgList        MessageType = "list"
)

type Message struct {
	ID       string      `json:"id,omitempty"`
	From     string      `json:"from"`
	To       string      `json:"to"`
	Type     MessageType `json:"type"`
	Text     *TextBody   `json:"text,omitempty"`
	Location *LocationBody `json:"location,omitempty"`
	Image    *MediaBody  `json:"image,omitempty"`
	Document *MediaBody  `json:"document,omitempty"`
	Audio    *MediaBody  `json:"audio,omitempty"`
	Template *TemplateBody `json:"template,omitempty"`
	Interactive *InteractiveBody `json:"interactive,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type TextBody struct {
	Body string `json:"body"`
}

type LocationBody struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      string  `json:"name,omitempty"`
	Address   string  `json:"address,omitempty"`
}

type MediaBody struct {
	ID       string `json:"id,omitempty"`
	Link     string `json:"link"`
	Caption  string `json:"caption,omitempty"`
	Filename string `json:"filename,omitempty"`
}

type TemplateBody struct {
	Name       string            `json:"name"`
	Language   string            `json:"language"`
	Components []TemplateComponent `json:"components,omitempty"`
}

type TemplateComponent struct {
	Type       string              `json:"type"`
	Parameters []TemplateParameter `json:"parameters,omitempty"`
}

type TemplateParameter struct {
	Type string `json:"type"` // text, currency, date_time, image, document, video
	Text string `json:"text,omitempty"`
}

type InteractiveBody struct {
	Type   string              `json:"type"` // button, list
	Header *InteractiveHeader  `json:"header,omitempty"`
	Body   *TextBody           `json:"body"`
	Footer *TextBody           `json:"footer,omitempty"`
	Action *InteractiveAction   `json:"action"`
}

type InteractiveHeader struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

type InteractiveAction struct {
	Buttons []InteractiveButton `json:"buttons,omitempty"`
	Button  string              `json:"button,omitempty"`
	Sections []ListSection      `json:"sections,omitempty"`
}

type InteractiveButton struct {
	Type string `json:"type"`
	Reply Reply `json:"reply"`
}

type Reply struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type ListSection struct {
	Title string     `json:"title,omitempty"`
	Rows  []ListRow  `json:"rows"`
}

type ListRow struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

type SendResponse struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
}

type WebhookEvent struct {
	ID      string           `json:"id"`
	Changes []WebhookChange  `json:"entry"`
}

type WebhookChange struct {
	ID      string            `json:"id"`
	Changes []WebhookFieldChange `json:"changes"`
}

type WebhookFieldChange struct {
	Field string      `json:"field"`
	Value WebhookValue `json:"value"`
}

type WebhookValue struct {
	MessagingProduct string     `json:"messaging_product"`
	Metadata         WebhookMeta `json:"metadata"`
	Contacts         []WebhookContact `json:"contacts,omitempty"`
	Messages         []WebhookMessage `json:"messages,omitempty"`
	Statuses         []WebhookStatus  `json:"statuses,omitempty"`
}

type WebhookMeta struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberID      string `json:"phone_number_id"`
}

type WebhookContact struct {
	Profile struct {
		Name string `json:"name"`
	} `json:"profile"`
	WaID string `json:"wa_id"`
}

type WebhookMessage struct {
	ID        string `json:"id"`
	From      string `json:"from"`
	Timestamp string `json:"timestamp"`
	Type      string `json:"type"`
	Text      *struct{ Body string } `json:"text,omitempty"`
	Location  *LocationBody `json:"location,omitempty"`
	Image     *MediaBody `json:"image,omitempty"`
	Interactive *struct{ Type string `json:"type"`; ButtonReply *struct{ID string `json:"id"`; Title string `json:"title"`} `json:"button_reply,omitempty"`; ListReply *struct{ID string `json:"id"`} `json:"list_reply,omitempty"`} `json:"interactive,omitempty"`
}

type WebhookStatus struct {
	ID           string `json:"id"`
	Status       string `json:"status"` // sent, delivered, read, failed
	Timestamp    string `json:"timestamp"`
	RecipientID  string `json:"recipient_id"`
	Conversation struct {
		ID string `json:"id"`
	} `json:"conversation,omitempty"`
}

// Provider interface
type Provider interface {
	Name() string
	SendMessage(ctx context.Context, msg Message) (*SendResponse, error)
	SendTemplate(ctx context.Context, to string, template TemplateBody) (*SendResponse, error)
	MarkAsRead(ctx context.Context, messageID string) error
	IsAvailable(ctx context.Context) bool
}

type Receiver interface {
	ProcessWebhook(ctx context.Context, payload []byte, signature string) ([]Message, error)
}

// Registry
type Registry struct {
	providers map[string]Provider
}

func NewRegistry() *Registry { return &Registry{providers: make(map[string]Provider)} }
func (r *Registry) Register(p Provider) { r.providers[p.Name()] = p }
func (r *Registry) Get(name string) (Provider, error) {
	p, ok := r.providers[name]
	if !ok { return nil, fmt.Errorf("provider %s not found", name) }
	return p, nil
}
