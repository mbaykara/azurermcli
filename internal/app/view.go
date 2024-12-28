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

	// Show context information in resources view
	if m.currentView == "resources" {
		contextInfo := fmt.Sprintf("Subscription: %s | Resource Group: %s", m.selectedSub, m.selectedRG)
		sb.WriteString(styles.HeaderStyle.Render(contextInfo))
		sb.WriteString("\n\n")

		// Show resource type tabs
		var renderedTabs []string
		for i, rType := range resourceTypes {
			tabStyle := styles.InactiveTabStyle
			if rType == m.selectedResourceType {
				tabStyle = styles.ActiveTabStyle
			}
			// Add number prefix to tabs
			tabText := fmt.Sprintf("%d: %s", (i+1)%len(resourceTypes), rType)
			renderedTabs = append(renderedTabs, tabStyle.Render(tabText))
		}

		tabBar := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
		sb.WriteString(tabBar)
		sb.WriteString("\n")

		// Add a separator line under the tabs
		separator := strings.Repeat("─", 100)
		sb.WriteString(styles.InactiveTabStyle.Render(separator))
		sb.WriteString("\n\n")

		// Show search bar if in search mode
		if m.searchMode {
			searchPrompt := fmt.Sprintf("Search: %s█", m.searchQuery)
			sb.WriteString(styles.SearchStyle.Render(searchPrompt))
			sb.WriteString("\n\n")
		}
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
		if m.searchMode {
			footerText += " • enter: finish search • esc: cancel search"
		} else {
			footerText += " • ←/→ or 1-5: switch resource type • /: search • esc: back to resource groups"
		}
	}

	sb.WriteString(styles.FooterStyle.Render(footerText))

	return sb.String()
}

// Add methods for updating table content...
