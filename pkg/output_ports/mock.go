package outputports

import (
	"fmt"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/usecases"
)

type MockViewer interface {
	ViewMock(interface{}) interface{}
}

type MockPresenter struct {
	Viewer MockViewer
}

type MockViewModel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (p *MockPresenter) PresentCreateSubject(res usecases.CreateSubjectResponse) {
	fmt.Printf("NOW PRESENTING: %+v", res)
}

func (p *MockPresenter) PresentReadSubject(res usecases.ReadSubjectResponse) error {
	viewModel := MockViewModel{
		ID:   res.Subject.Id.String(),
		Name: res.Subject.Name,
	}

	p.Viewer.ViewMock(viewModel)
	return nil
}
