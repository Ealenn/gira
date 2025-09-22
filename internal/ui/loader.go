package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

type Loader struct {
	Loading bool
	Message string
	Spinner spinner.Model
}

func NewLoader() *Loader {
	return &Loader{
		Loading: true,
		Message: "Loading ...",
		Spinner: spinner.New(
			spinner.WithSpinner(spinner.Dot),
			spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("63"))),
		),
	}
}

func (loader *Loader) View() string {
	return fmt.Sprintf("\n\n %s %s", loader.Spinner.View(), loader.Message)
}
