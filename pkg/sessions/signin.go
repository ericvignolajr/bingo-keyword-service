package sessions

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores"
	store "github.com/ericvignolajr/bingo-keyword-service/pkg/stores/sql"
)

const (
	ExternalID = iota + 1
	UserID
	Email
	User
)

var UserStore, _ = store.NewSQLUserStore()

func AddUserToContext(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := clerk.SessionFromContext(ctx)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			fmt.Println("could not retrieve session from context")
			return
		}

		clerkAuth := NewClerkAuthenticator()
		u, _ := clerkAuth.Client.Users().Read(claims.Subject)
		userEmail := u.EmailAddresses[0].EmailAddress

		user, err := UserStore.ReadByEmail(userEmail)
		if err != nil {
			var recordNotFoundErr *stores.RecordNotFoundError
			if errors.As(err, &recordNotFoundErr) {
				fmt.Printf("user with email address: %s, was not found in the database. attempting to create new user\n", userEmail)
			} else {
				fmt.Printf("unexpected fatal error when trying to retrieve user with email address: %s\n", userEmail)
				fmt.Print(err)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				return
			}
		}
		if user == nil {
			uID, err := UserStore.Create(userEmail, "")
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				fmt.Println("Failed to create user")
				return
			}

			user, err = UserStore.ReadById(uID)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				fmt.Println("Failed to read by id")
				return
			}
		}

		ctx = context.WithValue(ctx, User, user)
		r = r.WithContext(ctx)
		fmt.Println("GetClaimsFromClerkToken - ok")
		h.ServeHTTP(w, r)
	})
}
