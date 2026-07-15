package application

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/sekaishopml/cytaxi/backend/engines/policy/internal/policy/domain"
)

type ConditionEvaluator struct{}

func NewConditionEvaluator() *ConditionEvaluator {
	return &ConditionEvaluator{}
}

func (e *ConditionEvaluator) Evaluate(ctx context.Context, condition domain.Condition, decisionCtx domain.DecisionContext) (bool, error) {
	actual := decisionCtx.Get(condition.Field)
	if actual == nil {
		return false, nil
	}

	switch condition.Operator {
	case domain.OpEquals:
		return reflect.DeepEqual(actual, condition.Value), nil

	case domain.OpNotEquals:
		return !reflect.DeepEqual(actual, condition.Value), nil

	case domain.OpGreaterThan:
		return compare(actual, condition.Value, func(a, b float64) bool { return a > b })

	case domain.OpLessThan:
		return compare(actual, condition.Value, func(a, b float64) bool { return a < b })

	case domain.OpIn:
		return inSlice(actual, condition.Value)

	case domain.OpNotIn:
		ok, err := inSlice(actual, condition.Value)
		return !ok, err

	case domain.OpContains:
		str, ok := actual.(string)
		if !ok {
			return false, nil
		}
		sub, ok := condition.Value.(string)
		if !ok {
			return false, nil
		}
		return strings.Contains(str, sub), nil

	case domain.OpStartsWith:
		str, ok := actual.(string)
		if !ok {
			return false, nil
		}
		sub, ok := condition.Value.(string)
		if !ok {
			return false, nil
		}
		return strings.HasPrefix(str, sub), nil

	case domain.OpTrue:
		b, ok := actual.(bool)
		return ok && b, nil

	case domain.OpFalse:
		b, ok := actual.(bool)
		return ok && !b, nil

	default:
		return false, fmt.Errorf("%w: %s", domain.ErrUnsupportedOperator, condition.Operator)
	}
}

func compare(a, b any, fn func(float64, float64) bool) (bool, error) {
	fa, ok := toFloat64(a)
	if !ok {
		return false, nil
	}
	fb, ok := toFloat64(b)
	if !ok {
		return false, nil
	}
	return fn(fa, fb), nil
}

func toFloat64(v any) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	default:
		return 0, false
	}
}

func inSlice(actual any, value any) (bool, error) {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return false, nil
	}
	for i := 0; i < v.Len(); i++ {
		if reflect.DeepEqual(actual, v.Index(i).Interface()) {
			return true, nil
		}
	}
	return false, nil
}
