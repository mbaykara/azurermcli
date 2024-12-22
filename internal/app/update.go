package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mbaykara/azurermcli/internal/azure"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.updateLayout(msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			switch m.currentView {
			case "subscriptions":
				selected := m.table.SelectedRow()
				if len(selected) >= 2 {
					m.selectedSub = selected[1]
					m.currentView = "resourcegroups"
					m.loading = true
					return m, azure.FetchResourceGroups(m.selectedSub)
				}
			case "resourcegroups":
				selected := m.table.SelectedRow()
				if len(selected) >= 1 {
					m.selectedRG = selected[0]
					m.currentView = "resources"
					m.currentTab = "Clusters" // Default tab
					m.loading = true
					return m, azure.FetchResources(m.selectedSub, m.selectedRG)
				}
			}
		case "esc":
			switch m.currentView {
			case "resourcegroups":
				m.currentView = "subscriptions"
				m.updateTableWithSubscriptions()
			case "resources":
				m.currentView = "resourcegroups"
				m.updateTableWithResourceGroups()
			}
		case "tab":
			if m.currentView == "resources" {
				m.currentTab = nextTab(m.currentTab)
				m.updateTableWithResources()
			}
			return m, nil
		case "shift+tab":
			if m.currentView == "resources" {
				m.currentTab = prevTab(m.currentTab)
				m.updateTableWithResources()
			}
			return m, nil
		case "1", "2", "3", "4", "5", "0":
			if m.currentView == "resources" {
				switch msg.String() {
				case "1":
					m.currentTab = "Clusters"
				case "2":
					m.currentTab = "Compute"
				case "3":
					m.currentTab = "Network"
				case "4":
					m.currentTab = "Storage"
				case "5":
					m.currentTab = "Web"
				case "0":
					m.currentTab = "All"
				}
				m.updateTableWithResources()
				return m, nil
			}
		}

		// Handle table navigation
		if m.currentView != "resources" || (msg.String() != "tab" && msg.String() != "shift+tab") {
			m.table, cmd = m.table.Update(msg)
			return m, cmd
		}

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case azure.SubscriptionsMsg:
		m.loading = false
		m.subscriptions = msg.Subs
		m.updateTableWithSubscriptions()
		return m, nil

	case azure.ResourceGroupsMsg:
		m.loading = false
		m.resourceGroups[msg.SubscriptionID] = msg.Groups
		m.updateTableWithResourceGroups()
		return m, nil

	case azure.ResourcesMsg:
		m.loading = false
		m.resources[m.selectedRG] = msg.Resources
		m.updateTableWithResources()
		return m, nil
	}

	return m, nil
}

var tabs = []string{"Clusters", "Compute", "Network", "Storage", "Web", "All"}

func nextTab(current string) string {
	for i, tab := range tabs {
		if tab == current {
			if i == len(tabs)-1 {
				return tabs[0]
			}
			return tabs[i+1]
		}
	}
	return tabs[0]
}

func prevTab(current string) string {
	for i, tab := range tabs {
		if tab == current {
			if i == 0 {
				return tabs[len(tabs)-1]
			}
			return tabs[i-1]
		}
	}
	return tabs[len(tabs)-1]
}

func (m *Model) updateTableWithSubscriptions() {
	columns := []table.Column{
		{Title: "Name", Width: 40},
		{Title: "ID", Width: 40},
		{Title: "State", Width: 20},
	}
	m.table.SetColumns(columns)

	var rows []table.Row
	for _, sub := range m.subscriptions {
		rows = append(rows, table.Row{
			*sub.DisplayName,
			*sub.SubscriptionID,
			string(*sub.State),
		})
	}
	m.table.SetRows(rows)
}

func (m *Model) updateTableWithResourceGroups() {
	columns := []table.Column{
		{Title: "Name", Width: 30},
		{Title: "Location", Width: 20},
		{Title: "Status", Width: 20},
	}
	m.table.SetColumns(columns)

	var rows []table.Row
	if groups, ok := m.resourceGroups[m.selectedSub]; ok {
		for _, group := range groups {
			rows = append(rows, table.Row{
				*group.Name,
				*group.Location,
				"Available",
			})
		}
	}
	m.table.SetRows(rows)
}

func (m *Model) updateTableWithResources() {
	// Set up table columns specific to resource view
	columns := []table.Column{
		{Title: "Name", Width: 30},
		{Title: "Type", Width: 30},
		{Title: "Status", Width: 20},
	}
	m.table.SetColumns(columns)

	var rows []table.Row
	if resources, ok := m.resources[m.selectedRG]; ok {
		for _, resource := range resources {
			if m.currentTab == "All" || matchResourceType(*resource.Type, m.currentTab) {
				resourceType := formatResourceType(*resource.Type)
				rows = append(rows, table.Row{
					*resource.Name,
					resourceType,
					getResourceStatus(resource),
				})
			}
		}
	}
	
	// If no rows match the current tab, show an empty table with headers
	if len(rows) == 0 {
		rows = append(rows, table.Row{
			fmt.Sprintf("No %s found in this resource group", strings.ToLower(m.currentTab)),
			"-",
			"-",
		})
	}
	
	m.table.SetRows(rows)
}

func formatResourceType(resourceType string) string {
	// Remove the Microsoft.* prefix and convert to title case
	parts := strings.Split(resourceType, "/")
	if len(parts) >= 2 {
		return strings.Title(parts[len(parts)-1])
	}
	return resourceType
}

func getResourceStatus(resource interface{}) string {
	// You can add logic here to extract and return the status
	// of different resource types
	return "Running" // Default status
}

func matchResourceType(resourceType, tab string) bool {
	resourceType = strings.ToLower(resourceType)
	switch tab {
	case "Clusters":
		return strings.Contains(resourceType, "microsoft.containerservice/managedclusters") ||
			   strings.Contains(resourceType, "microsoft.container/containergroups")
	case "Compute":
		return strings.Contains(resourceType, "microsoft.compute/virtualmachines") ||
			   strings.Contains(resourceType, "microsoft.compute/vmscalesets") ||
			   strings.Contains(resourceType, "microsoft.compute/disks")
	case "Network":
		return strings.Contains(resourceType, "microsoft.network/virtualnetworks") ||
			   strings.Contains(resourceType, "microsoft.network/networksecuritygroups") ||
			   strings.Contains(resourceType, "microsoft.network/publicipaddresses") ||
			   strings.Contains(resourceType, "microsoft.network/loadbalancers") ||
			   strings.Contains(resourceType, "microsoft.network/applicationgateways")
	case "Storage":
		return strings.Contains(resourceType, "microsoft.storage/storageaccounts") ||
			   strings.Contains(resourceType, "microsoft.storage/fileservices") ||
			   strings.Contains(resourceType, "microsoft.storage/blobservices")
	case "Web":
		return strings.Contains(resourceType, "microsoft.web/sites") ||
			   strings.Contains(resourceType, "microsoft.web/serverfarms") ||
			   strings.Contains(resourceType, "microsoft.web/staticsites")
	default:
		return false
	}
}