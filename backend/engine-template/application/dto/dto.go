package dto

type Request struct {
	ID string `json:"id"`
}

type Response struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}
