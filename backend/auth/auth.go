package auth

type Principal struct {
	ID    string
	Role  string
	Scope []string
}

type Authenticator interface {
	Authenticate(token string) (*Principal, error)
}

type Authorizer interface {
	Authorize(principal *Principal, resource, action string) bool
}
