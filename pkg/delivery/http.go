package delivery

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/usecases"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/viewers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwt"
)

func NewHttpServer(rs usecases.ReadSubject, signIn usecases.Signin, rv *viewers.WebViewer) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.SetHeader("X-Frame-Options", "DENY"))

	r.Route("/signin", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			wd, err := os.Getwd()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			path := filepath.Join(wd, "/pkg", "/delivery", "/signin.html")

			http.ServeFile(w, r, path)
		})
		r.Post("/with-google", func(w http.ResponseWriter, r *http.Request) {
			// verify the id token (jwt) from google here
			token, err := jwt.ParseRequest(r, jwt.WithFormKey("credential"))
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusUnauthorized)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(token.Subject() + "\n"))
			claim, ok := token.Get("email")
			if ok {
				email, ok := claim.(string)
				if ok {
					w.Write([]byte(email + "\n"))
				}
			} else {
				w.Write([]byte("Email Unknown" + "\n"))
			}
			w.Write([]byte(token.Issuer() + "\n"))
			w.Write([]byte(token.IssuedAt().String() + "\n"))
			w.Write([]byte(token.Expiration().String() + "\n"))

			// build and sign a jwt here for storing via cookie on client
			// c := http.Cookie{
			// 	Name:     "cred",
			// 	HttpOnly: true,
			// 	Secure:   true,
			// 	SameSite: http.SameSiteLaxMode,
			// 	Value:    "new cookie",
			// }
			// http.SetCookie(w, &c)
		})
	})
	r.Route("/subjects", func(r chi.Router) {
		r.Put("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotImplemented)
			return
		})
		r.Get("/{subjectID}", func(w http.ResponseWriter, r *http.Request) {
			s := chi.URLParam(r, "subjectID")
			uid, err := uuid.Parse(s)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error() + "\n"))
				w.Write([]byte(s))
				return
			}
			req := usecases.ReadSubjectRequest{
				UserId:      uid,
				SubjectName: "Science",
			}
			res := rs.Exec(req)
			if res.Err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("<h1>There was an error</h1>"))
			}

			w.Header().Add("Content-Type", "text/html")
			w.Write([]byte(rv.View()))
		})
	})

	return r
}
