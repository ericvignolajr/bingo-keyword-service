package viewers

import (
	"bytes"
	"html/template"

	outputports "github.com/ericvignolajr/bingo-keyword-service/pkg/output_ports"
)

type WebViewer struct {
	view string
}

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

	w.view = b.String()
	return w.view, err
}

func (w *WebViewer) View() string {
	return w.view
}
