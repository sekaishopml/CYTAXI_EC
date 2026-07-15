package application

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/sekaishopml/cytaxi/backend/engines/policy/internal/policy/domain"
)

type PolicyRegistry struct {
	mu       sync.RWMutex
	policies map[domain.PolicyID]domain.Policy
	byDomain map[string][]domain.PolicyID
	loader   PolicyLoader
}

func NewPolicyRegistry(loader PolicyLoader) *PolicyRegistry {
	return &PolicyRegistry{
		policies: make(map[domain.PolicyID]domain.Policy),
		byDomain: make(map[string][]domain.PolicyID),
		loader:   loader,
	}
}

type PolicyLoader interface {
	Load(ctx context.Context) ([]domain.Policy, error)
	LoadByDomain(ctx context.Context, domain string) ([]domain.Policy, error)
}

func (r *PolicyRegistry) Register(policy domain.Policy) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.policies[policy.ID]; exists {
		return fmt.Errorf("%w: %s", domain.ErrDuplicatePolicy, policy.ID)
	}

	r.policies[policy.ID] = policy
	r.byDomain[policy.Domain] = append(r.byDomain[policy.Domain], policy.ID)
	return nil
}

func (r *PolicyRegistry) RegisterMany(policies []domain.Policy) error {
	for _, p := range policies {
		if err := r.Register(p); err != nil {
			return err
		}
	}
	return nil
}

func (r *PolicyRegistry) GetPolicy(id domain.PolicyID) (*domain.Policy, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	policy, ok := r.policies[id]
	if !ok {
		return nil, fmt.Errorf("%w: %s", domain.ErrPolicyNotFound, id)
	}
	return &policy, nil
}

func (r *PolicyRegistry) GetPolicies(domains ...string) []domain.Policy {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []domain.Policy

	if len(domains) == 0 {
		for _, p := range r.policies {
			result = append(result, p)
		}
	} else {
		seen := make(map[domain.PolicyID]bool)
		for _, d := range domains {
			for _, id := range r.byDomain[d] {
				if !seen[id] {
					result = append(result, r.policies[id])
					seen[id] = true
				}
			}
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Priority > result[j].Priority
	})

	return result
}

func (r *PolicyRegistry) Reload(ctx context.Context) error {
	policies, err := r.loader.Load(ctx)
	if err != nil {
		return fmt.Errorf("reload policies: %w", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.policies = make(map[domain.PolicyID]domain.Policy)
	r.byDomain = make(map[string][]domain.PolicyID)

	for _, p := range policies {
		r.policies[p.ID] = p
		r.byDomain[p.Domain] = append(r.byDomain[p.Domain], p.ID)
	}

	return nil
}

func (r *PolicyRegistry) RemovePolicy(id domain.PolicyID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	policy, ok := r.policies[id]
	if !ok {
		return fmt.Errorf("%w: %s", domain.ErrPolicyNotFound, id)
	}

	delete(r.policies, id)

	ids := r.byDomain[policy.Domain]
	for i, pid := range ids {
		if pid == id {
			r.byDomain[policy.Domain] = append(ids[:i], ids[i+1:]...)
			break
		}
	}

	return nil
}
