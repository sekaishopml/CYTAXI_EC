package adapters

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/sekaishopml/cytaxi/llm"
)

type RuleAdapter struct {
	classifier *IntentClassifier
	extractor  *EntityExtractor
	validator  *OutputValidator
}

func NewRuleAdapter() *RuleAdapter {
	return &RuleAdapter{
		classifier: NewIntentClassifier(),
		extractor:  NewEntityExtractor(),
		validator:  NewOutputValidator(),
	}
}

func (a *RuleAdapter) Name() string {
	return "rule-based"
}

func (a *RuleAdapter) Kind() llm.ProviderKind {
	return llm.ProviderRule
}

func (a *RuleAdapter) Complete(ctx context.Context, req llm.CompletionRequest) (*llm.CompletionResponse, error) {
	response := a.classifier.Classify(req.RawInput, req.Context)

	entities := a.extractor.Extract(req.RawInput, response.Intent.Kind)
	response.Entities = *entities

	if response.HasCompleteTripRequest() {
		response.NeedsClarification = false
		response.PendingQuestions = nil
		response.Confidence = 0.85
	} else if response.Intent.Kind == llm.IntentTripRequest {
		response = a.addMissingFields(response)
	}

	response = a.validator.Validate(response)

	jsonBytes, _ := json.Marshal(response)
	return &llm.CompletionResponse{
		Content:    string(jsonBytes),
		Confidence: response.Confidence,
		TokensUsed: len(req.RawInput) / 4,
		Provider:   llm.ProviderRule,
	}, nil
}

func (a *RuleAdapter) addMissingFields(resp *llm.LLMResponse) *llm.LLMResponse {
	if resp.Entities.Origin == nil {
		resp.NeedsClarification = true
		resp.ClarificationQuestion = "¿Desde dónde te gustaría que te recojan?"
		resp.PendingQuestions = append(resp.PendingQuestions, "origin")
		resp.Confidence = 0.5
		return resp
	}
	if resp.Entities.Destination == nil {
		resp.NeedsClarification = true
		resp.ClarificationQuestion = "¿Cuál es tu destino?"
		resp.PendingQuestions = append(resp.PendingQuestions, "destination")
		resp.Confidence = 0.6
		return resp
	}
	if resp.Entities.Passengers == 0 {
		resp.ClarificationQuestion = "¿Cuántos pasajeros son?"
		resp.PendingQuestions = append(resp.PendingQuestions, "passengers")
		resp.Confidence = 0.45
		return resp
	}
	return resp
}

type IntentClassifier struct {
	rules []IntentRule
}

type IntentRule struct {
	Keywords []string
	Intent   llm.IntentKind
	Priority int
}

func NewIntentClassifier() *IntentClassifier {
	return &IntentClassifier{
		rules: []IntentRule{
			{Keywords: []string{"cancelar", "cancel", "cancelación", "anular", "anula"}, Intent: llm.IntentCancel, Priority: 5},
			{Keywords: []string{"cambio", "cambiar", "modificar", "cambia"}, Intent: llm.IntentChange, Priority: 5},
			{Keywords: []string{"taxi", "viaje", "recoger", "llevar", "trip", "ride", "quiero ir", "necesito", "busco"}, Intent: llm.IntentTripRequest, Priority: 4},
			{Keywords: []string{"dónde está", "dónde", "estado", "llegada", "esperando", "status", "where", "tiempo"}, Intent: llm.IntentTripStatus, Priority: 3},
			{Keywords: []string{"ayuda", "soporte", "problema", "help", "support", "error", "no funciona"}, Intent: llm.IntentSupport, Priority: 2},
			{Keywords: []string{"hola", "buenos días", "buenas tardes", "buenas noches", "buenas", "hey", "hello", "hi", "que tal"}, Intent: llm.IntentGreeting, Priority: 1},
		},
	}
}

