package dto

type IncomingMessageRequest struct {
	Phone   string `json:"phone"`
	Content string `json:"content"`
}

type IncomingMessageResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

type ConversationResponse struct {
	ID        string `json:"id"`
	Phone     string `json:"phone"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}
