package externalapi

type Client struct{}

func New() *Client {
	return &Client{}
}
