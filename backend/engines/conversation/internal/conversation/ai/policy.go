package ai

import (
	"github.com/sekaishopml/cytaxi/backend/engines/conversation/internal/conversation/domain/entity"
)

type PolicyRule struct {
	Intent       IntentKind
	UseLLM       bool
	MinConfidence float64
	Priority     int
}

type PolicyEvaluator struct {
	rules []PolicyRule
}

func NewPolicyEvaluator(rules []PolicyRule) *PolicyEvaluator {
	return &PolicyEvaluator{rules: rules}
}

func NewDefaultPolicy() *PolicyEvaluator {
	return NewPolicyEvaluator([]PolicyRule{
		{Intent: IntentGreeting, UseLLM: false, MinConfidence: 0.0, Priority: 1},
		{Intent: IntentTripRequest, UseLLM: true, MinConfidence: 0.7, Priority: 10},
		{Intent: IntentTripStatus, UseLLM: false, MinConfidence: 0.0, Priority: 5},
		{Intent: IntentSupport, UseLLM: true, MinConfidence: 0.6, Priority: 8},
		{Intent: IntentCancel, UseLLM: false, MinConfidence: 0.0, Priority: 7},
		{Intent: IntentUnknown, UseLLM: true, MinConfidence: 0.4, Priority: 2},
	})
}

func (e *PolicyEvaluator) ShouldUseLLM(intent *Intent, session *entity.Session) bool {
	if intent == nil {
		return false
	}

	var matched *PolicyRule
	for _, rule := range e.rules {
		if rule.Intent == intent.Kind {
			if matched == nil || rule.Priority > matched.Priority {
				matched = &rule
			}
		}
	}

	if matched == nil {
		return false
	}

	if !matched.UseLLM {
		return false
	}

	return intent.Confidence >= matched.MinConfidence
}

func (e *PolicyEvaluator) AddRule(rule PolicyRule) {
	e.rules = append(e.rules, rule)
}
