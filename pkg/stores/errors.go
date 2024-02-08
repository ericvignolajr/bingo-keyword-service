package stores

type RecordNotFoundError struct {
	Err error
}

func (rnf *RecordNotFoundError) Error() string {
	return rnf.Err.Error()
}
