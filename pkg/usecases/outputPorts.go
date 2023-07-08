package usecases

type CreateSubjectPresenter interface {
	PresentCreateSubject(CreateSubjectResponse)
}

type ReadSubjectViewModel struct {
	Id   string
	Name string
}
type ReadSubjectPresenter interface {
	PresentReadSubject(ReadSubjectResponse) ReadSubjectViewModel
}
