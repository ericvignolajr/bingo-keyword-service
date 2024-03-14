package sql_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores/sql"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestReadById(t *testing.T) {
	uStore, err := sql.NewSQLUserStore()
	if err != nil {
		t.Error(err)
	}

	user, _ := domain.NewUser("foo@example.com", "foobar")
	subj, _ := domain.NewSubject("science", user.ID)
	unit, _ := domain.NewUnit("electro-magnets")
	subj.AddUnit(*unit)
	user.AddSubject(subj)
	uStore.DB.Create(user)

	cases := []struct {
		name     string
		id       uuid.UUID
		expected *domain.User
		err      error
	}{
		{"user found", user.ID, user, nil},
		{"user not found", uuid.Nil, nil, &stores.RecordNotFoundError{}},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			got, err := uStore.ReadById(v.id)
			if cmp.Equal(v.expected, got) == false {
				cmp.Diff(v.expected, got)
				t.Errorf("expected %v, got %v", v.expected, got)
			}

			switch err.(type) {
			case *stores.RecordNotFoundError:
				var recordNotFoundErr *stores.RecordNotFoundError
				if errors.As(v.err, &recordNotFoundErr) == false {
					t.Errorf("expected error %T got error %T", v.err, err)
				}
			case nil:
				if v.err != nil {
					t.Errorf("expected error %T, got nil error", v.err)
				}
			default:
				t.Errorf("unexpected error %T", err)
			}
		})
	}

	userFromDB, err := uStore.ReadById(user.ID)
	if err != nil {
		t.Error(err)
	}

	isEqual := cmp.Equal(user, userFromDB)
	assert.Equal(t, true, isEqual)
}

func TestReadByEmail(t *testing.T) {
	uStore, err := sql.NewSQLUserStore()
	if err != nil {
		t.Error(err)
	}

	userEmail := "test@example.com"
	user, _ := domain.NewUser(userEmail, "foobaz")
	uStore.DB.Create(user)

	cases := []struct {
		name     string
		email    string
		expected *domain.User
		err      error
	}{
		{"user found", userEmail, user, nil},
		{"user not found", "doesnotexist@test.com.invalid", nil, &stores.RecordNotFoundError{}},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			got, err := uStore.ReadByEmail(v.email)
			if cmp.Equal(got, v.expected) == false {
				fmt.Println(cmp.Diff(got, v.expected))
				t.Errorf("expected %v, got %v", v.expected, got)
			}

			switch err.(type) {
			case *stores.RecordNotFoundError:
				var recordNotFoundErr *stores.RecordNotFoundError
				if errors.As(v.err, &recordNotFoundErr) == false {
					t.Errorf("expected error %T got error %T", v.err, err)
				}
			case nil:
				if v.err != nil {
					t.Errorf("expected error %T, got nil error", v.err)
				}
			default:
				t.Errorf("unexpected error %T", err)
			}
		})
	}

	userFromDBByEmail, err := uStore.ReadByEmail(userEmail)
	if err != nil {
		t.Error(err)
	}

	if cmp.Equal(user, userFromDBByEmail) == false {
		fmt.Printf(cmp.Diff(user, userFromDBByEmail))
		t.FailNow()
	}
}

func TestSave(t *testing.T) {
	user, _ := domain.NewUser("foo@example.com", "baz")
	uStore, _ := sql.NewSQLUserStore()
	_, err := uStore.Save(user)
	if err != nil {
		t.Error(err)
	}

	userFromDB, _ := uStore.ReadById(user.ID)
	assert.Equal(t, true, cmp.Equal(user, userFromDB))
	newEmail := "updated@example.com"
	userFromDB.Email = newEmail

	uStore.Save(userFromDB)

	userAfterUpdate, _ := uStore.ReadById(user.ID)

	assert.Equal(t, newEmail, userAfterUpdate.Email)
}
