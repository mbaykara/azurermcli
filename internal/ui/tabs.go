package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/mbaykara/azurermcli/internal/styles"
)

var ResourceTabs = []string{
	"Clusters",
	"Compute",
	"Network",
	"Storage",
	"Web",
	"All",
}

func RenderTabs(currentTab string) string {
	var renderedTabs []string

	for _, tab := range ResourceTabs {
		if tab == currentTab {
			renderedTabs = append(renderedTabs, styles.ActiveTabStyle.Render(tab))
		} else {
			renderedTabs = append(renderedTabs, styles.InactiveTabStyle.Render(tab))
		}
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		renderedTabs...,
	)
} 