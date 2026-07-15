package entity

import "time"

type ContextEntry struct {
	Key       string
	Value     string
	UpdatedAt time.Time
}

type ConversationContext struct {
	ConversationID ConversationID
	Entries        map[string]ContextEntry
	UpdatedAt      time.Time
}

func NewConversationContext(convID ConversationID) *ConversationContext {
	return &ConversationContext{
		ConversationID: convID,
		Entries:        make(map[string]ContextEntry),
		UpdatedAt:      time.Now(),
	}
}

func (c *ConversationContext) Set(key, value string) {
	c.Entries[key] = ContextEntry{
		Key:       key,
		Value:     value,
		UpdatedAt: time.Now(),
	}
	c.UpdatedAt = time.Now()
}

func (c *ConversationContext) Get(key string) (string, bool) {
	entry, ok := c.Entries[key]
	if !ok {
		return "", false
	}
	return entry.Value, true
}

func (c *ConversationContext) Delete(key string) {
	delete(c.Entries, key)
	c.UpdatedAt = time.Now()
}

func (c *ConversationContext) Has(key string) bool {
	_, ok := c.Entries[key]
	return ok
}

func (c *ConversationContext) Clear() {
	c.Entries = make(map[string]ContextEntry)
	c.UpdatedAt = time.Now()
}
