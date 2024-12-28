package app

import (
	"testing"
)

func TestFormatResourceType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "AKS cluster",
			input:    "Microsoft.ContainerService/managedClusters",
			expected: "ManagedClusters",
		},
		{
			name:     "Virtual network",
			input:    "Microsoft.Network/virtualNetworks",
			expected: "VirtualNetworks",
		},
		{
			name:     "No slash",
			input:    "simpleResource",
			expected: "simpleResource",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatResourceType(tt.input)
			if result != tt.expected {
				t.Errorf("formatResourceType(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMatchResourceType(t *testing.T) {
	tests := []struct {
		name         string
		resourceType string
		tab          string
		expected     bool
	}{
		{
			name:         "AKS cluster in Clusters tab",
			resourceType: "Microsoft.ContainerService/managedClusters",
			tab:          "Clusters",
			expected:     true,
		},
		{
			name:         "VM in Compute tab",
			resourceType: "Microsoft.Compute/virtualMachines",
			tab:          "Compute",
			expected:     true,
		},
		{
			name:         "Storage account in Storage tab",
			resourceType: "Microsoft.Storage/storageAccounts",
			tab:          "Storage",
			expected:     true,
		},
		{
			name:         "Network in Network tab",
			resourceType: "Microsoft.Network/virtualNetworks",
			tab:          "Network",
			expected:     true,
		},
		{
			name:         "VM not in Network tab",
			resourceType: "Microsoft.Compute/virtualMachines",
			tab:          "Network",
			expected:     false,
		},
		{
			name:         "Any resource in All tab",
			resourceType: "Microsoft.AnyService/anyResource",
			tab:          "All",
			expected:     true,
		},
		{
			name:         "Resource in unknown tab",
			resourceType: "Microsoft.AnyService/anyResource",
			tab:          "Unknown",
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchResourceType(tt.resourceType, tt.tab)
			if result != tt.expected {
				t.Errorf("matchResourceType(%q, %q) = %v, want %v", tt.resourceType, tt.tab, result, tt.expected)
			}
		})
	}
}

func TestGetResourceStatus(t *testing.T) {
	t.Run("Default status", func(t *testing.T) {
		result := getResourceStatus(nil)
		expected := "Running"
		if result != expected {
			t.Errorf("getResourceStatus(nil) = %q, want %q", result, expected)
		}
	})
}
