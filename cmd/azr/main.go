package main

import (
    "fmt"
    "os"

    "github.com/mbaykara/azurermcli/internal/app"    // Replace 'yourusername' with your GitHub username
    "github.com/mbaykara/azurermcli/internal/azure"  // Replace 'yourusername' with your GitHub username
    tea "github.com/charmbracelet/bubbletea"
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