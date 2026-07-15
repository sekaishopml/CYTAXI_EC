# Trust & Identity Engine

Identity verification, trust scoring, and fraud detection for the CYTAXI platform.

## Purpose

The Trust & Identity Engine is the exclusive owner of identity management, verification, reputation, and trust scoring. No other Engine may validate users or assess risk directly.

## Architecture

DDD ✓ Clean Architecture ✓ CQRS ✓ Event Driven ✓ Contract First ✓ Zero Trust ✓

## Domain

### Aggregates
- **Identity** — Core identity (pending→approved→verified, 4 trust levels)
- **Verification** — Multi-type verifications (document, selfie, biometric, KYC, AML)
- **FraudAssessment** — Fraud check with flag-based risk scoring
- **TrustProfile** — Weighted trust score (verification 40%, activity 30%, community 20%, compliance 10%)

### Status Lifecycle
```
Identity: pending → approved → (trust levels: basic → verified → premium)
Verification: pending → in_review → approved/rejected/expired
```

### Entities
- **Document** — Uploaded verification documents (8 types)
- **RiskAssessment** — Risk profile assessment
- **BlacklistEntry** — Blacklisted identities with severity
- **WhitelistEntry** — Trusted identities

### Value Objects
IdentityID, VerificationID, DocumentID, FraudCheckID, RiskAssessmentID, VerificationStatus(5), DocumentType(8), IdentityType(3), TrustLevel(4), RiskLevel(4)

## Trust Score Calculation

```
Score = VerificationScore*0.4 + ActivityScore*0.3 + CommunityScore*0.2 + ComplianceScore*0.1

TrustNone      (<30)
TrustBasic     (30-59)
TrustVerified  (60-79)
TrustPremium   (80+)
```

## Verification Types
document, selfie, biometric, kyc, aml, address, phone

## Future Provider Adapters
- Registro Civil, OCR, Face Match, Liveness, KYC Provider, AML Provider

## CQRS

**Commands:** CreateIdentity, VerifyIdentity, UploadDocument, ApproveVerification, RejectVerification, CalculateTrustScore, RunFraudCheck, BlacklistIdentity, WhitelistIdentity

**Queries:** GetIdentity, GetVerification, GetTrustScore, GetRiskAssessment, GetDocuments, GetFraudHistory

## Events

| Event | Description |
|-------|-------------|
| `trust.identity_created` | Identity created |
| `trust.identity_verified` | Identity verified |
| `trust.verification_approved` | Verification approved |
| `trust.verification_rejected` | Verification rejected |
| `trust.document_uploaded` | Document uploaded |
| `trust.trust_score_updated` | Trust score changed |
| `trust.fraud_detected` | Fraud detected |
| `trust.blacklist_updated` | Blacklist modified |
| `trust.whitelist_updated` | Whitelist modified |

## API

| Method | Path | Description |
|--------|------|-------------|
| GET | /health | Health check |
| GET | /identity/{id} | Get identity |
| GET | /identity/{id}/trust | Get trust score |

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `TRUST_PORT` | 8092 | HTTP server port |

## Development

```bash
go run ./cmd/trust
```
