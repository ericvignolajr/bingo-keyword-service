package domain

import (
	"fmt"
	"net/mail"
	"strings"

	"github.com/google/uuid"
)

const (
	ErrSubjectDoesNotExist = "subject does not exist"
)

type User struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid"`
	ExternalIDs []string  `gorm:"type:string[]"`
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
		ID:          uuid.New(),
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

func (u *User) FindSubject(sID uuid.UUID) (*Subject, error) {
	for i, v := range u.Subjects {
		if v.Id == sID {
			return u.Subjects[i], nil
		}
	}

	return nil, fmt.Errorf("%s, subjectID: %s", ErrSubjectDoesNotExist, sID)
}

func (u *User) FindSubjectByName(subjectName string) (*Subject, error) {
	for i, v := range u.Subjects {
		if strings.EqualFold(v.Name, subjectName) {
			return u.Subjects[i], nil
		}
	}

	return nil, fmt.Errorf("%s, subjectName: %s", ErrSubjectDoesNotExist, subjectName)
}

func (u *User) DeleteSubject(subjectID uuid.UUID) error {
	var subjectSlice = u.Subjects
	var nullIndex struct {
		index int
		found bool // found is true if the index was set
	}
	for idx, subject := range subjectSlice {
		if subject.Id == subjectID {
			nullIndex = struct {
				index int
				found bool
			}{
				index: idx,
				found: true,
			}
		}
	}

	if !nullIndex.found {
		return nil
	}

	length := len(subjectSlice)
	subjectSlice[nullIndex.index] = subjectSlice[length-1]
	subjectSlice[length-1] = nil
	u.Subjects = subjectSlice[:length-1]
	return nil
}

func (u *User) IsDuplicateSubject(s Subject) bool {
	for _, v := range u.Subjects {
		if v.Name == s.Name {
			return true
		}
	}

	return false
}
