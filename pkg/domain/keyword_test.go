package domain

import (
	_ "image/jpeg"
	_ "image/png"
	"testing"

	"github.com/google/uuid"
)

func TestNewKeyword(t *testing.T) {
	expected := &Keyword{
		ID:           uuid.New(),
		Name:         "magent",
		Definition:   "magnet defintion",
		Translations: nil,
	}

	k, err := NewKeyword(expected.Name, expected.Definition)
	if err != nil {
		t.Error(err)
	}

	if expected.Name != k.Name || expected.Definition != k.Definition {
		t.Error("expected keyword should match the one returned by factory function")
	}
}
