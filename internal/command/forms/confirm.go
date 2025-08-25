package forms

import (
	"github.com/charmbracelet/huh"

	"github.com/Ealenn/gira/internal/log"
)

type ConfirmResult struct {
	Confirmed bool
}

type ConfirmType string

const (
	TypeYesNo   ConfirmType = "YES_NO"
	TypeConfirm ConfirmType = "CONFIRM"
	TypeApply   ConfirmType = "APPLY"
)

type Confirm struct {
	logger *log.Logger
	ui     *huh.Form
	Result *ConfirmResult
}

func NewConfirm(logger *log.Logger) *Confirm {
	return &Confirm{
		logger,
		nil,
		&ConfirmResult{},
	}
}

func (form Confirm) Ask(title string, description string, confirmType ConfirmType) *ConfirmResult {
	form.ui = form.getForm(title, description, confirmType)
	err := form.ui.Run()

	if err != nil {
		form.logger.Fatal("The operation was %s", "canceled")
	}

	form.ui.View()
	return form.Result
}

func (form Confirm) getForm(title string, description string, confirmType ConfirmType) *huh.Form {
	var (
		affirmativeText string
		negativeText    string
	)

	switch confirmType {
	case TypeYesNo:
		affirmativeText = "Yes!"
		negativeText = "No."
	case TypeConfirm:
		affirmativeText = "Confirm"
		negativeText = "Cancel"
	case TypeApply:
		affirmativeText = "Apply"
		negativeText = "Ignore"
	}

	return huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(title).
				Description(description).
				Affirmative(affirmativeText).
				Negative(negativeText).
				Value(&form.Result.Confirmed),
		),
	).WithTheme(huh.ThemeDracula())
}
