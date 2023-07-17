package viewers

import (
	"bytes"
	"html/template"

	outputports "github.com/ericvignolajr/bingo-keyword-service/pkg/output_ports"
)

type WebViewer struct{}

func (w *WebViewer) ViewReadSubject(viewModel outputports.WebViewModel) (string, error) {
	b := bytes.Buffer{}
	t, err := template.New("subject").Parse("<h2>{{.SubjectName}}</h2><p>{{.SubjectID}}</p>")
	if err != nil {
		return "", err
	}

	err = t.Execute(&b, viewModel)
	if err != nil {
		return "", err
	}

	return b.String(), err
}
