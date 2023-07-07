package usecases

type CreateSubjectPresenter interface {
	PresentCreateSubject(CreateSubjectResponse)
}

type ReadSubjectPresenter interface {
	PresentReadSubject(ReadSubjectResponse)
}
