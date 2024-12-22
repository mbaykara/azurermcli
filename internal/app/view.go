package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mbaykara/azurermcli/internal/styles"
)

func (m Model) View() string {
	var sb strings.Builder

	// Header
	sb.WriteString(styles.HeaderStyle.Render(m.header))
	sb.WriteString("\n\n")

	// Only show tabs in resources view
	if m.currentView == "resources" {
		tabs := []string{
			"Clusters",
			"Compute",
			"Network",
			"Storage",
			"Web",
			"All",
		}

		var renderedTabs []string
		for i, tab := range tabs {
			tabStyle := styles.InactiveTabStyle
			if tab == m.currentTab {
				tabStyle = styles.ActiveTabStyle
			}
			// Add number prefix to tabs
			tabText := fmt.Sprintf("%d: %s", (i+1)%6, tab) // Use modulo to make "All" tab "0"
			renderedTabs = append(renderedTabs, tabStyle.Render(tabText))
		}

		tabBar := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
		sb.WriteString(styles.TabContainerStyle.Render(tabBar))
		sb.WriteString("\n\n")
	}

	// Content
	if m.loading {
		sb.WriteString(m.spinner.View())
		sb.WriteString(" Loading...")
	} else if m.err != nil {
		sb.WriteString(styles.ErrorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
	} else {
		sb.WriteString(m.table.View())
	}

	// Footer
	sb.WriteString("\n")
	footerText := "q: quit"
	switch m.currentView {
	case "subscriptions":
		footerText += " • enter: select subscription"
	case "resourcegroups":
		footerText += " • enter: view resources • esc: back to subscriptions"
	case "resources":
		footerText += " • tab/shift+tab or 1-5,0: switch view • esc: back to resource groups"
	}
	
	sb.WriteString(styles.FooterStyle.Render(footerText))

	return sb.String()
}

// Add methods for updating table content... 