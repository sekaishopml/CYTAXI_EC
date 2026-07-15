package query

import "github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/domain/valueobject"

type GetIdentity struct {
	IdentityID valueobject.IdentityID
}

type GetVerification struct {
	VerificationID valueobject.VerificationID
}

type GetTrustScore struct {
	IdentityID valueobject.IdentityID
}

type GetRiskAssessment struct {
	IdentityID valueobject.IdentityID
}

type GetDocuments struct {
	IdentityID valueobject.IdentityID
}

type GetFraudHistory struct {
	IdentityID valueobject.IdentityID
}
