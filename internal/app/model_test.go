package app

import (
	"testing"
)

func TestNewModel(t *testing.T) {
	model := New()

	// Test initial state
	if model.currentView != "subscriptions" {
		t.Errorf("Initial view = %q, want %q", model.currentView, "subscriptions")
	}

	if model.loading != true {
		t.Errorf("Initial loading = %v, want %v", model.loading, true)
	}

	if model.selectedResourceType != "" {
		t.Errorf("Initial selectedResourceType = %q, want %q", model.selectedResourceType, "")
	}

	if model.searchMode != false {
		t.Errorf("Initial searchMode = %v, want %v", model.searchMode, false)
	}

	// Test map initialization
	if model.resourceGroups == nil {
		t.Error("resourceGroups map not initialized")
	}

	if model.resources == nil {
		t.Error("resources map not initialized")
	}
}

func TestUpdateLayout(t *testing.T) {
	tests := []struct {
		name          string
		width         int
		height        int
		currentView   string
		searchMode    bool
		wantMinHeight int
	}{
		{
			name:          "Subscriptions view",
			width:         100,
			height:        30,
			currentView:   "subscriptions",
			searchMode:    false,
			wantMinHeight: 26, // height - 4 (header/footer)
		},
		{
			name:          "Resources view with search",
			width:         100,
			height:        30,
			currentView:   "resources",
			searchMode:    true,
			wantMinHeight: 20, // height - 8 (header/tabs/footer) - 2 (search)
		},
		{
			name:          "Small window",
			width:         50,
			height:        5,
			currentView:   "resources",
			searchMode:    false,
			wantMinHeight: 3, // minimum height
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := New()
			model.currentView = tt.currentView
			model.searchMode = tt.searchMode

			model.updateLayout(tt.width, tt.height)

			if model.width != tt.width {
				t.Errorf("width = %d, want %d", model.width, tt.width)
			}

			if model.height != tt.height {
				t.Errorf("height = %d, want %d", model.height, tt.height)
			}

			tableHeight := model.table.Height()
			if tableHeight < tt.wantMinHeight {
				t.Errorf("table height = %d, want >= %d", tableHeight, tt.wantMinHeight)
			}
		})
	}
}
