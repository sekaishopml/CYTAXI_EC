package llm

import (
	"strings"
	"unicode"
)

type SecurityFilter struct {
	maxInputLength   int
	maxTokens        int
	blockedPatterns  []string
	sensitiveDomains []string
}

func NewSecurityFilter() *SecurityFilter {
	return &SecurityFilter{
		maxInputLength:  2000,
		maxTokens:       2048,
		blockedPatterns: []string{
			"ignore previous instructions",
			"ignore all instructions",
			"forget everything",
			"you are now",
			"act as a",
			" pretend you are",
			"system prompt",
			"<script",
			"javascript:",
			"onerror=",
			"onclick=",
			"onload=",
		},
		sensitiveDomains: []string{
			"database", "password", "secret", "token", "api_key",
			"private_key", "credential", "jwt_secret", "auth_secret",
		},
	}
}

func (f *SecurityFilter) SanitizeInput(input string) string {
	var result strings.Builder
	for _, r := range input {
		if r == '\n' || r == '\r' || r == '\t' {
			result.WriteRune(' ')
		} else if unicode.IsPrint(r) {
			result.WriteRune(r)
		}
	}
	return strings.TrimSpace(result.String())
}

func (f *SecurityFilter) IsPromptInjection(input string) bool {
	lower := strings.ToLower(input)
	for _, pattern := range f.blockedPatterns {
		if strings.Contains(lower, pattern) {
			return true
		}
	}
	return false
}

func (f *SecurityFilter) ContainsSensitiveData(input string) bool {
	lower := strings.ToLower(input)
	for _, domain := range f.sensitiveDomains {
		if strings.Contains(lower, domain) {
			return true
		}
	}
	return false
}

func (f *SecurityFilter) Truncate(input string) string {
	runes := []rune(input)
	if len(runes) > f.maxInputLength {
		return string(runes[:f.maxInputLength])
	}
	return input
}

func (f *SecurityFilter) ValidateAndSanitize(input string) (string, error) {
	sanitized := f.SanitizeInput(input)

	if len(sanitized) == 0 {
		return "", ErrEmptyInput
	}

	if f.IsPromptInjection(sanitized) {
		return "", ErrPromptInjection
	}

	if f.ContainsSensitiveData(sanitized) {
		return "", ErrSensitiveData
	}

	return f.Truncate(sanitized), nil
}
