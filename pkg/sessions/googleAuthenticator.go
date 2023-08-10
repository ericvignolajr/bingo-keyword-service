package sessions

import (
	"fmt"
	"os"

	"github.com/lestrrat-go/jwx/jwt"
)

const gissname1 = "accounts.google.com"
const gissname2 = "https://accounts.google.com"
const gcredname = "credential"

type GoogleAuthenticator struct{}

func (g *GoogleAuthenticator) Authenticate(s string) AuthenticatorResponse {
	//TODO: retrieve jwk from google and verify jwt when parsing
	t, err := jwt.Parse([]byte(s))
	if err != nil {
		return AuthenticatorResponse{
			IsAuthenticated: false,
			ID:              "",
			Email:           "",
			Err:             err,
		}
	}
	_, err = g.expectedAudience(t)
	if err != nil {
		return AuthenticatorResponse{
			IsAuthenticated: false,
			ID:              "",
			Email:           "",
			Err:             err,
		}
	}
	_, err = g.expectedIssuer(t)
	if err != nil {
		return AuthenticatorResponse{
			IsAuthenticated: false,
			ID:              "",
			Email:           "",
			Err:             err,
		}
	}

	// any final validations like checking if token is expired
	err = jwt.Validate(t)
	if err != nil {
		return AuthenticatorResponse{
			IsAuthenticated: false,
			ID:              "",
			Email:           "",
			Err:             err,
		}
	}

	email := ""
	claim, ok := t.Get("email")
	if ok {
		email = claim.(string)
	}
	return AuthenticatorResponse{
		IsAuthenticated: true,
		ID:              t.Subject(),
		Email:           email,
	}
}

func (g *GoogleAuthenticator) expectedAudience(t jwt.Token) (bool, error) {
	fmt.Println("expectedAudience...")
	err := jwt.Validate(t, jwt.WithAudience(os.Getenv("G_APPLICATION_ID")))
	if err != nil {
		return false, fmt.Errorf("in audience check: %w", err)
	}

	return true, nil
}

func (g *GoogleAuthenticator) expectedIssuer(t jwt.Token) (bool, error) {
	fmt.Println("expectedIssuer...")
	issErr1 := jwt.Validate(t, jwt.WithIssuer(gissname1))
	issErr2 := jwt.Validate(t, jwt.WithIssuer(gissname2))
	if issErr1 != nil && issErr2 != nil {
		return false, fmt.Errorf("in issuer check: %w | %w", issErr1, issErr2)
	}

	return true, nil
}
