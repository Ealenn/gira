package forms

import (
	"github.com/charmbracelet/huh"

	"github.com/Ealenn/gira/internal/issue"
	"github.com/Ealenn/gira/internal/log"
)

type CreateIssueResult struct {
	Project     string
	Type        issue.Type
	Title       string
	Description string
}

type CreateIssue struct {
	logger *log.Logger
	ui     *huh.Form
	Result *CreateIssueResult
}

func NewCreateIssue(logger *log.Logger) *CreateIssue {
	return &CreateIssue{
		logger,
		nil,
		&CreateIssueResult{},
	}
}

func (form CreateIssue) Ask(haveProject bool) *CreateIssueResult {
	form.ui = form.getForm(haveProject)
	err := form.ui.Run()

	if err != nil {
		form.logger.Fatal("❌ The operation was %s", "canceled")
	}

	form.ui.View()
	return form.Result
}

func (form CreateIssue) getForm(haveProject bool) *huh.Form {
	var steps []*huh.Group

	if haveProject {
		steps = append(steps, huh.NewGroup(
			huh.NewInput().
				Title("Project").
				Value(&form.Result.Project),
		))
	}

	steps = append(steps, huh.NewGroup(
		huh.NewSelect[issue.Type]().
			Title("Type").
			Options(
				huh.NewOption(string(issue.TypeBug), issue.TypeBug),
				huh.NewOption(string(issue.TypeFeature), issue.TypeFeature),
			).
			Value(&form.Result.Type),
		huh.NewInput().
			Title("Title").
			Value(&form.Result.Title),
		huh.NewText().
			Title("Description").
			CharLimit(1024).
			Value(&form.Result.Description),
	))

	return huh.NewForm(steps...).WithTheme(huh.ThemeDracula())
}
