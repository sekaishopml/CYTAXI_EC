package ai

import "strings"

type RuleClassifier struct {
	rules []ClassificationRule
}

type ClassificationRule struct {
	Keywords []string
	Intent   IntentKind
	Priority int
}

func NewRuleClassifier(rules []ClassificationRule) *RuleClassifier {
	return &RuleClassifier{rules: rules}
}

func (c *RuleClassifier) Classify(ctx Context, input string) (*Intent, error) {
	lower := strings.ToLower(input)

	var best *ClassificationRule
	for _, rule := range c.rules {
		for _, kw := range rule.Keywords {
			if strings.Contains(lower, kw) {
				if best == nil || rule.Priority > best.Priority {
					best = &rule
					break
				}
			}
		}
	}

	if best == nil {
		return nil, nil
	}

	return &Intent{
		Kind:       best.Intent,
		Confidence: 0.8,
		Entities:   extractEntities(lower, best.Intent),
	}, nil
}

func NewDefaultClassifier() *RuleClassifier {
	return NewRuleClassifier([]ClassificationRule{
		{
			Keywords: []string{"hola", "buenos d\u00edas", "buenas", "hey", "hello", "hi"},
			Intent:   IntentGreeting,
			Priority: 1,
		},
		{
			Keywords: []string{"taxi", "viaje", "recoger", "llevar", "trip", "ride"},
			Intent:   IntentTripRequest,
			Priority: 5,
		},
		{
			Keywords: []string{"dónde", "estado", "llegada", "esperando", "status", "where"},
			Intent:   IntentTripStatus,
			Priority: 3,
		},
		{
			Keywords: []string{"cancelar", "cancel", "cancelación"},
			Intent:   IntentCancel,
			Priority: 4,
		},
		{
			Keywords: []string{"ayuda", "soporte", "problema", "help", "support", "error"},
			Intent:   IntentSupport,
			Priority: 2,
		},
	})
}

func extractEntities(input string, kind IntentKind) map[string]string {
	entities := make(map[string]string)
	entities["raw_input"] = input
	return entities
}
