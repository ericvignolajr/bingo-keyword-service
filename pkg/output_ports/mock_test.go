package outputports_test

import (
	"testing"

	outputports "github.com/ericvignolajr/bingo-keyword-service/pkg/output_ports"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores/inmemory"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/usecases"
	"github.com/stretchr/testify/assert"
)

func TestPresentReadSubject(t *testing.T) {
	userStore := inmemory.UserStore{}
	uid, err := userStore.Create("stephv@com.invalid", "fakepassword")
	if err != nil {
		t.Error(err)
	}

	subjectStore := inmemory.NewSubjectStore()
	mockPresenter := outputports.MockPresenter{}
	readSubject := usecases.ReadSubject{
		SubjectStore: subjectStore,
		Presenter:    &mockPresenter,
	}

	createSubject := usecases.CreateSubject{
		UserStore:    &userStore,
		SubjectStore: subjectStore,
		Presenter:    &mockPresenter,
	}
	createSubject.Exec(usecases.CreateSubjectRequest{
		UserId:      uid,
		SubjectName: "Science",
	})

	req := usecases.ReadSubjectRequest{
		UserId:      uid,
		SubjectName: "Science",
	}
	res := readSubject.Exec(req)
	if res.Err != nil {
		t.Error(res.Err)
	}

	actual := mockPresenter.PresentReadSubject(res)
	expected := usecases.ReadSubjectViewModel{
		Id:   res.Subject.Id.String(),
		Name: "Science",
	}

	assert.Equal(t, expected, actual)
}
