package command

import (
	"strconv"

	"github.com/Ealenn/gira/internal/branch"
	"github.com/Ealenn/gira/internal/configuration"
	"github.com/Ealenn/gira/internal/git"
	"github.com/Ealenn/gira/internal/issue"
	"github.com/Ealenn/gira/internal/log"

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
	logger  *log.Logger
	tracker issue.Tracker
	git     *git.Git
	branch  *branch.Manager
	profile *configuration.Profile

	enableAI bool
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
		logger:       logger,
		tracker:      tracker,
		profile:      profile,
		headerHeight: 1,
		footerHeight: 1,
	}
}

func (cmd Dash) Init() tea.Cmd { return nil }

func (cmd Dash) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var teacmd tea.Cmd
	cmd.table, teacmd = cmd.table.Update(msg)

	if selectedIssue, ok := cmd.issues[cmd.table.SelectedRow()[0]]; ok {
		cmd.selected = selectedIssue
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
			if cmd.selected != nil {
				NewIssue(cmd.logger, cmd.tracker, cmd.git, cmd.branch).RunWithIssue(cmd.selected, cmd.enableAI)
			}
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

	return cmd, teacmd
}

func (cmd Dash) View() string {
	content := cmd.table.View()

	sel := cmd.table.Cursor()
	totalItems := len(cmd.issues)
	footerText := "ESC/Q Quit | ↑/↓ Scroll | Enter View | b Branch | o Open"
	right := strconv.Itoa(sel+1) + "/" + strconv.Itoa(totalItems) + " "
	footer := lipgloss.JoinHorizontal(
		lipgloss.Top,
		footerStyle.Width(max(0, cmd.width-lipgloss.Width(right))).Render(footerText),
		footerStyle.Align(lipgloss.Right).Render(right),
	)

	// Frame around content to improve aesthetic
	framed := frameStyle.Width(cmd.width).Height(cmd.height - 1 - cmd.headerHeight - cmd.footerHeight).Render(content)

	// Vertical layout
	ui := lipgloss.JoinVertical(lipgloss.Left, framed, footer)
	return ui
}

/* ----------------------
   Runner
-----------------------*/

func (cmd *Dash) Run(dashboardStatusFlag *string, enableAI bool) {
	cmd.enableAI = enableAI
	if cmd.profile.Type == configuration.ProfileTypeJira && cmd.profile.Jira.Board == "" {
		cmd.logger.Fatal("❌ %s\nYou can configure new dashboard with %s", "No dashboard configured", "gira config -p "+cmd.profile.Name)
	}

	cmd.issues = cmd.tracker.SearchIssues(*dashboardStatusFlag)

	// Build rows
	rows := make([]table.Row, 0, len(cmd.issues))
	for _, issue := range cmd.issues {
		rows = append(rows, table.Row{
			issue.ID,
			issue.Title,
			issue.Status,
		})
	}

	// Initial columns (will be resized on first WindowSizeMsg)
	columns := []table.Column{
		{Title: "#", Width: 8},
		{Title: "Title", Width: 50},
		{Title: "Status", Width: 12},
	}

	// Create table
	cmd.table = table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10), // temporary; will be recalculated on first resize
	)

	// Style the table nicely
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true).
		Foreground(lipgloss.Color("252")).
		Background(lipgloss.Color("236"))

	s.Selected = s.Selected.
		Foreground(lipgloss.Color("230")). // light
		Background(lipgloss.Color("57")).  // indigo
		Bold(false)

	cmd.table.SetStyles(s)

	// Run
	p := tea.NewProgram(cmd, tea.WithAltScreen(), tea.WithMouseCellMotion())
	finalModel, err := p.Run()
	if err != nil {
		cmd.logger.Fatal("Gira fatal exception : %v", err)
	}

	// Actions
	if dash, ok := finalModel.(Dash); ok {
		switch dash.action {
		case "branch":
			if dash.selected != nil {
				NewBranch(dash.logger, dash.tracker, dash.git, dash.branch).
					Run(dash.selected.ID, true, false, dash.enableAI)
			}
		}
	}
}

func (cmd *Dash) resize() {
	// Compute inner content area (inside title/footer bars and frame padding/border)
	// Outer height/width come from the terminal
	outerW, outerH := cmd.width, cmd.height

	// Space taken by title bar + footer bar
	usableH := max(0, outerH-cmd.headerHeight-cmd.footerHeight)

	// Frame padding & borders
	// Rounded border adds 1px each side; padding is 1 each side per frameStyle.
	hPad := 2 /* left+right padding */ + 2 /* left+right border */
	vPad := 2 /* top+bottom padding */ + 2 /* top+bottom border */

	tableW := max(20, outerW-hPad) // ensure a minimum width
	tableH := max(3, usableH-vPad) // ensure a minimum height

	cmd.tableWidth = tableW
	cmd.tableHeight = tableH

	// Recompute responsive columns:
	// Keep ID ~10, Status ~14, Title takes the rest.
	idW := 10
	statusW := 14
	titleW := max(20, tableW-idW-statusW-4) // small fudge for internal padding

	cmd.table.SetColumns([]table.Column{
		{Title: "#", Width: idW},
		{Title: "Title", Width: titleW},
		{Title: "Status", Width: statusW},
	})

	// Page size == visible height (minus header row inside the table)
	cmd.table.SetHeight(tableH)
	cmd.table.SetWidth(tableW)
}
