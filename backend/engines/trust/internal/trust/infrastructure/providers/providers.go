package providers

import "context"

type VerificationProvider interface {
	Name() string
	Verify(ctx context.Context, providerRef string) (*VerificationResult, error)
}

type VerificationResult struct {
	Passed   bool
	Score    float64
	Ref      string
	Details  map[string]string
}

type KYCProvider interface {
	RunKYC(ctx context.Context, identityID string, docData []byte) (*KYCResult, error)
}

type KYCResult struct {
	Passed    bool
	RiskScore float64
	Flags     []string
}
