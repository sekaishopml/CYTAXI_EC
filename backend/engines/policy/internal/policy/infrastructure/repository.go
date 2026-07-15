package infrastructure

import (
	"context"

	"github.com/sekaishopml/cytaxi/backend/engines/policy/internal/policy/domain"
)

type PolicyRepository interface {
	FindByID(ctx context.Context, id domain.PolicyID) (*domain.Policy, error)
	FindByDomain(ctx context.Context, domain string) ([]domain.Policy, error)
	FindAll(ctx context.Context) ([]domain.Policy, error)
	Save(ctx context.Context, policy domain.Policy) error
	Update(ctx context.Context, policy domain.Policy) error
	Delete(ctx context.Context, id domain.PolicyID) error
}

type VersionRepository interface {
	FindByPolicyID(ctx context.Context, policyID domain.PolicyID) ([]domain.VersionRecord, error)
	Save(ctx context.Context, record domain.VersionRecord) error
	GetLatest(ctx context.Context, policyID domain.PolicyID) (*domain.VersionRecord, error)
}
