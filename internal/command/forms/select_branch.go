package forms

import (
	"github.com/Ealenn/gira/internal/branch"
	"github.com/Ealenn/gira/internal/log"
	"github.com/charmbracelet/huh"
)

type SelectBranchResult struct {
	Branch branch.Branch
}

type SelectBranch struct {
	logger *log.Logger
	ui     *huh.Form
	Result *SelectBranchResult
}

func NewSelect(logger *log.Logger) *SelectBranch {
	return &SelectBranch{
		logger,
		nil,
		&SelectBranchResult{},
	}
}

func (form SelectBranch) Ask(title string, description string, branches ...*branch.Branch) *SelectBranchResult {
	form.ui = form.getForm(title, description, branches)
	err := form.ui.Run()

	if err != nil {
		form.logger.Fatal("The operation was %s", "canceled")
	}

	form.ui.View()
	return form.Result
}

func (form SelectBranch) getForm(title string, description string, options []*branch.Branch) *huh.Form {
	var opts []huh.Option[branch.Branch]

	for _, opt := range options {
		opts = append(opts, huh.NewOption(opt.Raw, *opt))
	}

	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[branch.Branch]().
				Title(title).
				Description(description).
				Options(
					opts...,
				).
				Value(&form.Result.Branch),
		),
	).WithTheme(huh.ThemeDracula())
}
