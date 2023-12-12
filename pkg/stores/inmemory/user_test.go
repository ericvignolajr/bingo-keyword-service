package inmemory_test

import (
	"testing"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores/inmemory"
	"github.com/stretchr/testify/assert"
)

// func (u *UserStore) ReadById(id uuid.UUID) (*domain.User, error) {
// 	for i, v := range u.store {
// 		if v.Id == id {
// 			return &u.store[i], nil
// 		}
// 	}

// 	return nil, fmt.Errorf("%s, userID: %s", stores.ErrUserDoesNotExist, id)
// }

func TestSave(t *testing.T) {
	uStore := inmemory.UserStore{}
	uID, err := uStore.Create("bingoboard@example.com", "")
	if err != nil {
		t.Error(err)
	}

	u, err := uStore.ReadById(uID)
	if err != nil {
		t.Error(err)
	}

	newSubject, err := domain.NewSubject("science", u.ID)
	if err != nil {
		t.Error(err)
	}

	assert.NotContains(t, u.Subjects, newSubject)

	_, err = u.AddSubject(*newSubject)
	if err != nil {
		t.Error(err)
	}

	_, err = uStore.Save(u)
	if err != nil {
		t.Error(err)
	}
	updatedUser, err := uStore.ReadById(u.ID)
	if err != nil {
		t.Error(err)
	}

	assert.Contains(t, updatedUser.Subjects, newSubject)
	assert.Equal(t, u.ID, updatedUser.ID)
}
