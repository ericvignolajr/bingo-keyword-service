package outputports

import (
	"github.com/ericvignolajr/bingo-keyword-service/pkg/usecases"
)

type WebPresenter struct {
	Viewer WebViewer
}

type WebViewer interface {
	ViewReadSubject(WebViewModel) (string, error)
}

type WebViewModel struct {
	SubjectID   string
	SubjectName string
}

func (w *WebPresenter) PresentCreateSubject(res usecases.CreateSubjectResponse) {
	panic("PresentCreateSubject not implemented on WebViewer")
}

func (w *WebPresenter) PresentReadSubject(res usecases.ReadSubjectResponse) error {
	viewModel := WebViewModel{
		SubjectID:   res.Subject.Id.String(),
		SubjectName: res.Subject.Name,
	}

	w.Viewer.ViewReadSubject(viewModel)
	return nil
}
