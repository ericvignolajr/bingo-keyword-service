package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {

	/*
		by the time a user is instantiated, the password should be hashed
		TODO: decice where to enforce this in the code base. I think it
		adds unnecessary functionality to the user class to have it be
		concerned with hashing so it should be done elsewhere.
	*/
	expected := &User{
		Email:    "foo@example.com",
		Password: "supersecret",
	}
	u, err := NewUser("foo@example.com", "supersecret")
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, expected.Email, u.Email)
	assert.Equal(t, expected.Password, u.Password)
}

func TestAddSubject(t *testing.T) {

	u, err := NewUser("foo@example.com", "supersecret")
	if err != nil {
		t.Error(err)
	}

	s, err := NewSubject("Science", u.ID)
	if err != nil {
		t.Error(err)
	}
	_, err = u.AddSubject(*s)
	if err != nil {
		t.Error(err)
	}
}

func TestIsDuplicateSubject(t *testing.T) {
	u, err := NewUser("foo@example.com", "supersecret")
	if err != nil {
		t.Error(err)
	}

	s1, err := NewSubject("Electromagnets", u.ID)
	if err != nil {
		t.Error(err)
	}
	s2, err := NewSubject("Electromagnets", u.ID)
	if err != nil {
		t.Error(err)
	}

	_, err = u.AddSubject(*s1)
	if err != nil {
		t.Error(err)
	}

	isDuplicate := u.IsDuplicateSubject(*s2)
	assert.Equal(t, true, isDuplicate)
}

func TestFindSubjectByName(t *testing.T) {
	user, err := NewUser("example@bingoboard.com", "")
	if err != nil {
		t.Error(err)
	}

	s, _ := NewSubject("Science", user.ID)
	user.AddSubject(*s)

	res, err := user.FindSubjectByName(s.Name)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, s, res)
}

func TestFindSubject(t *testing.T) {
	user, err := NewUser("example@bingoboard.com", "")
	if err != nil {
		t.Error(err)
	}

	subject, _ := NewSubject("science", user.ID)
	user.AddSubject(*subject)

	actual, _ := user.FindSubject(subject.Id)

	assert.Equal(t, subject, actual)
}
