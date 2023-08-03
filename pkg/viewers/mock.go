package viewers

import (
	"fmt"

	outputports "github.com/ericvignolajr/bingo-keyword-service/pkg/output_ports"
)

type MockViewer struct {
	view outputports.MockViewModel
}

func (m *MockViewer) ViewMock(viewModel interface{}) interface{} {
	fmt.Println(viewModel)
	m.view = viewModel.(outputports.MockViewModel)
	return nil
}

func (m *MockViewer) View(viewModel interface{}) interface{} {
	fmt.Println(viewModel)
	return m.view
}
