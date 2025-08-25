package forms

import (
	"github.com/charmbracelet/huh"

	"github.com/Ealenn/gira/internal/branch"
	"github.com/Ealenn/gira/internal/log"
)

type EditBranch struct {
	logger *log.Logger
	ui     *huh.Form
}

func NewEditBranch(logger *log.Logger) *EditBranch {
	return &EditBranch{
		logger,
		nil,
	}
}

func (form EditBranch) Ask(title string, description string, branch *branch.Branch) {
	form.ui = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title(title).
				Description(description).
				Value(&branch.Raw),
		),
	).WithTheme(huh.ThemeDracula())

	err := form.ui.Run()

	if err != nil {
		form.logger.Fatal("The operation was %s", "canceled")
	}

	form.ui.View()
}
