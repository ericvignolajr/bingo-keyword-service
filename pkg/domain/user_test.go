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
	s, err := NewSubject("Science")
	if err != nil {
		t.Error(err)
	}

	u, err := NewUser("foo@example.com", "supersecret")
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

	s1, err := NewSubject("Electromagnets")
	if err != nil {
		t.Error(err)
	}
	s2, err := NewSubject("Electromagnets")
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
