package main

import (
	"net/http"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/delivery"
	outputports "github.com/ericvignolajr/bingo-keyword-service/pkg/output_ports"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores/inmemory"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/usecases"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/viewers"
)

func main() {
	ustore := inmemory.UserStore{}
	sstore := inmemory.NewSubjectStore()

	v := viewers.WebViewer{}
	p := outputports.WebPresenter{
		Viewer: &v,
	}
	u := usecases.ReadSubject{
		SubjectStore: sstore,
		Presenter:    &p,
	}
	signIn := usecases.Signin{
		UserStore: ustore,
	}
	s := delivery.NewHttpServer(u, signIn, &v)
	http.ListenAndServe(":8080", s)
}
