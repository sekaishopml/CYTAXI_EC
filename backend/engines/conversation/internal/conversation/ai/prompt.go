package ai

import (
	"fmt"
	"strings"
)

type PromptTemplate struct {
	SystemMessage string
	Template      string
	MaxTokens     int
}

type BuiltPrompt struct {
	System string
	User   string
	Full   string
}

func BuildPrompt(intent *Intent, ctx Context) CompletionRequest {
	system := buildSystemMessage(intent)
	user := buildUserMessage(intent, ctx)

	return CompletionRequest{
		Prompt:    user,
		SystemMsg: system,
		Messages: []Message{
			{Role: "system", Content: system},
			{Role: "user", Content: user},
		},
		MaxTokens: 1024,
	}
}

func buildSystemMessage(intent *Intent) string {
	parts := []string{
		"Eres un asistente de movilidad.",
		"Responde de forma clara y profesional.",
		"No inventes información.",
		"Si no sabes la respuesta, solicita más datos al usuario.",
	}

	switch intent.Kind {
	case IntentTripRequest:
		parts = append(parts, "Ayuda al usuario a solicitar un viaje.")
	case IntentTripStatus:
		parts = append(parts, "Proporciona el estado del viaje solicitado.")
	case IntentSupport:
		parts = append(parts, "Ofrece ayuda al usuario con su problema.")
	case IntentCancel:
		parts = append(parts, "Ayuda al usuario a cancelar el viaje.")
	case IntentGreeting:
		parts = append(parts, "Saluda al usuario y pregunta cómo puedes ayudarle.")
	}

	return strings.Join(parts, "\n")
}

func buildUserMessage(intent *Intent, ctx Context) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Intención: %s\n", intent.Kind))

	if len(ctx.Entries) > 0 {
		b.WriteString("Contexto de la conversación:\n")
		for k, v := range ctx.Entries {
			b.WriteString(fmt.Sprintf("  %s: %s\n", k, v))
		}
	}

	if len(intent.Entities) > 0 {
		b.WriteString("Entidades detectadas:\n")
		for k, v := range intent.Entities {
			b.WriteString(fmt.Sprintf("  %s: %s\n", k, v))
		}
	}

	b.WriteString("Genera una respuesta apropiada.")
	return b.String()
}
