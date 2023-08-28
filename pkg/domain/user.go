package domain

import (
	"fmt"
	"net/mail"

	"github.com/google/uuid"
)

type User struct {
	Id          uuid.UUID
	ExternalIDs []string
	Email       string
	Password    string
	Subjects    []*Subject
}

func NewUser(email, password string) (*User, error) {
	e, err := mail.ParseAddress(email)
	if err != nil {
		return nil, err
	}

	return &User{
		Id:          uuid.New(),
		Email:       e.Address,
		ExternalIDs: make([]string, 0),
		Password:    password,
		Subjects:    []*Subject{},
	}, nil
}

func (u *User) AddSubject(s Subject) (*Subject, error) {
	isDuplicate := u.IsDuplicateSubject(s)
	if isDuplicate {
		return nil, fmt.Errorf("user %s already has subject %s", u.Email, s.Name)
	}

	u.Subjects = append(u.Subjects, &s)
	return &s, nil
}

func (u *User) IsDuplicateSubject(s Subject) bool {
	for _, v := range u.Subjects {
		if v.Name == s.Name {
			return true
		}
	}

	return false
}
