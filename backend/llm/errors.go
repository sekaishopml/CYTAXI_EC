package llm

import "errors"

var (
	ErrEmptyInput       = errors.New("llm: empty input after sanitization")
	ErrPromptInjection  = errors.New("llm: prompt injection detected")
	ErrSensitiveData    = errors.New("llm: sensitive data detected in input")
	ErrLowConfidence    = errors.New("llm: confidence below threshold")
	ErrProviderFailed   = errors.New("llm: provider request failed")
	ErrNoProvider       = errors.New("llm: no provider available")
	ErrRateLimit        = errors.New("llm: rate limit exceeded")
)
