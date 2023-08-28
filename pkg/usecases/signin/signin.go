package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/sessions"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores/inmemory"
	"github.com/google/uuid"
)

const (
	externalID = iota + 1
	userID
	email
)

func SigninHandler(w http.ResponseWriter, r *http.Request) {
	wd, err := os.Getwd()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	path := filepath.Join(wd, "/pkg", "/usecases", "signin", "/signin.html")

	http.ServeFile(w, r, path)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func MyInfoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fmt.Println(ctx)
	uID := ctx.Value(userID).(uuid.UUID)
	email := ctx.Value(email).(string)

	out := struct {
		UserID uuid.UUID
		Email  string
	}{
		UserID: uID,
		Email:  email,
	}

	json.NewEncoder(w).Encode(out)
}

type SigninRequest struct {
	Email string
}
type SigninResponse struct {
	Ok  bool
	Err error
}

type Signin struct {
	UserStore inmemory.UserStore
}

func GetClaimsFromClerkToken(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := clerk.SessionFromContext(ctx)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		ctx = context.WithValue(ctx, externalID, claims.ID)

		clerkAuth := sessions.NewClerkAuthenticator()
		u, _ := clerkAuth.Client.Users().Read(claims.Subject)
		userEmail := u.EmailAddresses[0].EmailAddress
		ctx = context.WithValue(ctx, email, userEmail)

		r = r.WithContext(ctx)
		fmt.Println("GetClaimsFromClerkToken - ok")
		h.ServeHTTP(w, r)
	})
}

type ReadUserByEmailRequest struct {
	ExternalID string
	Email      string
}

type ReadUserByEmailResponse struct {
	ID    uuid.UUID
	Email string
}

type ReadUserByEmail struct {
	UserStore stores.User
}

func (r *ReadUserByEmail) Exec(req ReadUserByEmailRequest) (ReadUserByEmailResponse, error) {
	u, err := r.UserStore.ReadByEmail(req.Email)
	if err != nil {
		return ReadUserByEmailResponse{
			ID:    uuid.Nil,
			Email: "",
		}, err
	}

	if u == nil {
		return ReadUserByEmailResponse{
			ID:    uuid.Nil,
			Email: "",
		}, nil
	}

	return ReadUserByEmailResponse{
		ID:    u.Id,
		Email: u.Email,
	}, nil
}

func (ru *ReadUserByEmail) AsMiddleWare(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		eID, ok := ctx.Value(externalID).(string)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}
		emailAddress, ok := ctx.Value(email).(string)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}
		req := ReadUserByEmailRequest{
			ExternalID: eID,
			Email:      emailAddress,
		}
		res, err := ru.Exec(req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx = context.WithValue(ctx, userID, res.ID)
		ctx = context.WithValue(ctx, emailAddress, res.Email)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}

type CreateUserAccount struct {
	UserStore stores.User
}

func (c *CreateUserAccount) Exec(email string) error {
	u, err := c.UserStore.ReadByEmail(email)
	if err != nil {
		return err
	}
	if u != nil {
		return nil
	}

	err = c.UserStore.CreateAccount(email)
	if err != nil {
		return err
	}

	return nil
}

func (c *CreateUserAccount) AsMiddleWare(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		emailAddress := ctx.Value(email).(string)

		err := c.Exec(emailAddress)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		h.ServeHTTP(w, r)
	})
}
