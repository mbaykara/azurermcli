package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mbaykara/azurermcli/internal/app"
	"github.com/mbaykara/azurermcli/internal/azure"
)

func main() {
	p := tea.NewProgram(app.New(), tea.WithAltScreen())

	// Start by fetching subscriptions
	go func() {
		msg := azure.FetchSubscriptions()
		p.Send(msg)
	}()

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
