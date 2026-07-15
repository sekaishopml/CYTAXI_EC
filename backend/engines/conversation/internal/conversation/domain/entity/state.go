package entity

import "fmt"

type ConversationState string

const (
	StateNew          ConversationState = "new"
	StateActive       ConversationState = "active"
	StateWaitingInput ConversationState = "waiting_input"
	StateProcessing   ConversationState = "processing"
	StateIdle         ConversationState = "idle"
	StateExpired      ConversationState = "expired"
	StateClosed       ConversationState = "closed"
)

type StateTransition struct {
	From ConversationState
	To   ConversationState
}

var validTransitions = map[ConversationState][]ConversationState{
	StateNew:          {StateActive},
	StateActive:       {StateWaitingInput, StateIdle, StateClosed},
	StateWaitingInput: {StateProcessing, StateIdle, StateClosed},
	StateProcessing:   {StateWaitingInput, StateActive, StateClosed},
	StateIdle:         {StateActive, StateExpired, StateClosed},
	StateExpired:      {StateClosed},
	StateClosed:       {},
}

func ValidateTransition(from, to ConversationState) error {
	allowed, ok := validTransitions[from]
	if !ok {
		return fmt.Errorf("invalid source state: %s", from)
	}
	for _, s := range allowed {
		if s == to {
			return nil
		}
	}
	return fmt.Errorf("transition from %s to %s is not allowed", from, to)
}

func CanTransition(from, to ConversationState) bool {
	return ValidateTransition(from, to) == nil
}
