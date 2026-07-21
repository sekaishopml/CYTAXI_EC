package llm

// LLMResponse is the structured output contract.
// The LLM ALWAYS returns this JSON. No business logic, no pricing, no DB access.
type LLMResponse struct {
	Intent          Intent        `json:"intent"`
	Entities        Entities      `json:"entities"`
	Confidence      float64       `json:"confidence"`
	PendingQuestions []string     `json:"pending_questions,omitempty"`
	NeedsClarification bool       `json:"needs_clarification"`
	ClarificationQuestion string `json:"clarification_question,omitempty"`
	RawInput        string        `json:"raw_input"`
}

type Intent struct {
	Kind        IntentKind `json:"kind"`
	Description string     `json:"description"`
}

type IntentKind string

const (
	IntentGreeting     IntentKind = "greeting"
	IntentTripRequest  IntentKind = "trip_request"
	IntentTripStatus   IntentKind = "trip_status"
	IntentSupport      IntentKind = "support"
	IntentCancel       IntentKind = "cancel"
	IntentChange       IntentKind = "change"
	IntentUnknown      IntentKind = "unknown"
)

// Entities extracted from natural language. LLM fills what it detects.
// Backend validates everything — LLM values are suggestions, not authority.
type Entities struct {
	Origin      *Place       `json:"origin,omitempty"`
	Destination *Place       `json:"destination,omitempty"`
	Passengers  int          `json:"passengers,omitempty"`
	Luggage     string       `json:"luggage,omitempty"`
	Schedule    *Schedule    `json:"schedule,omitempty"`
	VehicleType string       `json:"vehicle_type,omitempty"`
	Preferences []string     `json:"preferences,omitempty"`
	TripID      string       `json:"trip_id,omitempty"`
	Phone       string       `json:"phone,omitempty"`
	Name        string       `json:"name,omitempty"`
}

type Place struct {
	Name       string  `json:"name"`
	Address    string  `json:"address,omitempty"`
	Lat        float64 `json:"lat,omitempty"`
	Lng        float64 `json:"lng,omitempty"`
	IsCurrent  bool    `json:"is_current,omitempty"`
}

type Schedule struct {
	IsNow  bool   `json:"is_now"`
	Date   string `json:"date,omitempty"`
	Time   string `json:"time,omitempty"`
	Flexible bool `json:"flexible,omitempty"`
}

func (r *LLMResponse) HasHighConfidence() bool {
	return r.Confidence >= 0.7 && !r.NeedsClarification
}

func (r *LLMResponse) HasMediumConfidence() bool {
	return r.Confidence >= 0.4 && r.Confidence < 0.7
}

func (r *LLMResponse) HasCompleteTripRequest() bool {
	if r.Intent.Kind != IntentTripRequest {
		return false
	}
	return r.Entities.Origin != nil && r.Entities.Destination != nil
}
