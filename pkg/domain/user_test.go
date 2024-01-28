package domain

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
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
	_, err = user.AddSubject(*subject)
	if err != nil {
		t.Error(err)
	}

	actual, err := user.FindSubject(subject.ID)
	if err != nil {
		t.Error(err)
	}

	if isEqual := cmp.Equal(subject, actual); isEqual != true {
		t.Errorf("%s\n", cmp.Diff(subject, actual))
	}
}

func TestEqual(t *testing.T) {
	type testCase struct {
		name     string
		user1    *User
		user2    *User
		expected bool
	}

	u1, _ := NewUser("foo@example.com", "12345")
	u2, _ := NewUser("foo@example.com", "12345")

	testCase1 := testCase{
		name:     "IDs don't match",
		user1:    u1,
		user2:    u2,
		expected: false,
	}

	testCase2 := testCase{
		name:     "one is nil",
		user1:    u1,
		user2:    nil,
		expected: false,
	}

	mismatchEmailUser, _ := NewUser("bar@example.com", "12345")
	mismatchEmailUser.ID = u1.ID
	testCase3 := testCase{
		name:     "emails don't match",
		user1:    u1,
		user2:    mismatchEmailUser,
		expected: false,
	}

	testCases := []testCase{
		testCase1,
		testCase2,
		testCase3,
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := testCase.user1.Equal(testCase.user2)
			if result != testCase.expected {
				fmt.Println(cmp.Diff(testCase.user1, testCase.user2))
				t.Errorf("on test %s, got %t, expected %t", testCase.name, result, testCase.expected)
			}
		})
	}
}

func TestDeleteSubject(t *testing.T) {
	u, err := NewUser("foo@example.com", "fake pass")
	if err != nil {
		t.Error(err)
	}

	s1, err := NewSubject("science", u.ID)
	if err != nil {
		t.Error(err)
	}
	s2, err := NewSubject("math", u.ID)
	if err != nil {
		t.Error(err)
	}

	_, err = u.AddSubject(*s1)
	if err != nil {
		t.Error(err)
	}
	_, err = u.AddSubject(*s2)
	if err != nil {
		t.Error(err)
	}

	err = u.DeleteSubject(s1.ID)
	if err != nil {
		t.Error(err)
	}

	assert.NotContains(t, u.Subjects, s1)
	_, ok := u.subjectsMap[s1.ID]
	if ok {
		t.Errorf("the key %s for subject %s should have been removed from the user's subjects map, but the key is still present", s1.ID, s1.Name)
	}

}
