package usecases

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores"
	"github.com/google/uuid"
)

type ReadSubjects struct {
	UserStore stores.User
}

func (r *ReadSubjects) Exec(userID uuid.UUID, subjectID *uuid.UUID) ([]domain.Subject, error) {
	user, err := r.UserStore.ReadById(userID)
	if err != nil {
		return nil, fmt.Errorf("in readSubjects: %w", err)
		return nil, fmt.Errorf("in readSubjects: %w", err)
	}

	subjectOutput := make([]domain.Subject, 0, len(user.Subjects))
	if subjectID != nil {
		for _, v := range user.Subjects {
			if v.ID == *subjectID {
				// there should only be one subject with id matching subjectID
				subjectOutput = append(subjectOutput, *v)
				break
			}
		}
	} else {
		for _, v := range user.Subjects {
			subjectOutput = append(subjectOutput, *v)
		}
	}

	sort.Slice(subjectOutput, func(i, j int) bool {
		return strings.ToLower(subjectOutput[i].Name) < strings.ToLower(subjectOutput[j].Name)
	})

	return subjectOutput, nil
	return subjectOutput, nil
}
