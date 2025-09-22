package command

import (
	"fmt"
	"time"

	"github.com/Ealenn/gira/internal/ai"
	"github.com/Ealenn/gira/internal/branch"
	"github.com/Ealenn/gira/internal/git"
	"github.com/Ealenn/gira/internal/issue"
	"github.com/Ealenn/gira/internal/log"
	"github.com/Ealenn/gira/internal/ui"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type Issue struct {
	loader *ui.Loader

	logger  *log.Logger
	tracker issue.Tracker
	git     *git.Git
	branch  *branch.Manager

	enableAI bool
	agent    *ai.OpenAI

	aiSummary string
	issueID   string
	issue     *issue.Issue

	action string
	width  int
	height int

	focus                    string
	componentAttributes      viewport.Model
	componentAttributesValue string
	componentContent         viewport.Model
	componentContentValue    string
}

func NewIssue(logger *log.Logger, tracker issue.Tracker, git *git.Git, branch *branch.Manager) *Issue {
	return &Issue{
		loader: ui.NewLoader(),

		logger:  logger,
		tracker: tracker,
		git:     git,
		branch:  branch,

		issue:               nil,
		componentAttributes: viewport.New(0, 0),
		componentContent:    viewport.New(0, 0),
	}
}

func (cmd *Issue) Init() tea.Cmd {
	return tea.Batch(
		cmd.loader.Spinner.Tick,
		cmd.fetchIssue(),
	)
}

func (cmd *Issue) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var teaCmd tea.Cmd

	if cmd.loader.Loading {
		cmd.loader.Spinner, teaCmd = cmd.loader.Spinner.Update(msg)

		switch m := msg.(type) {
		case tea.WindowSizeMsg:
			cmd.width = m.Width
			cmd.height = m.Height
			cmd.resize()
		case issueLoadedMsg:
			cmd.issue = m.issue
			cmd.aiSummary = m.aiSummary
			cmd.loader.Loading = false
			cmd.render()
			return cmd, func() tea.Msg {
				return tea.WindowSizeMsg{Width: cmd.width, Height: cmd.height}
			}
		case tea.KeyMsg:
			if m.String() == "ctrl+c" || m.String() == "q" {
				return cmd, tea.Quit
			}
		}
		return cmd, teaCmd
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		cmd.width = msg.Width
		cmd.height = msg.Height
		cmd.resize()
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
		case "ctrl+c", "esc", "q":
			cmd.action = "quit"
			return cmd, tea.Quit
		case "o":
			cmd.action = "open"
			return cmd, tea.Quit
		case "a":
			cmd.action = "assign"
			return cmd, tea.Quit
		case "b":
			cmd.action = "branch"
			return cmd, tea.Quit
		}
	}

	cmd.componentAttributes.Update(msg)
	cmd.componentContent.Update(msg)

	switch cmd.focus {
	case "attributes":
		cmd.componentAttributes, teaCmd = cmd.componentAttributes.Update(msg)
	case "content":
		cmd.componentContent, teaCmd = cmd.componentContent.Update(msg)
	}

	return cmd, teaCmd
}

func (cmd *Issue) View() string {
	if cmd.loader.Loading {
		return cmd.loader.View()
	}

	if cmd.action != "" {
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

	helpBar := helpStyle.Render("ESC/Q Quit | ↑/↓ Scroll | a Assign | b Branch | o Open")
	return lipgloss.JoinVertical(lipgloss.Left, mainContent, helpBar)
}

/*
	----------------------
	  Runner
  -----------------------
*/

func (cmd *Issue) RunWithIssue(issue *issue.Issue, enableAI bool) {
	cmd.issue = issue
	cmd.Run(&cmd.issue.ID, enableAI)
}

func (cmd *Issue) Run(optionalIssueID *string, enableAI bool) {
	cmd.enableAI = enableAI
	if cmd.enableAI {
		cmd.agent = ai.NewOpenAI(cmd.logger)
	}

	if optionalIssueID != nil {
		cmd.issueID = *optionalIssueID
	} else {
		cmd.issueID = cmd.branch.GetCurrentBranch().IssueID
	}

	p := tea.NewProgram(cmd,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		cmd.logger.Fatal("Gira fatal exception : %v", err)
	}

	switch cmd.action {
	case "open":
		NewOpen(cmd.logger, cmd.branch, cmd.tracker).Run(&cmd.issue.ID)
	case "assign":
		cmd.tracker.SelfAssignIssue(cmd.issue.ID)
		NewIssue(cmd.logger, cmd.tracker, cmd.git, cmd.branch).RunWithIssue(cmd.issue, enableAI)
	case "branch":
		NewBranch(cmd.logger, cmd.tracker, cmd.git, cmd.branch).Run(cmd.issue.ID, false, enableAI, false)
	}
}

/*
	----------------------
	  Messages
  -----------------------
*/

type issueLoadedMsg struct {
	issue     *issue.Issue
	aiSummary string
}

func (cmd *Issue) fetchIssue() tea.Cmd {
	return func() tea.Msg {
		var (
			issue     *issue.Issue
			aiSummary string
		)

		if cmd.issue == nil {
			issue = cmd.tracker.GetIssue(cmd.issueID)
		} else {
			issue = cmd.issue
		}

		aiSummary = ""
		if cmd.enableAI {
			if response, err := cmd.agent.IssueSummary(issue); err == nil {
				aiSummary = response
			}
		}

		return issueLoadedMsg{issue: issue, aiSummary: aiSummary}
	}
}

func (cmd *Issue) resize() {
	mainHeight := max(cmd.height-3, 3)

	cmd.componentAttributes.Width = 25
	cmd.componentAttributes.Height = (mainHeight) - 2

	cmd.componentContent.Width = max(cmd.width-25-2, 30)
	cmd.componentContent.Height = mainHeight
}

func (cmd *Issue) render() {
	cmd.componentContentValue = fmt.Sprintf("# %s\n\r\n\r%s", cmd.issue.Title, cmd.issue.Description)

	if cmd.enableAI && cmd.aiSummary != "" {
		cmd.componentContentValue = fmt.Sprintf("# %s\n> 🤖 %s\n\r\n\r---\n\r\n\r%s", cmd.issue.Title, cmd.aiSummary, cmd.issue.Description)
	}

	cmd.componentAttributesValue = fmt.Sprintf("> #%s\n\nStatus: %s\n\n", cmd.issue.ID, cmd.issue.Status)
	cmd.componentAttributesValue += "\n> Types \n\n"
	for _, tag := range cmd.issue.Types {
		cmd.componentAttributesValue += fmt.Sprintf("- %s\n\n", tag)
	}
	cmd.componentAttributesValue += "\n> Assignees\n\n"
	for _, assignee := range cmd.issue.Assignees {
		cmd.componentAttributesValue += fmt.Sprintf("- [%s](%s) \n\n", assignee.Name, assignee.Email)
	}
	cmd.componentAttributesValue += fmt.Sprintf("\n\n%s", cmd.issue.CreatedAt.Format(time.RFC822))

	cmd.componentAttributes.SetContent(cmd.renderMarkdown(cmd.componentAttributesValue, cmd.componentAttributes.Width))
	cmd.componentContent.SetContent(cmd.renderMarkdown(cmd.componentContentValue, cmd.componentContent.Width))
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
