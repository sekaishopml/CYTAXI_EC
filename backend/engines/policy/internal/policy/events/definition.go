package events

const (
	EventPolicyCreated    = "policy.created"
	EventPolicyUpdated    = "policy.updated"
	EventPolicyDeleted    = "policy.deleted"
	EventPolicyEnabled    = "policy.enabled"
	EventPolicyDisabled   = "policy.disabled"
	EventPolicyEvaluated  = "policy.evaluated"
	EventPolicyError      = "policy.error"
	EventRuleMatched      = "policy.rule_matched"
	EventRuleNotMatched   = "policy.rule_not_matched"
	EventPoliciesReloaded = "policy.reloaded"
)

type PolicyCreatedPayload struct {
	PolicyID string `json:"policy_id"`
	Name     string `json:"name"`
	Domain   string `json:"domain"`
	Version  string `json:"version"`
}

type PolicyEvaluatedPayload struct {
	PolicyID string `json:"policy_id"`
	Matched  bool   `json:"matched"`
	Duration int64  `json:"duration_ms"`
	Domain   string `json:"domain"`
}

type RuleMatchedPayload struct {
	PolicyID string `json:"policy_id"`
	RuleID   string `json:"rule_id"`
	Actions  int    `json:"actions_count"`
}
