package main

import (
	"fmt"
	"net/http"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/delivery"
	outputports "github.com/ericvignolajr/bingo-keyword-service/pkg/output_ports"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores/inmemory"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/usecases"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/viewers"
)

func main() {
	ustore := inmemory.UserStore{}
	uid, _ := ustore.Create("foo@example.com", "supersecret")
	fmt.Println(uid)
	sstore := inmemory.NewSubjectStore()
	c := usecases.CreateSubject{
		UserStore:    &ustore,
		SubjectStore: sstore,
		Presenter:    &outputports.MockPresenter{&viewers.MockViewer{}},
	}
	c.Exec(usecases.CreateSubjectRequest{uid, "Science"})

	v := viewers.WebViewer{}
	p := outputports.WebPresenter{
		Viewer: &v,
	}
	u := usecases.ReadSubject{
		SubjectStore: sstore,
		Presenter:    &p,
	}
	s := delivery.NewHttpServer(u, &v)
	http.ListenAndServe(":8080", s)
}
