package domain

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrSubjectDoesNotExist = errors.New("subject does not exist")
)

type User struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid"`
	ExternalIDs []string  `gorm:"type:string[]"`
	Email       string
	Password    string
	Subjects    []*Subject
	subjectsMap map[uuid.UUID]int `gorm:"-:all"` // tell gorm to ignore this field
}

func NewUser(email, password string) (*User, error) {
	e, err := mail.ParseAddress(email)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:          uuid.New(),
		Email:       e.Address,
		ExternalIDs: nil,
		Password:    password,
		Subjects:    nil,
		subjectsMap: make(map[uuid.UUID]int),
	}, nil
}

func (u *User) AddSubject(s *Subject) error {
	isDuplicate := u.IsDuplicateSubject(*s)
	if isDuplicate {
		return fmt.Errorf("user %s already has subject %s", u.Email, s.Name)
	}

	u.Subjects = append(u.Subjects, s)
	u.subjectsMap[s.ID] = len(u.Subjects) - 1
	return nil
}

func (u *User) FindSubject(sID uuid.UUID) (*Subject, error) {
	subjectIdx, ok := u.subjectsMap[sID]
	if !ok {
		return nil, fmt.Errorf("%w, subjectID: %s", ErrSubjectDoesNotExist, sID)
	}

	if subjectIdx >= 0 && subjectIdx < len(u.Subjects) {
		if u.Subjects[subjectIdx].ID == sID {
			return u.Subjects[subjectIdx], nil
		}
	} else {
		return nil, fmt.Errorf("tried to index subjects slice with index %d, out of range, check subjectsMap", subjectIdx)
	}

	return nil, fmt.Errorf("%w, subjectID: %s", ErrSubjectDoesNotExist, sID)
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
	subjectIdx, ok := u.subjectsMap[subjectID]
	if !ok {
		return nil
	}

	length := len(u.Subjects)
	elementToDelete := u.Subjects[subjectIdx]
	lastElement := u.Subjects[length-1]

	u.subjectsMap[lastElement.ID] = subjectIdx
	u.Subjects[subjectIdx] = lastElement
	delete(u.subjectsMap, elementToDelete.ID)
	u.Subjects[length-1] = nil
	u.Subjects = u.Subjects[:length-1]

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

func (u *User) Equal(other *User) bool {
	if u == nil || other == nil {
		if u == nil && other == nil {
			return true
		}
		return false
	}

	if u.ID != other.ID {
		return false
	}
	if u.Email != other.Email {
		return false
	}
	if len(u.ExternalIDs) != len(other.ExternalIDs) {
		return false
	}
	for i := range u.ExternalIDs {
		if u.ExternalIDs[i] != other.ExternalIDs[i] {
			return false
		}
	}
	return true
}

func (u *User) AfterFind(tx *gorm.DB) error {
	if u.subjectsMap == nil {
		u.subjectsMap = make(map[uuid.UUID]int)
	}

	for i, v := range u.Subjects {
		u.subjectsMap[v.ID] = i
	}

	return nil
}
