package app

import (
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mbaykara/azurermcli/internal/azure"
	"github.com/mbaykara/azurermcli/internal/styles"
)

type Model struct {
	table                table.Model
	spinner              spinner.Model
	loading              bool
	width                int
	height               int
	header               string
	subscriptions        []armsubscription.Subscription
	resourceGroups       map[string][]armresources.ResourceGroup
	resources            map[string][]armresources.GenericResourceExpanded
	currentView          string
	currentTab           string
	selectedSub          string
	selectedRG           string
	selectedResourceType string
	err                  error
	showTabs             bool
	searchMode           bool
	searchQuery          string
}

func New() Model {
	return Model{
		table:                initTable(),
		spinner:              initSpinner(),
		loading:              true,
		resourceGroups:       make(map[string][]armresources.ResourceGroup),
		resources:            make(map[string][]armresources.GenericResourceExpanded),
		currentView:          "subscriptions",
		currentTab:           "All",
		selectedResourceType: "",
		showTabs:             false,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, azure.FetchSubscriptions)
}

func initTable() table.Model {
	columns := []table.Column{
		{Title: "Name", Width: 40},
		{Title: "Description", Width: 60},
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

func (m *Model) updateLayout(width, height int) {
	m.width = width
	m.height = height

	// Adjust header width to terminal width
	styles.HeaderStyle = styles.HeaderStyle.Width(width)
	styles.FooterStyle = styles.FooterStyle.Width(width)
	styles.TabContainerStyle = styles.TabContainerStyle.Width(width)

	// Calculate table height: total height minus space for header, context, tabs, footer
	tableHeight := height
	if m.currentView == "resources" {
		tableHeight -= 8 // Subtract space for header, context info, tabs, footer, and spacing
		if m.searchMode {
			tableHeight -= 2 // Additional space for search bar
		}
	} else {
		tableHeight -= 4 // Just header and footer for other views
	}
	if tableHeight < 3 {
		tableHeight = 3 // Minimum height for table
	}

	// Update table dimensions
	m.table.SetWidth(width)
	m.table.SetHeight(tableHeight)

	// Adjust column widths based on terminal width
	switch m.currentView {
	case "subscriptions":
		m.updateTableWithSubscriptions()
	case "resourcegroups":
		m.updateTableWithResourceGroups()
	case "resources":
		m.updateTableWithResources()
	}
}

// Resource type menu options
var resourceTypes = []string{
	"Clusters",
	"Compute",
	"Network",
	"Storage",
	"All",
}

func (m *Model) updateResourceTypeMenu() {
	m.table.SetRows([]table.Row{})

	columns := []table.Column{
		{Title: "Resource Type", Width: 30},
		{Title: "Description", Width: 50},
	}
	m.table.SetColumns(columns)

	var rows []table.Row
	for _, rType := range resourceTypes {
		description := getResourceTypeDescription(rType)
		rows = append(rows, table.Row{rType, description})
	}
	m.table.SetRows(rows)

	if len(rows) > 0 {
		m.table.SetCursor(0)
	}
}

func getResourceTypeDescription(resourceType string) string {
	switch resourceType {
	case "Clusters":
		return "Container Services and Managed Clusters"
	case "Compute":
		return "Virtual Machines, VM Scale Sets, and Disks"
	case "Network":
		return "Virtual Networks, NSGs, Load Balancers, and Gateways"
	case "Storage":
		return "Storage Accounts, File Services, and Blob Storage"
	case "All":
		return "All Azure Resources"
	default:
		return ""
	}
}
