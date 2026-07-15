package events

const (
	EventIdentityCreated      = "trust.identity_created"
	EventIdentityVerified     = "trust.identity_verified"
	EventVerificationApproved = "trust.verification_approved"
	EventVerificationRejected = "trust.verification_rejected"
	EventDocumentUploaded     = "trust.document_uploaded"
	EventTrustScoreUpdated    = "trust.trust_score_updated"
	EventFraudDetected        = "trust.fraud_detected"
	EventBlacklistUpdated     = "trust.blacklist_updated"
	EventWhitelistUpdated     = "trust.whitelist_updated"
)

type IdentityCreatedPayload struct {
	IdentityID string `json:"identity_id"`
	OwnerID    string `json:"owner_id"`
	Type       string `json:"type"`
	Phone      string `json:"phone"`
}

type IdentityVerifiedPayload struct {
	IdentityID string  `json:"identity_id"`
	Score      float64 `json:"score"`
}

type VerificationApprovedPayload struct {
	VerificationID string  `json:"verification_id"`
	IdentityID     string  `json:"identity_id"`
	Type           string  `json:"type"`
	Score          float64 `json:"score"`
}

type VerificationRejectedPayload struct {
	VerificationID string `json:"verification_id"`
	IdentityID     string `json:"identity_id"`
	Reason         string `json:"reason"`
}

type DocumentUploadedPayload struct {
	DocumentID string `json:"document_id"`
	IdentityID string `json:"identity_id"`
	Type       string `json:"type"`
}

type TrustScoreUpdatedPayload struct {
	IdentityID string  `json:"identity_id"`
	OldScore   float64 `json:"old_score"`
	NewScore   float64 `json:"new_score"`
	Level      int     `json:"level"`
}

type FraudDetectedPayload struct {
	IdentityID string  `json:"identity_id"`
	RiskLevel  string  `json:"risk_level"`
	Score      float64 `json:"score"`
	Flags      int     `json:"flags_count"`
}

type BlacklistUpdatedPayload struct {
	IdentityID string `json:"identity_id"`
	Reason     string `json:"reason"`
	Severity   string `json:"severity"`
}
