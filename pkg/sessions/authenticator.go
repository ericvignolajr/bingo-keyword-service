package sessions

type AuthenticatorResponse struct {
	IsAuthenticated bool   `json:"isAuthenticated"`
	ID              string `json:"id"`
	Email           string `json:"email"`
	Err             error  `json:"error,omitempty"`
}

type Authenticator interface {
	Authenticate(string) AuthenticatorResponse
}
