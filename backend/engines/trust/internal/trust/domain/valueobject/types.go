package valueobject

import (
	"fmt"
	"time"
)

type IdentityID string
type VerificationID string
type DocumentID string
type FraudCheckID string
type RiskAssessmentID string

type VerificationStatus string

const (
	VerifPending   VerificationStatus = "pending"
	VerifInReview  VerificationStatus = "in_review"
	VerifApproved  VerificationStatus = "approved"
	VerifRejected  VerificationStatus = "rejected"
	VerifExpired   VerificationStatus = "expired"
)

type DocumentType string

const (
	DocNationalID      DocumentType = "national_id"
	DocPassport        DocumentType = "passport"
	DocDriverLicense   DocumentType = "driver_license"
	DocVehicleReg      DocumentType = "vehicle_registration"
	DocInsurance       DocumentType = "insurance"
	DocProofOfAddress  DocumentType = "proof_of_address"
	DocSelfie          DocumentType = "selfie"
	DocBackgroundCheck DocumentType = "background_check"
)

type IdentityType string

const (
	IDTypeCustomer IdentityType = "customer"
	IDTypeDriver   IdentityType = "driver"
	IDTypeAdmin    IdentityType = "admin"
)

type TrustLevel int

const (
	TrustNone     TrustLevel = 0
	TrustBasic    TrustLevel = 1
	TrustVerified TrustLevel = 2
	TrustPremium  TrustLevel = 3
)

type RiskLevel string

const (
	RiskLow       RiskLevel = "low"
	RiskMedium    RiskLevel = "medium"
	RiskHigh      RiskLevel = "high"
	RiskCritical  RiskLevel = "critical"
)

func NewIdentityID() IdentityID       { return IdentityID(fmt.Sprintf("id_%d", time.Now().UnixNano())) }
func NewVerificationID() VerificationID { return VerificationID(fmt.Sprintf("ver_%d", time.Now().UnixNano())) }
func NewDocumentID() DocumentID       { return DocumentID(fmt.Sprintf("doc_%d", time.Now().UnixNano())) }
func NewFraudCheckID() FraudCheckID   { return FraudCheckID(fmt.Sprintf("fr_%d", time.Now().UnixNano())) }
func NewRiskAssessmentID() RiskAssessmentID { return RiskAssessmentID(fmt.Sprintf("risk_%d", time.Now().UnixNano())) }
