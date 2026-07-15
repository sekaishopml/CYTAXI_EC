package command

import "github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/domain/valueobject"

type CreateIdentity struct {
	OwnerID string
	Type    valueobject.IdentityType
	Phone   string
	Email   string
}

type VerifyIdentity struct {
	IdentityID valueobject.IdentityID
	Type       string
}

type UploadDocument struct {
	IdentityID valueobject.IdentityID
	DocType    valueobject.DocumentType
	URL        string
}

type ApproveVerification struct {
	VerificationID valueobject.VerificationID
	Score          float64
}

type RejectVerification struct {
	VerificationID valueobject.VerificationID
	Reason         string
}

type CalculateTrustScore struct {
	IdentityID valueobject.IdentityID
	VerifyScore float64
	Activity    float64
	Community   float64
	Compliance  float64
}

type RunFraudCheck struct {
	IdentityID valueobject.IdentityID
}

type BlacklistIdentity struct {
	IdentityID valueobject.IdentityID
	Reason     string
	Severity   string
}

type WhitelistIdentity struct {
	IdentityID valueobject.IdentityID
	Reason     string
}
