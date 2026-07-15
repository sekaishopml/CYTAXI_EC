package domain

import "errors"

var (
	ErrPolicyNotFound      = errors.New("policy: not found")
	ErrPolicyDisabled      = errors.New("policy: disabled")
	ErrRuleNotFound        = errors.New("policy: rule not found")
	ErrInvalidVersion      = errors.New("policy: invalid version")
	ErrDuplicatePolicy     = errors.New("policy: duplicate id")
	ErrCyclicDependency    = errors.New("policy: cyclic dependency detected")
	ErrInvalidCondition    = errors.New("policy: invalid condition")
	ErrEvaluationFailed    = errors.New("policy: evaluation failed")
	ErrUnsupportedOperator = errors.New("policy: unsupported operator")
)
