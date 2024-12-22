package app

import (
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mbaykara/azurermcli/internal/azure"
)

type Model struct {
	table           table.Model
	spinner         spinner.Model
	loading         bool
	width          int
	height         int
	header         string
	footer         string
	subscriptions  []armsubscription.Subscription
	resourceGroups map[string][]armresources.ResourceGroup
	resources      map[string][]armresources.GenericResourceExpanded
	currentView    string
	currentTab     string
	selectedSub    string
	selectedRG     string
	err           error
}

func New() Model {
	return Model{
		table:           initTable(),
		spinner:         initSpinner(),
		loading:        true,
		resourceGroups: make(map[string][]armresources.ResourceGroup),
		resources:      make(map[string][]armresources.GenericResourceExpanded),
		currentView:    "subscriptions",
		currentTab:     "subscriptions",
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, azure.FetchSubscriptions)
}

func initTable() table.Model {
	columns := []table.Column{
		{Title: "Name", Width: 40},
		{Title: "ID", Width: 40},
		{Title: "State", Width: 20},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)
	return t
}

func initSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return s
}

func (m *Model) updateHeaderFooter() {
	m.header = "Azure Resource Manager CLI"
	m.footer = "Press q to quit • Press h for help"
}

func (m *Model) updateLayout(width, height int) {
	m.width = width
	m.height = height
	m.table.SetWidth(width)
	m.table.SetHeight(height - 4) // Account for header and footer
} 