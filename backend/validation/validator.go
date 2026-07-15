package validation

import (
	"fmt"
	"reflect"
	"strings"
)

type RuleFunc func(value any) error

var rules = map[string]RuleFunc{
	"required": requiredRule,
	"email":    emailRule,
}

func requiredRule(value any) error {
	if value == nil {
		return fmt.Errorf("value is required")
	}
	s, ok := value.(string)
	if ok && strings.TrimSpace(s) == "" {
		return fmt.Errorf("value is required")
	}
	return nil
}

func emailRule(value any) error {
	s, ok := value.(string)
	if !ok {
		return nil
	}
	if !strings.Contains(s, "@") || !strings.Contains(s, ".") {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

type structValidator struct{}

func New() Validator {
	return &structValidator{}
}

func (v *structValidator) Validate(input any) error {
	val := reflect.ValueOf(input)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("validation: input must be a struct")
	}

	errs := &Errors{}
	t := val.Type()

	for i := range t.NumField() {
		field := t.Field(i)
		tag := field.Tag.Get("validate")
		if tag == "" {
			continue
		}

		rulesList := strings.Split(tag, ",")
		fieldVal := val.Field(i).Interface()

		for _, ruleName := range rulesList {
			ruleName = strings.TrimSpace(ruleName)
			if fn, ok := rules[ruleName]; ok {
				if err := fn(fieldVal); err != nil {
					errs.Add(field.Name, err.Error())
				}
			}
		}
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}

func ValidateStruct(input any) error {
	v := New()
	return v.Validate(input)
}
