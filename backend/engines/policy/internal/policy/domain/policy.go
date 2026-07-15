package domain

import "time"

type PolicyID string
type RuleID string

type Policy struct {
	ID          PolicyID     `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Domain      string       `json:"domain"` // pricing, dispatch, zones, etc.
	Version     Version      `json:"version"`
	Rules       []Rule       `json:"rules"`
	Priority    int          `json:"priority"`
	Enabled     bool         `json:"enabled"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type Rule struct {
	ID          RuleID          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Conditions  []Condition     `json:"conditions"`
	Actions     []Action        `json:"actions"`
	Priority    int             `json:"priority"`
	Enabled     bool            `json:"enabled"`
	Metadata    map[string]any  `json:"metadata,omitempty"`
}

type Condition struct {
	Field    string      `json:"field"`
	Operator Operator    `json:"operator"`
	Value    any         `json:"value"`
}

type Operator string

const (
	OpEquals       Operator = "equals"
	OpNotEquals    Operator = "not_equals"
	OpGreaterThan  Operator = "greater_than"
	OpLessThan     Operator = "less_than"
	OpIn           Operator = "in"
	OpNotIn        Operator = "not_in"
	OpBetween      Operator = "between"
	OpContains     Operator = "contains"
	OpStartsWith   Operator = "starts_with"
	OpMatches      Operator = "matches"
	OpTrue         Operator = "is_true"
	OpFalse        Operator = "is_false"
)

type Action struct {
	Type   string         `json:"type"`
	Params map[string]any `json:"params"`
}

type PolicyStatus string

const (
	PolicyDraft     PolicyStatus = "draft"
	PolicyActive    PolicyStatus = "active"
	PolicyDeprecated PolicyStatus = "deprecated"
	PolicyArchived  PolicyStatus = "archived"
)
