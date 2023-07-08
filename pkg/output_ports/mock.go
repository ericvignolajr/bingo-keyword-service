package outputports

import (
	"fmt"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/usecases"
)

type MockPresenter struct{}

func (p *MockPresenter) PresentCreateSubject(res usecases.CreateSubjectResponse) {
	fmt.Printf("NOW PRESENTING: %+v", res)
}

func (p *MockPresenter) PresentReadSubject(res usecases.ReadSubjectResponse) usecases.ReadSubjectViewModel {
	viewModel := usecases.ReadSubjectViewModel{
		Id:   res.Subject.Id.String(),
		Name: res.Subject.Name,
	}
	return viewModel
}
