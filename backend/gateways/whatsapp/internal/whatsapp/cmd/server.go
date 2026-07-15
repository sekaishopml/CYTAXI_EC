package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/sekaishopml/cytaxi/backend/gateways/whatsapp/internal/whatsapp/providers"
	"github.com/sekaishopml/cytaxi/backend/gateways/whatsapp/internal/whatsapp/webhook"
)

type WABServer struct {
	registry *providers.Registry
	receiver *webhook.Receiver
	logger   *slog.Logger
	convStore sync.Map
}

func NewWABServer(registry *providers.Registry, receiver *webhook.Receiver, logger *slog.Logger) *WABServer {
	s := &WABServer{registry: registry, receiver: receiver, logger: logger}
	s.setupHandlers()
	return s
}

func (s *WABServer) setupHandlers() {
	s.receiver.Handle(providers.MsgText, func(ctx *context.Context, msg providers.Message) error {
		s.logger.Info("whatsapp message received", "from", msg.From, "type", msg.Type)
		s.storeConversation(msg.From, msg)
		return nil
	})
	s.receiver.Handle(providers.MsgLocation, func(ctx *context.Context, msg providers.Message) error {
		s.logger.Info("whatsapp location received", "from", msg.From)
		s.storeConversation(msg.From, msg)
		return nil
	})
	s.receiver.Handle(providers.MsgInteractive, func(ctx *context.Context, msg providers.Message) error {
		s.logger.Info("whatsapp interactive received", "from", msg.From)
		s.storeConversation(msg.From, msg)
		return nil
	})
}

func (s *WABServer) storeConversation(phone string, msg providers.Message) {
	existing, _ := s.convStore.Load(phone)
	var history []providers.Message
	if existing != nil {
		history = existing.([]providers.Message)
	}
	history = append(history, msg)
	if len(history) > 100 {
		history = history[1:]
	}
	s.convStore.Store(phone, history)
}

func (s *WABServer) HandleSendMessage(w http.ResponseWriter, r *http.Request) {
	var req struct {
		To   string `json:"to"`
		Type string `json:"type"`
		Text string `json:"text,omitempty"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	provider, _ := s.registry.Get("mock")
	msg := providers.Message{
		From: "bot", To: req.To, Type: providers.MsgText,
		Text: &providers.TextBody{Body: req.Text}, Timestamp: time.Now(),
	}

	resp, err := provider.SendMessage(r.Context(), msg)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	s.logger.Info("whatsapp message sent", "to", req.To)
	writeJSON(w, http.StatusOK, resp)
}

func (s *WABServer) HandleSendTemplate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		To       string `json:"to"`
		Template string `json:"template"`
		Language string `json:"language"`
		Params   []providers.TemplateParameter `json:"params"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	provider, _ := s.registry.Get("mock")
	tpl := providers.TemplateBody{
		Name: req.Template, Language: req.Language,
		Components: []providers.TemplateComponent{
			{Type: "body", Parameters: req.Params},
		},
	}

	resp, err := provider.SendTemplate(r.Context(), req.To, tpl)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	s.logger.Info("whatsapp template sent", "to", req.To, "template", req.Template)
	writeJSON(w, http.StatusOK, resp)
}

func (s *WABServer) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	// Verification (GET)
	if r.Method == "GET" {
		mode := r.URL.Query().Get("hub.mode")
		token := r.URL.Query().Get("hub.verify_token")
		challenge := r.URL.Query().Get("hub.challenge")
		if resp, ok := s.receiver.VerifyWebhookToken(mode, token, challenge); ok {
			w.Write([]byte(resp))
			return
		}
		http.Error(w, "verification failed", http.StatusForbidden)
		return
	}

	// Process webhook (POST)
	var evt providers.WebhookEvent
	if err := json.NewDecoder(r.Body).Decode(&evt); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid payload"})
		return
	}

	payload, _ := json.Marshal(evt)
	signature := r.Header.Get("X-Hub-Signature-256")

	if !s.receiver.VerifySignature(payload, signature) {
		s.logger.Warn("webhook signature invalid")
		writeJSON(w, http.StatusForbidden, map[string]any{"error": "invalid signature"})
		return
	}

	messages, err := s.receiver.ProcessWebhook(r.Context(), payload)
	if err != nil {
		s.logger.Error("webhook process failed", "error", err)
		writeJSON(w, http.StatusOK, map[string]any{"status": "partial", "error": err.Error()})
		return
	}

	s.logger.Info("webhook processed", "messages", len(messages))
	writeJSON(w, http.StatusOK, map[string]any{"status": "received", "messages_processed": len(messages)})
}

func (s *WABServer) HandleGetConversation(w http.ResponseWriter, r *http.Request) {
	phone := r.PathValue("phone")
	conv, _ := s.convStore.Load(phone)
	if conv == nil {
		writeJSON(w, http.StatusOK, map[string]any{"messages": []providers.Message{}, "phone": phone})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"phone":    phone,
		"messages": conv,
		"count":    len(conv.([]providers.Message)),
	})
}

func (s *WABServer) HandleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":  "ok",
		"service": "whatsapp-gateway",
		"provider": fmt.Sprintf("ready"),
	})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
