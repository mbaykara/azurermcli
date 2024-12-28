# azurermcli

A terminal UI for managing Azure resources, inspired by k9s.

## Features

- Navigate Azure resources with an intuitive terminal interface
- Filter resources by type (Clusters, Compute, Network, Storage)
- Real-time search within resources
- Responsive design that adapts to terminal size

## Installation

### From Releases
Download the latest binary from [GitHub Releases](https://github.com/mbaykara/azurermcli/releases) for your platform:

```bash
# Linux (amd64)
curl -L -o azr https://github.com/mbaykara/azurermcli/releases/latest/download/azr-linux-amd64
chmod +x azr
sudo mv azr /usr/local/bin/

# macOS (Apple Silicon)
curl -L -o azr https://github.com/mbaykara/azurermcli/releases/latest/download/azr-darwin-arm64
chmod +x azr
sudo mv azr /usr/local/bin/
```

### From Source
```bash
go install github.com/mbaykara/azurermcli/cmd/azr@latest
```

## Usage

1. Ensure you're logged in to Azure CLI:
```bash
az login
```

2. Run azurermcli:
```bash
azr
```

### Navigation

- Use arrow keys to navigate
- Enter to select
- ESC to go back
- 1-5 or ←/→ to switch resource types
- / to search within current view
- q to quit
