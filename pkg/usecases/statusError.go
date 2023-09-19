package usecases

type Status int

type StatusErr struct {
	Status  Status
	Message string
}

func (s *StatusErr) Error() string {
	return s.Message
}
