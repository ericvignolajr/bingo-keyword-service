package viewers_test

import (
	"testing"

	outputports "github.com/ericvignolajr/bingo-keyword-service/pkg/output_ports"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/viewers"
	"github.com/stretchr/testify/assert"
)

func TestViewModelToHTML(t *testing.T) {
	viewModel := outputports.WebViewModel{
		SubjectID:   "123456",
		SubjectName: "Science",
	}

	view := viewers.WebViewer{}
	actual, err := view.ViewReadSubject(viewModel)
	if err != nil {
		t.Error(err)
	}

	expected := "<h2>Science</h2><p>123456</p>"

	assert.Equal(t, expected, actual)
}
