package repository

import (
	"context"

	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/domain/blacklist"
	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/domain/document"
	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/domain/fraud"
	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/domain/identity"
	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/domain/trustscore"
	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/domain/valueobject"
	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/domain/verification"
)

type IdentityRepository interface {
	FindByID(ctx context.Context, id valueobject.IdentityID) (*identity.Identity, error)
	FindByOwnerID(ctx context.Context, ownerID string) (*identity.Identity, error)
	Save(ctx context.Context, i *identity.Identity) error
	Update(ctx context.Context, i *identity.Identity) error
}

type VerificationRepository interface {
	FindByID(ctx context.Context, id valueobject.VerificationID) (*verification.Verification, error)
	FindByIdentityID(ctx context.Context, identityID valueobject.IdentityID) ([]verification.Verification, error)
	Save(ctx context.Context, v *verification.Verification) error
	Update(ctx context.Context, v *verification.Verification) error
}

type DocumentRepository interface {
	FindByID(ctx context.Context, id valueobject.DocumentID) (*document.Document, error)
	FindByIdentityID(ctx context.Context, identityID valueobject.IdentityID) ([]document.Document, error)
	Save(ctx context.Context, d *document.Document) error
}

type FraudRepository interface {
	FindByID(ctx context.Context, id valueobject.FraudCheckID) (*fraud.FraudAssessment, error)
	FindByIdentityID(ctx context.Context, identityID valueobject.IdentityID) ([]fraud.FraudAssessment, error)
	Save(ctx context.Context, f *fraud.FraudAssessment) error
}

type TrustRepository interface {
	FindByID(ctx context.Context, identityID valueobject.IdentityID) (*trustscore.TrustProfile, error)
	Save(ctx context.Context, tp *trustscore.TrustProfile) error
}

type BlacklistRepository interface {
	FindByIdentityID(ctx context.Context, identityID valueobject.IdentityID) (*blacklist.BlacklistEntry, error)
	Save(ctx context.Context, e *blacklist.BlacklistEntry) error
	Delete(ctx context.Context, identityID valueobject.IdentityID) error
}