func (c *IntentClassifier) Classify(input string, ctx map[string]string) *llm.LLMResponse {
	lower := strings.ToLower(input)
	response := &llm.LLMResponse{
		RawInput: input,
		Intent: llm.Intent{
			Kind:        llm.IntentUnknown,
			Description: "No se pudo determinar la intención",
		},
		Confidence: 0.0,
	}

	var bestRule *IntentRule
	for i := range c.rules {
		rule := c.rules[i]
		for _, kw := range rule.Keywords {
			if strings.Contains(lower, kw) {
				if bestRule == nil || rule.Priority > bestRule.Priority {
					bestRule = &rule
					break
				}
			}
		}
	}

	if bestRule != nil {
		response.Intent.Kind = bestRule.Intent
		response.Confidence = 0.7 + float64(bestRule.Priority)*0.04
		if response.Confidence > 0.95 {
			response.Confidence = 0.95
		}
		switch bestRule.Intent {
		case llm.IntentGreeting:
			response.Intent.Description = "El usuario saluda e inicia una conversación"
		case llm.IntentTripRequest:
			response.Intent.Description = "El usuario solicita un servicio de transporte"
		case llm.IntentTripStatus:
			response.Intent.Description = "El usuario consulta el estado de un viaje"
		case llm.IntentCancel:
			response.Intent.Description = "El usuario desea cancelar un viaje"
		case llm.IntentSupport:
			response.Intent.Description = "El usuario necesita ayuda o soporte"
		case llm.IntentChange:
			response.Intent.Description = "El usuario desea modificar un viaje existente"
		}
	}

	return response
}

type EntityExtractor struct {
	numberPattern *regexp.Regexp
	timePattern   *regexp.Regexp
	datePattern   *regexp.Regexp
}

func NewEntityExtractor() *EntityExtractor {
	return &EntityExtractor{
		numberPattern: regexp.MustCompile(`(\d+)\s*(pasajero|persona|gente|personas)`),
		timePattern:   regexp.MustCompile(`(\d{1,2}):(\d{2})\s*(am|pm|AM|PM)?`),
		datePattern:   regexp.MustCompile(`(mañana|hoy|pasado\s*mañana|lunes|martes|miércoles|jueves|viernes|sábado|domingo|este\s*\w+|\d{1,2}\s*de\s*[a-z]+)`),
	}
}

func (e *EntityExtractor) Extract(input string, intent llm.IntentKind) *llm.Entities {
	entities := &llm.Entities{}

	switch intent {
	case llm.IntentTripRequest:
		e.extractPlaces(input, entities)
		e.extractPassengers(input, entities)
		e.extractSchedule(input, entities)
		e.extractLuggage(input, entities)
	case llm.IntentCancel, llm.IntentChange, llm.IntentTripStatus:
		e.extractTripID(input, entities)
	}

	return entities
}

func (e *EntityExtractor) extractPlaces(input string, entities *llm.Entities) {
	lower := strings.ToLower(input)

	prepPattern := regexp.MustCompile(`(?:de|desde)\s+([A-ZÁÉÍÓÚÑa-záéíóúñ][A-ZÁÉÍÓÚÑa-záéíóúñ\s]{2,30})\s+(?:a|hacia|para|hasta)\s+([A-ZÁÉÍÓÚÑa-záéíóúñ][A-ZÁÉÍÓÚÑa-záéíóúñ\s]{2,30})`)
	prepMatches := prepPattern.FindStringSubmatch(lower)
	if len(prepMatches) > 2 {
		entities.Origin = &llm.Place{Name: strings.TrimSpace(prepMatches[1])}
		entities.Destination = &llm.Place{Name: strings.TrimSpace(prepMatches[2])}
		return
	}

	if entities.Origin == nil {
		fromPattern := regexp.MustCompile(`(?:desde|de|salir\s*de)\s+([A-ZÁÉÍÓÚÑa-záéíóúñ][A-ZÁÉÍÓÚÑa-záéíóúñ\s]{2,30})`)
		fromMatches := fromPattern.FindStringSubmatch(lower)
		if len(fromMatches) > 1 {
			entities.Origin = &llm.Place{Name: strings.TrimSpace(fromMatches[1])}
		}
	}

	if entities.Destination == nil {
		toPattern := regexp.MustCompile(`(?:a|hacia|para|hasta|llevar\s*a|destino\s*)\s+([A-ZÁÉÍÓÚÑa-záéíóúñ][A-ZÁÉÍÓÚÑa-záéíóúñ\s]{2,30})`)
		toMatches := toPattern.FindStringSubmatch(lower)
		if len(toMatches) > 1 {
			entities.Destination = &llm.Place{Name: strings.TrimSpace(toMatches[1])}
		}
	}

	if strings.Contains(lower, "aquí") || strings.Contains(lower, "donde estoy") || strings.Contains(lower, "mi ubicación") {
		if entities.Origin == nil {
			entities.Origin = &llm.Place{Name: "Ubicación actual", IsCurrent: true}
		}
	}

	landmarkPattern := regexp.MustCompile(`(aeropuerto|terminal|centro|mall|universidad|hospital|parque|plaza|estación)`)
	landmarkMatches := landmarkPattern.FindStringSubmatch(lower)
	if len(landmarkMatches) > 1 {
		name := landmarkMatches[1]
		cappedName := strings.ToUpper(name[:1]) + name[1:]
		if entities.Origin == nil {
			entities.Origin = &llm.Place{Name: cappedName}
		} else if entities.Destination == nil {
			entities.Destination = &llm.Place{Name: cappedName}
		}
	}
}

