package application

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sekaishopml/cytaxi/backend/engines/policy/internal/policy/domain"
)

type fileLoader struct {
	paths []string
}

func NewFileLoader(paths ...string) PolicyLoader {
	return &fileLoader{paths: paths}
}

func (l *fileLoader) Load(ctx context.Context) ([]domain.Policy, error) {
	var allPolicies []domain.Policy

	for _, path := range l.paths {
		policies, err := l.loadFile(path)
		if err != nil {
			return nil, fmt.Errorf("load %s: %w", path, err)
		}
		allPolicies = append(allPolicies, policies...)
	}

	return allPolicies, nil
}

func (l *fileLoader) LoadByDomain(ctx context.Context, domainName string) ([]domain.Policy, error) {
	all, err := l.Load(ctx)
	if err != nil {
		return nil, err
	}

	var filtered []domain.Policy
	for _, p := range all {
		if p.Domain == domainName {
			filtered = append(filtered, p)
		}
	}
	return filtered, nil
}

func (l *fileLoader) loadFile(path string) ([]domain.Policy, error) {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}

	var policies []domain.Policy
	if err := json.Unmarshal(data, &policies); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return policies, nil
}

type memoryLoader struct {
	policies []domain.Policy
}

func NewMemoryLoader(policies []domain.Policy) PolicyLoader {
	return &memoryLoader{policies: policies}
}

func (l *memoryLoader) Load(ctx context.Context) ([]domain.Policy, error) {
	return l.policies, nil
}

func (l *memoryLoader) LoadByDomain(ctx context.Context, domainName string) ([]domain.Policy, error) {
	var filtered []domain.Policy
	for _, p := range l.policies {
		if p.Domain == domainName {
			filtered = append(filtered, p)
		}
	}
	return filtered, nil
}
