// validate_test.go

/*
 * Copyright (c) - All Rights Reserved.
 *
 * See the LICENCE file for more information.
 */

package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name      string
		cfg       Config
		wantError bool
	}{
		{"missing name", Config{Name: "", Theme: "auto"}, true},
		{"valid config", Config{Name: "Ada", Theme: "auto"}, false},
		{"valid config with any theme", Config{Name: "Ada", Theme: "invalid"}, false}, // No theme validation without themesDir
		{
			name: "valid content blocks",
			cfg: Config{
				Name:  "Test",
				Theme: "auto",
				Content: []ContentBlock{
					{
						Type: "vertical-list-text",
						Collections: map[string][]Item{
							"links": {{Title: "Test", URL: "https://test.com"}},
						},
					},
				},
			},
			wantError: false,
		},
		{
			name: "invalid block type",
			cfg: Config{
				Name:  "Test",
				Theme: "auto",
				Content: []ContentBlock{
					{Type: "invalid-type"},
				},
			},
			wantError: true,
		},
		{
			name: "item without content",
			cfg: Config{
				Name:  "Test",
				Theme: "auto",
				Content: []ContentBlock{
					{
						Type: "vertical-list-text",
						Collections: map[string][]Item{
							"links": {{Title: "Empty"}}, // No url, copy-text, or text
						},
					},
				},
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantError {
				t.Errorf("Validate() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateWithThemes(t *testing.T) {
	// Create a temporary themes directory for testing
	tmpDir := t.TempDir()
	themesDir := filepath.Join(tmpDir, "themes")
	if err := os.MkdirAll(themesDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create some test theme files
	for _, theme := range []string{"auto", "nord", "gruvbox", "catppuccin"} {
		if err := os.WriteFile(filepath.Join(themesDir, theme+".css"), []byte("/* test */"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name      string
		cfg       Config
		themesDir string
		wantError bool
	}{
		{"missing name", Config{Name: "", Theme: "auto"}, themesDir, true},
		{"invalid theme", Config{Name: "Ada", Theme: "invalid"}, themesDir, true},
		{"valid auto theme", Config{Name: "Ada", Theme: "auto"}, themesDir, false},
		{"valid nord theme", Config{Name: "Ada", Theme: "nord"}, themesDir, false},
		{"valid gruvbox theme", Config{Name: "Ada", Theme: "gruvbox"}, themesDir, false},
		{"valid catppuccin theme", Config{Name: "Ada", Theme: "catppuccin"}, themesDir, false},
		{"no themes dir", Config{Name: "Ada", Theme: "anything"}, "", false}, // Skip validation if no dir
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.ValidateWithThemes(tt.themesDir)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateWithThemes() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestGetAvailableThemes(t *testing.T) {
	// Create a temporary themes directory for testing
	tmpDir := t.TempDir()
	themesDir := filepath.Join(tmpDir, "themes")
	if err := os.MkdirAll(themesDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create some test theme files
	testThemes := []string{"auto", "nord", "gruvbox"}
	for _, theme := range testThemes {
		if err := os.WriteFile(filepath.Join(themesDir, theme+".css"), []byte("/* test */"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Create a non-CSS file (should be ignored)
	if err := os.WriteFile(filepath.Join(themesDir, "readme.txt"), []byte("readme"), 0644); err != nil {
		t.Fatal(err)
	}

	themes, err := GetAvailableThemes(themesDir)
	if err != nil {
		t.Fatalf("GetAvailableThemes() error = %v", err)
	}

	if len(themes) != len(testThemes) {
		t.Errorf("GetAvailableThemes() got %d themes, want %d", len(themes), len(testThemes))
	}

	// Check that all expected themes are present
	themeMap := make(map[string]bool)
	for _, theme := range themes {
		themeMap[theme] = true
	}

	for _, expectedTheme := range testThemes {
		if !themeMap[expectedTheme] {
			t.Errorf("Expected theme %q not found in available themes", expectedTheme)
		}
	}

	// Test non-existent directory
	themes, err = GetAvailableThemes(filepath.Join(tmpDir, "nonexistent"))
	if err != nil {
		t.Errorf("GetAvailableThemes() should not error on non-existent dir, got %v", err)
	}
	if len(themes) != 0 {
		t.Errorf("GetAvailableThemes() should return empty list for non-existent dir, got %d themes", len(themes))
	}
}