func (e *EntityExtractor) extractPassengers(input string, entities *llm.Entities) {
	matches := e.numberPattern.FindStringSubmatch(strings.ToLower(input))
	if len(matches) > 1 {
		count := 0
		for _, r := range matches[1] {
			count = count*10 + int(r-'0')
		}
		entities.Passengers = count
		return
	}
	if strings.Contains(strings.ToLower(input), "solo") ||
		strings.Contains(strings.ToLower(input), "sola") ||
		strings.Contains(strings.ToLower(input), "uno") {
		entities.Passengers = 1
	}
}

func (e *EntityExtractor) extractSchedule(input string, entities *llm.Entities) {
	if entities.Schedule == nil {
		entities.Schedule = &llm.Schedule{}
	}

	if strings.Contains(strings.ToLower(input), "ahora") ||
		strings.Contains(strings.ToLower(input), "ya") ||
		strings.Contains(strings.ToLower(input), "urgente") ||
		strings.Contains(strings.ToLower(input), "inmediato") ||
		strings.Contains(strings.ToLower(input), "lo más pronto") {
		entities.Schedule.IsNow = true
		return
	}

	timeMatches := e.timePattern.FindStringSubmatch(input)
	if len(timeMatches) > 2 {
		entities.Schedule.Time = timeMatches[1] + ":" + timeMatches[2]
		if len(timeMatches) > 3 && timeMatches[3] != "" {
			entities.Schedule.Time += " " + timeMatches[3]
		}
		entities.Schedule.IsNow = false
	}

	dateMatches := e.datePattern.FindStringSubmatch(strings.ToLower(input))
	if len(dateMatches) > 1 {
		entities.Schedule.Date = dateMatches[1]
		entities.Schedule.IsNow = false
	}
}

func (e *EntityExtractor) extractLuggage(input string, entities *llm.Entities) {
	lower := strings.ToLower(input)
	if strings.Contains(lower, "maleta") || strings.Contains(lower, "equipaje") || strings.Contains(lower, "valija") {
		entities.Luggage = "Con equipaje"
		numPattern := regexp.MustCompile(`(\d+)\s*(maleta|valija|equipaje)`)
		matches := numPattern.FindStringSubmatch(lower)
		if len(matches) > 1 {
			entities.Luggage = matches[1] + " maleta(s)"
		}
	}
}

func (e *EntityExtractor) extractTripID(input string, entities *llm.Entities) {
	idPattern := regexp.MustCompile(`(?:viaje|trip|id|#)\s*[:\s]*([A-Za-z0-9_-]{4,20})`)
	matches := idPattern.FindStringSubmatch(input)
	if len(matches) > 1 {
		entities.TripID = matches[1]
	}
}

type OutputValidator struct{}

func NewOutputValidator() *OutputValidator {
	return &OutputValidator{}
}

func (v *OutputValidator) Validate(resp *llm.LLMResponse) *llm.LLMResponse {
	if resp.Intent.Kind == llm.IntentTripRequest && resp.Confidence > 0.3 {
		if resp.Entities.Origin != nil && resp.Entities.Destination != nil &&
			resp.Entities.Origin.Name == resp.Entities.Destination.Name {
			resp.Confidence = 0.2
			resp.NeedsClarification = true
			resp.ClarificationQuestion = "El origen y destino parecen ser el mismo. ¿Puedes indicarme el destino correcto?"
		}
	}
	if resp.Confidence < 0.3 && !resp.NeedsClarification {
		resp.NeedsClarification = true
		resp.ClarificationQuestion = "No entendí bien tu mensaje. ¿Podrías repetirlo con más detalles?"
	}
	return resp
}
