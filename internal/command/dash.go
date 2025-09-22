package command

import (
	"strconv"

	"github.com/Ealenn/gira/internal/branch"
	"github.com/Ealenn/gira/internal/configuration"
	"github.com/Ealenn/gira/internal/git"
	"github.com/Ealenn/gira/internal/issue"
	"github.com/Ealenn/gira/internal/log"
	"github.com/Ealenn/gira/internal/ui"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	footerStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("236")).
			Foreground(lipgloss.Color("244")).
			Padding(0, 1)

	frameStyle = lipgloss.NewStyle().
			Margin(0).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240"))
)

type Dash struct {
	loader *ui.Loader

	logger  *log.Logger
	tracker issue.Tracker
	git     *git.Git
	branch  *branch.Manager
	profile *configuration.Profile

	enableAI     bool
	statusFilter string

	issues   map[string]*issue.Issue
	table    table.Model
	selected *issue.Issue
	action   string

	width        int
	height       int
	tableWidth   int
	tableHeight  int
	headerHeight int
	footerHeight int
}

func NewDashboard(logger *log.Logger, profile *configuration.Profile, tracker issue.Tracker) *Dash {
	return &Dash{
		loader: ui.NewLoader(),

		logger:       logger,
		tracker:      tracker,
		profile:      profile,
		headerHeight: 1,
		footerHeight: 1,
	}
}

func (cmd *Dash) Init() tea.Cmd {
	return tea.Batch(
		cmd.loader.Spinner.Tick,
		cmd.fetchIssues(),
	)
}

func (cmd *Dash) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var teaCmd tea.Cmd

	if cmd.loader.Loading {
		cmd.loader.Spinner, teaCmd = cmd.loader.Spinner.Update(msg)

		switch m := msg.(type) {
		case tea.WindowSizeMsg:
			cmd.width = m.Width
			cmd.height = m.Height
		case issuesLoadedMsg:
			cmd.issues = m.issues
			cmd.buildTable()
			cmd.resize()
			cmd.loader.Loading = false
		case tea.KeyMsg:
			if m.String() == "ctrl+c" || m.String() == "q" {
				return cmd, tea.Quit
			}
		}
		return cmd, teaCmd
	}

	cmd.table, teaCmd = cmd.table.Update(msg)

	if selectedRow := cmd.table.SelectedRow(); len(selectedRow) > 0 {
		if selectedIssue, ok := cmd.issues[selectedRow[0]]; ok && selectedIssue != nil {
			cmd.selected = selectedIssue
		}
	}

	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		cmd.width = m.Width
		cmd.height = m.Height
		cmd.resize()
	case tea.KeyMsg:
		switch m.String() {
		case "ctrl+c", "esc", "q":
			cmd.action = ""
			return cmd, tea.Quit
		case "enter":
			cmd.action = "enter"
			return cmd, tea.Quit
		case "o":
			if cmd.selected != nil {
				NewOpen(cmd.logger, cmd.branch, cmd.tracker).Run(&cmd.selected.ID)
			}
			return cmd, func() tea.Msg {
				return tea.WindowSizeMsg{Width: cmd.width, Height: cmd.height}
			}
		case "b":
			if cmd.selected != nil {
				cmd.action = "branch"
				return cmd, tea.Quit
			}
			return cmd, nil
		}
	}

	return cmd, teaCmd
}

func (cmd *Dash) View() string {
	if cmd.loader.Loading {
		return cmd.loader.View()
	}

	content := cmd.table.View()
	sel := cmd.table.Cursor()
	totalItems := len(cmd.issues)
	current := 0
	if totalItems > 0 {
		current = sel + 1
	}
	footerText := "ESC/Q Quit | ↑/↓ Scroll | Enter View | b Branch | o Open"
	right := strconv.Itoa(current) + "/" + strconv.Itoa(totalItems) + " "
	footer := lipgloss.JoinHorizontal(
		lipgloss.Top,
		footerStyle.Width(max(0, cmd.width-lipgloss.Width(right))).Render(footerText),
		footerStyle.Align(lipgloss.Right).Render(right),
	)

	framed := frameStyle.
		Width(cmd.width).
		Height(cmd.height - 1 - cmd.headerHeight - cmd.footerHeight).
		Render(content)

	return lipgloss.JoinVertical(lipgloss.Left, framed, footer)
}

/* ----------------------
   Runner
-----------------------*/

func (cmd *Dash) Run(dashboardStatusFlag *string, enableAI bool) {
	cmd.enableAI = enableAI

	if dashboardStatusFlag != nil {
		cmd.statusFilter = *dashboardStatusFlag
	}

	if cmd.profile.Type == configuration.ProfileTypeJira && cmd.profile.Jira.Board == "" {
		cmd.logger.Fatal("❌ %s\nYou can configure new dashboard with %s", "No dashboard configured", "gira config -p "+cmd.profile.Name)
	}

	p := tea.NewProgram(cmd, tea.WithAltScreen(), tea.WithMouseCellMotion())
	finalModel, err := p.Run()
	if err != nil {
		cmd.logger.Fatal("Gira fatal exception : %v", err)
	}

	if dash, ok := finalModel.(*Dash); ok {
		switch dash.action {
		case "branch":
			if dash.selected != nil {
				NewBranch(dash.logger, dash.tracker, dash.git, dash.branch).
					Run(dash.selected.ID, true, false, dash.enableAI)
			}
		case "enter":
			if dash.selected != nil {
				NewIssue(cmd.logger, cmd.tracker, cmd.git, cmd.branch).RunWithIssue(cmd.selected, cmd.enableAI)
			}
		}
	}
}

func (cmd *Dash) resize() {
	outerW, outerH := cmd.width, cmd.height
	usableH := max(0, outerH-cmd.headerHeight-cmd.footerHeight)

	hPad := 2 + 2
	vPad := 2 + 2

	tableW := max(20, outerW-hPad)
	tableH := max(3, usableH-vPad)

	cmd.tableWidth = tableW
	cmd.tableHeight = tableH

	idW := 10
	statusW := 14
	titleW := max(20, tableW-idW-statusW-4)

	cmd.table.SetColumns([]table.Column{
		{Title: "#", Width: idW},
		{Title: "Title", Width: titleW},
		{Title: "Status", Width: statusW},
	})

	cmd.table.SetHeight(tableH)
	cmd.table.SetWidth(tableW)
}

/*
	----------------------
	Internal helpers
  -----------------------
*/

type issuesLoadedMsg struct {
	issues map[string]*issue.Issue
}

func (cmd *Dash) fetchIssues() tea.Cmd {
	return func() tea.Msg {
		issues := cmd.tracker.SearchIssues(cmd.statusFilter)
		return issuesLoadedMsg{issues: issues}
	}
}

func (cmd *Dash) buildTable() {
	rows := make([]table.Row, 0, len(cmd.issues))
	for _, issue := range cmd.issues {
		rows = append(rows, table.Row{
			issue.ID,
			issue.Title,
			issue.Status,
		})
	}

	columns := []table.Column{
		{Title: "#", Width: 8},
		{Title: "Title", Width: 50},
		{Title: "Status", Width: 12},
	}

	cmd.table = table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true).
		Foreground(lipgloss.Color("252")).
		Background(lipgloss.Color("236"))
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("57")).
		Bold(false)
	cmd.table.SetStyles(s)
}
