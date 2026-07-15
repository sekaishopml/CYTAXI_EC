package kyc

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type VerificationStatus string

const (
	StatusPending   VerificationStatus = "PENDING"
	StatusInReview  VerificationStatus = "IN_REVIEW"
	StatusApproved  VerificationStatus = "APPROVED"
	StatusRejected  VerificationStatus = "REJECTED"
	StatusSuspended VerificationStatus = "SUSPENDED"
	StatusExpired   VerificationStatus = "EXPIRED"
)

type DocumentType string

const (
	DocDriverLicense     DocumentType = "driver_license"
	DocVehicleRegistration DocumentType = "vehicle_registration"
	DocNationalID        DocumentType = "national_id"
	DocBackgroundCheck   DocumentType = "background_check"
	DocInsurance         DocumentType = "insurance"
	DocVehiclePhotos     DocumentType = "vehicle_photos"
	DocProfilePhoto      DocumentType = "profile_photo"
)

type Document struct {
	ID         string       `json:"id"`
	DriverID   string       `json:"driver_id"`
	Type       DocumentType `json:"type"`
	URL        string       `json:"url"`
	Filename   string       `json:"filename"`
	UploadedAt time.Time    `json:"uploaded_at"`
	Verified   bool         `json:"verified"`
	ReviewedBy string       `json:"reviewed_by,omitempty"`
	Notes      string       `json:"notes,omitempty"`
}

type Verification struct {
	ID          string             `json:"id"`
	DriverID    string             `json:"driver_id"`
	DriverName  string             `json:"driver_name"`
	Phone       string             `json:"phone"`
	Status      VerificationStatus `json:"status"`
	Documents   []Document         `json:"documents"`
	Notes       string             `json:"notes,omitempty"`
	ReviewedBy  string             `json:"reviewed_by,omitempty"`
	StartedAt   time.Time          `json:"started_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	CompletedAt *time.Time         `json:"completed_at,omitempty"`
}

type KYCProvider interface {
	Name() string
	VerifyIdentity(ctx context.Context, doc Document) (*KYCResult, error)
	IsAvailable(ctx context.Context) bool
}

type KYCResult struct {
	Passed  bool    `json:"passed"`
	Score   float64 `json:"score"`
	Ref     string  `json:"ref"`
	Details map[string]string `json:"details,omitempty"`
}

type Registry struct {
	providers map[string]KYCProvider
}

func NewRegistry() *Registry { return &Registry{providers: make(map[string]KYCProvider)} }
func (r *Registry) Register(p KYCProvider) { r.providers[p.Name()] = p }
func (r *Registry) Get(name string) (KYCProvider, error) {
	p, ok := r.providers[name]
	if !ok { return nil, fmt.Errorf("KYC provider %s not found", name) }
	return p, nil
}

type MockKYC struct{ name string }

func NewMockKYC() KYCProvider { return &MockKYC{name: "mock_kyc"} }
func (m *MockKYC) Name() string { return m.name }
func (m *MockKYC) IsAvailable(ctx context.Context) bool { return true }
func (m *MockKYC) VerifyIdentity(ctx context.Context, doc Document) (*KYCResult, error) {
	return &KYCResult{Passed: true, Score: 0.85, Ref: fmt.Sprintf("kyc_%d", time.Now().UnixNano())}, nil
}

type OCRProvider struct{ name string }
func NewOCRProvider() KYCProvider { return &OCRProvider{name: "ocr"} }
func (o *OCRProvider) Name() string { return o.name }
func (o *OCRProvider) IsAvailable(ctx context.Context) bool { return false }
func (o *OCRProvider) VerifyIdentity(ctx context.Context, doc Document) (*KYCResult, error) {
	return nil, fmt.Errorf("OCR provider not configured")
}

type Manager struct {
	verifications sync.Map
	docs          sync.Map
}

func NewManager() *Manager { return &Manager{} }

func (m *Manager) StartVerification(driverID, name, phone string) *Verification {
	v := &Verification{
		ID: fmt.Sprintf("ver_%d", time.Now().UnixNano()),
		DriverID: driverID, DriverName: name, Phone: phone,
		Status: StatusPending, StartedAt: time.Now(), UpdatedAt: time.Now(),
	}
	m.verifications.Store(v.ID, v)
	return v
}

func (m *Manager) GetVerification(id string) (*Verification, error) {
	v, ok := m.verifications.Load(id)
	if !ok { return nil, fmt.Errorf("verification %s not found", id) }
	return v.(*Verification), nil
}

func (m *Manager) GetByDriver(driverID string) (*Verification, error) {
	var found *Verification
	m.verifications.Range(func(_, v any) bool {
		if v.(*Verification).DriverID == driverID {
			found = v.(*Verification)
			return false
		}
		return true
	})
	if found == nil { return nil, fmt.Errorf("no verification for driver %s", driverID) }
	return found, nil
}

func (m *Manager) AddDocument(verID string, doc Document) error {
	v, err := m.GetVerification(verID)
	if err != nil { return err }
	v.Documents = append(v.Documents, doc)
	v.Status = StatusInReview
	v.UpdatedAt = time.Now()
	m.verifications.Store(verID, v)
	m.docs.Store(doc.ID, doc)
	return nil
}

func (m *Manager) Approve(verID, reviewer, notes string) error {
	v, err := m.GetVerification(verID)
	if err != nil { return err }
	now := time.Now()
	v.Status = StatusApproved
	v.ReviewedBy = reviewer
	v.Notes = notes
	v.CompletedAt = &now
	v.UpdatedAt = now
	m.verifications.Store(verID, v)
	return nil
}

func (m *Manager) Reject(verID, reviewer, reason string) error {
	v, err := m.GetVerification(verID)
	if err != nil { return err }
	now := time.Now()
	v.Status = StatusRejected
	v.ReviewedBy = reviewer
	v.Notes = reason
	v.CompletedAt = &now
	v.UpdatedAt = now
	m.verifications.Store(verID, v)
	return nil
}

func (m *Manager) Suspend(verID, reviewer, reason string) error {
	v, err := m.GetVerification(verID)
	if err != nil { return err }
	v.Status = StatusSuspended
	v.ReviewedBy = reviewer
	v.Notes = reason
	v.UpdatedAt = time.Now()
	m.verifications.Store(verID, v)
	return nil
}

func (m *Manager) ListByStatus(status VerificationStatus) []*Verification {
	var result []*Verification
	m.verifications.Range(func(_, v any) bool {
		if v.(*Verification).Status == status {
			result = append(result, v.(*Verification))
		}
		return true
	})
	return result
}
