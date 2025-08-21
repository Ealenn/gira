package command

import (
	"fmt"
	"time"

	"github.com/Ealenn/gira/internal/ai"
	"github.com/Ealenn/gira/internal/branch"
	"github.com/Ealenn/gira/internal/git"
	"github.com/Ealenn/gira/internal/issue"
	"github.com/Ealenn/gira/internal/log"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type Issue struct {
	logger  *log.Logger
	tracker issue.Tracker
	git     *git.Git
	branch  *branch.Manager

	width    int
	height   int
	quitting bool
	open     bool

	focus                    string
	componentAttributes      viewport.Model
	componentAttributesValue string
	componentContent         viewport.Model
	componentContentValue    string
}

func NewIssue(logger *log.Logger, tracker issue.Tracker, git *git.Git, branch *branch.Manager) *Issue {
	return &Issue{
		logger:              logger,
		tracker:             tracker,
		git:                 git,
		branch:              branch,
		componentAttributes: viewport.New(0, 0),
		componentContent:    viewport.New(0, 0),
	}
}

func (cmd *Issue) Init() tea.Cmd {
	return nil
}

func (cmd *Issue) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var teacmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		cmd.width = msg.Width
		cmd.height = msg.Height
		cmd.renderComponents()

	case tea.MouseMsg:
		if msg.Action == tea.MouseActionPress {
			if msg.X < cmd.componentAttributes.Width+2 {
				cmd.focus = "attributes"
			} else {
				cmd.focus = "content"
			}
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			cmd.quitting = true
			return cmd, tea.Quit
		case "o":
			cmd.open = true
			return cmd, tea.Quit
		}
	}

	cmd.componentAttributes.Update(msg)
	cmd.componentContent.Update(msg)

	switch cmd.focus {
	case "attributes":
		cmd.componentAttributes, teacmd = cmd.componentAttributes.Update(msg)
	case "content":
		cmd.componentContent, teacmd = cmd.componentContent.Update(msg)
	}

	return cmd, teacmd
}

func (cmd *Issue) View() string {
	if cmd.quitting {
		return ""
	}

	mainHeight := max(cmd.height-3, 3)

	// Styles
	attributesStyle := lipgloss.NewStyle().
		Width(cmd.componentAttributes.Width).
		Height(cmd.componentAttributes.Height + 2).
		Border(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("35"))

	contentStyle := lipgloss.NewStyle().
		Width(cmd.componentContent.Width - 2).
		Height(mainHeight).
		Border(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("35"))

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Background(lipgloss.Color("235")).
		Height(1).
		Width(cmd.width).
		Align(lipgloss.Center)

	// Focus
	switch cmd.focus {
	case "attributes":
		attributesStyle = attributesStyle.Border(lipgloss.DoubleBorder())
	case "content":
		contentStyle = contentStyle.Border(lipgloss.DoubleBorder())
	}

	// Render boxes
	leftBox := lipgloss.JoinVertical(
		lipgloss.Left,
		attributesStyle.Render(cmd.componentAttributes.View()),
	)
	leftBox = lipgloss.NewStyle().
		Height(mainHeight).
		Render(leftBox)

	rightBox := contentStyle.Render(cmd.componentContent.View())
	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)

	helpBar := helpStyle.Render("ESC/CTRL+C Quit | ↑/↓ Scroll | o Open ")
	return lipgloss.JoinVertical(lipgloss.Left, mainContent, helpBar)
}

/* ----------------------
   Runner
-----------------------*/

func (cmd *Issue) Run(optionalIssueID *string, enableAI bool) {
	var issueID string
	if optionalIssueID != nil {
		issueID = *optionalIssueID
	} else {
		issueID = cmd.branch.GetCurrentBranch().IssueID
	}
	issue := cmd.tracker.GetIssue(issueID)

	cmd.componentContentValue = fmt.Sprintf("# %s\n\r\n\r%s", issue.Title, issue.Description)
	if enableAI {
		agent := ai.NewOpenAI(cmd.logger)
		response, err := agent.IssueSummary(issue)

		if err == nil {
			cmd.componentContentValue = fmt.Sprintf("# %s\n> 🤖 %s\n\r\n\r---\n\r\n\r%s", issue.Title, response, issue.Description)
		}
	}

	cmd.componentAttributesValue = fmt.Sprintf("# %s (%s)\n\r", issue.ID, issue.Status)
	cmd.componentAttributesValue += "\n> Types \n\n"
	for _, tag := range issue.Types {
		cmd.componentAttributesValue += fmt.Sprintf("- %s\n\n", tag)
	}
	cmd.componentAttributesValue += "\n> Assignees\n\n"
	for _, assignee := range issue.Assignees {
		cmd.componentAttributesValue += fmt.Sprintf("- [%s](%s) \n\n", assignee.Name, assignee.Email)
	}

	cmd.componentAttributesValue += fmt.Sprintf("\n\n%s", issue.CreatedAt.Format(time.RFC822))

	p := tea.NewProgram(cmd,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		cmd.logger.Fatal("Gira fatal exception : %v", err)
	}

	if cmd.open {
		NewOpen(cmd.logger, cmd.branch, cmd.tracker).Run(optionalIssueID)
	}
}

func (cmd *Issue) renderComponents() {
	mainHeight := max(cmd.height-3, 3)

	cmd.componentAttributes.Width = 25
	cmd.componentAttributes.Height = (mainHeight) - 2

	cmd.componentContent.Width = max(cmd.width-25-2, 30)
	cmd.componentContent.Height = mainHeight

	cmd.componentAttributes.SetContent(cmd.renderMarkdown(cmd.componentAttributesValue, cmd.componentAttributes.Width))
	cmd.componentContent.SetContent(cmd.renderMarkdown(cmd.componentContentValue, cmd.componentContent.Width))

	cmd.Update(nil)
	cmd.View()
}

func (cmd *Issue) renderMarkdown(markdown string, wrap int) string {
	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(wrap),
		glamour.WithEmoji(),
	)
	out, _ := renderer.Render(markdown)

	return out
}
