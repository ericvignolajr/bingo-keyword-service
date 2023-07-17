package viewers

import "fmt"

type MockViewer struct{}

func (*MockViewer) ViewMock(viewModel interface{}) interface{} {
	fmt.Println(viewModel)
	return nil
}

func (*MockViewer) View(viewModel interface{}) interface{} {
	fmt.Println(viewModel)
	return nil
}
