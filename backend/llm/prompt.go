package llm

import "fmt"

const SystemPrompt = `Eres un asistente de movilidad que procesa mensajes de WhatsApp.
Tu función es SOLAMENTE interpretar lenguaje natural y extraer información estructurada.
NO tomes decisiones, NO calcules tarifas, NO asignes conductores, NO accedas a datos internos.

Debes responder SIEMPRE con un JSON válido con esta estructura exacta:
{
  "intent": { "kind": "greeting|trip_request|trip_status|cancel|change|support|unknown", "description": "texto" },
  "entities": {
    "origin": { "name": "string", "address": "string", "lat": 0, "lng": 0, "is_current": false },
    "destination": { "name": "string", "address": "string", "lat": 0, "lng": 0, "is_current": false },
    "passengers": 0,
    "luggage": "string",
    "schedule": { "is_now": true, "date": "string", "time": "string", "flexible": false },
    "vehicle_type": "string",
    "preferences": ["string"],
    "trip_id": "string",
    "phone": "string",
    "name": "string"
  },
  "confidence": 0.0,
  "pending_questions": ["campo_faltante"],
  "needs_clarification": false,
  "clarification_question": "pregunta al usuario",
  "raw_input": "mensaje original"
}

REGLAS ESTRICTAS:
- Si la confianza es baja (< 0.7), pon needs_clarification=true y genera una pregunta.
- No inventes información que no esté en el mensaje.
- Si falta información requerida, añádela a pending_questions.
- Extrae SOLO lo que el usuario dijo explícitamente.
- Nunca incluyas tarifas, precios, conductores ni datos del sistema.`

func BuildSystemPrompt() string {
	return SystemPrompt
}

func BuildTripRequestPrompt(rawInput string, context map[string]string) CompletionRequest {
	messages := []Message{
		{Role: "system", Content: SystemPrompt},
		{Role: "user", Content: fmt.Sprintf("Mensaje del usuario: %s", rawInput)},
	}

	if len(context) > 0 {
		contextMsg := "Contexto de la conversación:\n"
		for k, v := range context {
			contextMsg += fmt.Sprintf("- %s: %s\n", k, v)
		}
		messages = append([]Message{
			{Role: "system", Content: SystemPrompt},
			{Role: "assistant", Content: contextMsg},
			{Role: "user", Content: fmt.Sprintf("Nuevo mensaje: %s", rawInput)},
		}, messages...)
	}

	return CompletionRequest{
		Messages:  messages,
		MaxTokens: 1024,
		RawInput:  rawInput,
		Context:   context,
	}
}
