package port

import (
	"context"

	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/application/command"
	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/application/query"
	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/domain/blacklist"
	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/domain/document"
	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/domain/fraud"
	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/domain/identity"
	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/domain/trustscore"
	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/domain/verification"
)

type TrustService interface {
	CreateIdentity(ctx context.Context, cmd command.CreateIdentity) (*identity.Identity, error)
	VerifyIdentity(ctx context.Context, cmd command.VerifyIdentity) (*verification.Verification, error)
	UploadDocument(ctx context.Context, cmd command.UploadDocument) (*document.Document, error)
	ApproveVerification(ctx context.Context, cmd command.ApproveVerification) error
	RejectVerification(ctx context.Context, cmd command.RejectVerification) error
	CalculateTrustScore(ctx context.Context, cmd command.CalculateTrustScore) (*trustscore.TrustProfile, error)
	RunFraudCheck(ctx context.Context, cmd command.RunFraudCheck) (*fraud.FraudAssessment, error)
	Blacklist(ctx context.Context, cmd command.BlacklistIdentity) (*blacklist.BlacklistEntry, error)
	Whitelist(ctx context.Context, cmd command.WhitelistIdentity) (*whitelist.WhitelistEntry, error)
	GetIdentity(ctx context.Context, q query.GetIdentity) (*identity.Identity, error)
	GetTrustScore(ctx context.Context, q query.GetTrustScore) (*trustscore.TrustProfile, error)
	GetDocuments(ctx context.Context, q query.GetDocuments) ([]document.Document, error)
	GetFraudHistory(ctx context.Context, q query.GetFraudHistory) ([]fraud.FraudAssessment, error)
}
