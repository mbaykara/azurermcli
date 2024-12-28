package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mbaykara/azurermcli/internal/azure"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Resource types available in the UI
var resourceTypes = []string{
	"Clusters",
	"Compute",
	"Network",
	"Storage",
	"All",
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.updateLayout(msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		// Handle search mode
		if m.searchMode {
			switch msg.Type {
			case tea.KeyEsc:
				m.searchMode = false
				m.searchQuery = ""
				m.updateTableWithResources()
				return m, nil
			case tea.KeyBackspace:
				if len(m.searchQuery) > 0 {
					m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
					m.updateTableWithResources()
				}
				return m, nil
			case tea.KeyEnter:
				m.searchMode = false
				if m.searchQuery != "" {
					m.updateTableWithResources()
				}
				return m, nil
			default:
				if msg.Type == tea.KeyRunes {
					m.searchQuery += string(msg.Runes)
					m.updateTableWithResources()
				}
				return m, nil
			}
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "/":
			if m.currentView == "resources" && m.selectedResourceType != "" {
				m.searchMode = true
				m.searchQuery = ""
				return m, nil
			}
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
					m.selectedResourceType = "All"
					m.loading = true
					return m, azure.FetchResources(m.selectedSub, m.selectedRG)
				}
			}
		case "right", "left":
			if m.currentView == "resources" && !m.searchMode {
				oldType := m.selectedResourceType
				if msg.String() == "right" {
					for i, rType := range resourceTypes {
						if rType == m.selectedResourceType {
							if i == len(resourceTypes)-1 {
								m.selectedResourceType = resourceTypes[0]
							} else {
								m.selectedResourceType = resourceTypes[i+1]
							}
							break
						}
					}
				} else {
					for i, rType := range resourceTypes {
						if rType == m.selectedResourceType {
							if i == 0 {
								m.selectedResourceType = resourceTypes[len(resourceTypes)-1]
							} else {
								m.selectedResourceType = resourceTypes[i-1]
							}
							break
						}
					}
				}
				if oldType != m.selectedResourceType {
					m.loading = true
					return m, azure.FetchResources(m.selectedSub, m.selectedRG)
				}
			}
		case "1", "2", "3", "4", "5":
			if m.currentView == "resources" && !m.searchMode {
				idx := int(msg.String()[0] - '1')
				if idx >= 0 && idx < len(resourceTypes) {
					oldType := m.selectedResourceType
					m.selectedResourceType = resourceTypes[idx]
					if oldType != m.selectedResourceType {
						m.loading = true
						return m, azure.FetchResources(m.selectedSub, m.selectedRG)
					}
				}
			}
		case "6":
			if m.currentView == "resources" && !m.searchMode {
				idx := int(msg.String()[0] - '1')
				if idx >= 0 && idx < len(resourceTypes) {
					oldType := m.selectedResourceType
					m.selectedResourceType = resourceTypes[idx]
					if oldType != m.selectedResourceType {
						m.loading = true
						return m, azure.FetchResources(m.selectedSub, m.selectedRG)
					}
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
		}

		// Handle table navigation
		if !m.searchMode {
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

		// Find first tab that has resources
		foundResources := false
		for _, tab := range tabs {
			if tab == "All" {
				continue // Skip "All" tab in initial search
			}
			for _, resource := range msg.Resources {
				if matchResourceType(*resource.Type, tab) {
					m.currentTab = tab
					foundResources = true
					break
				}
			}
			if foundResources {
				break
			}
		}

		// If no specific resource type found, default to "All"
		if !foundResources {
			m.currentTab = "All"
		}

		m.updateTableWithResources()
		return m, nil
	}

	return m, nil
}

var tabs = []string{"Clusters", "Compute", "Network", "Storage", "All"}

func (m *Model) updateTableWithSubscriptions() {
	// First clear the rows
	m.table.SetRows([]table.Row{})

	// Calculate responsive column widths
	nameWidth := int(float64(m.width) * 0.4)  // 40% of width
	idWidth := int(float64(m.width) * 0.4)    // 40% of width
	stateWidth := int(float64(m.width) * 0.2) // 20% of width

	// Update columns
	columns := []table.Column{
		{Title: "Name", Width: nameWidth},
		{Title: "ID", Width: idWidth},
		{Title: "State", Width: stateWidth},
	}
	m.table.SetColumns(columns)

	// Set rows
	var rows []table.Row
	for _, sub := range m.subscriptions {
		rows = append(rows, table.Row{
			*sub.DisplayName,
			*sub.SubscriptionID,
			string(*sub.State),
		})
	}
	m.table.SetRows(rows)
	if len(rows) > 0 {
		m.table.SetCursor(0)
	}
}

func (m *Model) updateTableWithResourceGroups() {
	// First clear the rows
	m.table.SetRows([]table.Row{})

	// Calculate responsive column widths
	nameWidth := int(float64(m.width) * 0.5)     // 50% of width
	locationWidth := int(float64(m.width) * 0.3) // 30% of width
	statusWidth := int(float64(m.width) * 0.2)   // 20% of width

	// Update columns
	columns := []table.Column{
		{Title: "Name", Width: nameWidth},
		{Title: "Location", Width: locationWidth},
		{Title: "Status", Width: statusWidth},
	}
	m.table.SetColumns(columns)

	// Set rows
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
	if len(rows) > 0 {
		m.table.SetCursor(0)
	}
}

func (m *Model) updateTableWithResources() {
	// First clear the rows
	m.table.SetRows([]table.Row{})

	// Calculate responsive column widths
	nameWidth := int(float64(m.width) * 0.4)   // 40% of width
	typeWidth := int(float64(m.width) * 0.4)   // 40% of width
	statusWidth := int(float64(m.width) * 0.2) // 20% of width

	// Update columns
	columns := []table.Column{
		{Title: "Name", Width: nameWidth},
		{Title: "Type", Width: typeWidth},
		{Title: "Status", Width: statusWidth},
	}
	m.table.SetColumns(columns)

	// Set rows
	var rows []table.Row
	if resources, ok := m.resources[m.selectedRG]; ok {
		for _, resource := range resources {
			matchesTab := m.selectedResourceType == "All" || matchResourceType(*resource.Type, m.selectedResourceType)
			matchesSearch := !m.searchMode || strings.Contains(strings.ToLower(*resource.Name), strings.ToLower(m.searchQuery))

			if matchesTab && matchesSearch {
				resourceType := *resource.Type
				if m.selectedResourceType != "All" {
					resourceType = formatResourceType(resourceType)
				}
				rows = append(rows, table.Row{
					*resource.Name,
					resourceType,
					getResourceStatus(resource),
				})
			}
		}
	}

	var message string
	if len(rows) == 0 {
		if m.searchMode {
			message = fmt.Sprintf("No matches for '%s'", m.searchQuery)
		} else {
			message = fmt.Sprintf("No %s found in this resource group", strings.ToLower(m.selectedResourceType))
		}
		rows = append(rows, table.Row{message, "-", "-"})
	}

	m.table.SetRows(rows)
	if len(rows) > 0 {
		m.table.SetCursor(0)
	}
}

func formatResourceType(resourceType string) string {
	// Remove the Microsoft.* prefix and convert to title case
	parts := strings.Split(resourceType, "/")
	if len(parts) >= 2 {
		lastPart := parts[len(parts)-1]
		caser := cases.Title(language.English)
		return caser.String(lastPart)
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
	case "All":
		return true // Show all resource types
	default:
		return false
	}
}
