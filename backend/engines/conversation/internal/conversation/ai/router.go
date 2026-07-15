package ai

import (
	"github.com/sekaishopml/cytaxi/backend/engines/conversation/internal/conversation/domain/entity"
)

type IntentKind string

const (
	IntentGreeting    IntentKind = "greeting"
	IntentTripRequest IntentKind = "trip_request"
	IntentTripStatus  IntentKind = "trip_status"
	IntentSupport     IntentKind = "support"
	IntentCancel      IntentKind = "cancel"
	IntentUnknown     IntentKind = "unknown"
)

type Intent struct {
	Kind                 IntentKind
	Confidence           float64
	RequiredCapabilities []string
	Entities             map[string]string
}

type Classifier interface {
	Classify(ctx Context, input string) (*Intent, error)
}

type Router struct {
	classifiers []Classifier
}

func NewRouter(classifiers ...Classifier) *Router {
	return &Router{classifiers: classifiers}
}

func (r *Router) AddClassifier(c Classifier) {
	r.classifiers = append(r.classifiers, c)
}

func (r *Router) Route(ctx Context, input string, session *entity.Session) *Intent {
	for _, c := range r.classifiers {
		intent, err := c.Classify(ctx, input)
		if err == nil && intent != nil && intent.Confidence > 0.5 {
			return intent
		}
	}

	if session.CurrentState == entity.StateWaitingInput {
		return &Intent{
			Kind:       IntentUnknown,
			Confidence: 0.0,
			Entities:   make(map[string]string),
		}
	}

	return &Intent{
		Kind:       IntentGreeting,
		Confidence: 0.3,
		Entities:   make(map[string]string),
	}
}
