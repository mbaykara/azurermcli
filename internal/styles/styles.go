package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	BaseStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240"))

	HeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		Background(lipgloss.Color("236")).
		Width(100).
		Align(lipgloss.Center).
		Padding(0, 1)

	FooterStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Background(lipgloss.Color("236")).
		Width(100).
		Align(lipgloss.Left).
		Padding(0, 1)

	TitleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Bold(true)

	ErrorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true)

	TabContainerStyle = lipgloss.NewStyle().
		Width(100).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1)

	ActiveTabStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Background(lipgloss.Color("236")).
		Bold(true).
		Padding(0, 2)

	InactiveTabStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("250")).
		Background(lipgloss.Color("236")).
		Padding(0, 2)

	ResourceStyles = map[string]lipgloss.Style{
		"AKS Cluster": lipgloss.NewStyle().Foreground(lipgloss.Color("87")),
		"Microsoft.ContainerService/managedClusters": lipgloss.NewStyle().Foreground(lipgloss.Color("87")),
		"Microsoft.Network/virtualNetworks": lipgloss.NewStyle().Foreground(lipgloss.Color("39")),
		"Microsoft.Storage/storageAccounts": lipgloss.NewStyle().Foreground(lipgloss.Color("220")),
		"Microsoft.Compute/virtualMachines": lipgloss.NewStyle().Foreground(lipgloss.Color("82")),
		"Microsoft.Web/sites": lipgloss.NewStyle().Foreground(lipgloss.Color("213")),
		"default": lipgloss.NewStyle().Foreground(lipgloss.Color("250")),
	}
) 