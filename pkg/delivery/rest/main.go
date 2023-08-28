package rest

import (
	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/sessions"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores/inmemory"
	usecases "github.com/ericvignolajr/bingo-keyword-service/pkg/usecases/signin"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewServer() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.SetHeader("X-Frame-Options", "DENY"))
	var clerkAuth = sessions.NewClerkAuthenticator()
	requireSession := clerk.RequireSessionV2(clerkAuth.Client)

	r.Get("/", usecases.SigninHandler)
	r.Get("/signin", usecases.SigninHandler)

	userStore := inmemory.UserStore{}
	readUserFeat := usecases.ReadUserByEmail{
		UserStore: &userStore,
	}
	createUserAccountFeat := usecases.CreateUserAccount{
		UserStore: &userStore,
	}

	r.Group(func(r chi.Router) {
		r.Use(requireSession)
		r.Use(usecases.GetClaimsFromClerkToken)
		r.Use(createUserAccountFeat.AsMiddleWare)
		r.Get("/register", usecases.RegisterHandler)
	})
	r.Group(func(r chi.Router) {
		r.Use(requireSession)
		r.Use(usecases.GetClaimsFromClerkToken)
		r.Use(readUserFeat.AsMiddleWare)
		r.Get("/dashboard", usecases.MyInfoHandler)
	})

	return r
}
