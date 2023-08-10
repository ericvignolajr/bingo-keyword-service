package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/sessions"
	"github.com/go-chi/chi/v5"
)

const GCSRFTOKEN = "g_csrf_token"
const gCredName = "credential"
const parsedCredName = "parsedCredential"

func userRoutes(authenticator sessions.Authenticator) chi.Router {
	r := chi.NewRouter()

	r.Get("/", signin)
	r.Group(func(r chi.Router) {
		r.Use(authenticate(authenticator))
		r.Use(validateGcsrf)
		r.Post("/with-google", handleGoogleAuth)
	})

	return r
}

func signin(w http.ResponseWriter, r *http.Request) {
	wd, err := os.Getwd()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	path := filepath.Join(wd, "/pkg", "/delivery", "/signin.html")

	http.ServeFile(w, r, path)
}

func validateGcsrf(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(GCSRFTOKEN)
		if err != nil {
			fmt.Println("no csrf token in cookies")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		err = r.ParseForm()
		ok := r.PostForm.Has(GCSRFTOKEN)
		if err != nil || !ok {
			fmt.Println("no csrf token in request body")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if c.Value != r.FormValue(GCSRFTOKEN) {
			fmt.Println("cookies and body csrf tokens do not match")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func authenticate(a sessions.Authenticator) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := r.ParseForm()
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
			}
			cred := r.PostFormValue(gCredName)
			authResult := a.Authenticate(cred)
			ctx := context.WithValue(r.Context(), parsedCredName, authResult)
			r = r.WithContext(ctx)

			h.ServeHTTP(w, r)
		})
	}
}

func handleGoogleAuth(w http.ResponseWriter, r *http.Request) {
	authInfo := r.Context().Value(parsedCredName).(sessions.AuthenticatorResponse)
	res, err := json.Marshal(authInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	if !authInfo.IsAuthenticated {
		w.WriteHeader(http.StatusUnauthorized)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
