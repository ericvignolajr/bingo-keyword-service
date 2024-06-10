package stores

type ErrRecordNotFound struct {
	Err error
}

func (rnf *ErrRecordNotFound) Error() string {
	return rnf.Err.Error()
}
