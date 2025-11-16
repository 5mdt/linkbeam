// config_test.go

/*
 * Copyright (c) - All Rights Reserved.
 *
 * See the LICENSE file for more information.
 */

package config

import (
	"testing"
)

func loadTestConfig(t *testing.T) *Config {
	t.Helper()
	cfg, err := Load("testdata/config.yaml")
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	return cfg
}

func TestLoad(t *testing.T) {
	cfg := loadTestConfig(t)

	if cfg.Name != "Ada Lovelace" {
		t.Errorf("name: got %v, want %v", cfg.Name, "Ada Lovelace")
	}

	if len(cfg.Content) == 0 {
		t.Fatal("expected content blocks, got none")
	}

	// Find the links collection
	var linksFound bool
	for _, block := range cfg.Content {
		if block.Type == "vertical-list-text" {
			for collectionName, items := range block.Collections {
				if collectionName == "links" {
					linksFound = true
					if len(items) != 4 {
						t.Errorf("links count: got %d, want 4", len(items))
					}
					if len(items) > 0 && items[0].URL != "https://ada.blog/research" {
						t.Errorf("first link URL: got %v, want https://ada.blog/research", items[0].URL)
					}
				}
			}
		}
	}

	if !linksFound {
		t.Error("expected to find links collection, but didn't")
	}
}

func TestThemeValidation(t *testing.T) {
	cfg := loadTestConfig(t)
	if cfg.Theme != "light" && cfg.Theme != "dark" && cfg.Theme != "auto" {
		t.Errorf("invalid theme: %s", cfg.Theme)
	}
}

func TestLinksCountMethod(t *testing.T) {
	cfg := loadTestConfig(t)
	// Config has 4 links in "links" collection + 1 in "gaming" + 5 in "socials" = 10 URLs total
	if got, want := cfg.LinksCount(), 10; got != want {
		t.Errorf("LinksCount() = %d, want %d", got, want)
	}
}

func TestContentBlocks(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
		want int
	}{
		{
			name: "with content blocks",
			cfg: Config{
				Name: "Test",
				Content: []ContentBlock{
					{
						Type: "footer",
						Collections: map[string][]Item{
							"footer": {
								{Text: "© 2025 Test"},
								{Text: "Made with LinkBeam"},
							},
						},
					},
				},
			},
			want: 1,
		},
		{
			name: "empty content",
			cfg:  Config{Name: "Test", Content: []ContentBlock{}},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := len(tt.cfg.Content); got != tt.want {
				t.Errorf("Content length = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestItemTypes(t *testing.T) {
	tests := []struct {
		name string
		item Item
		desc string
	}{
		{
			name: "URL item",
			item: Item{Title: "Test Link", URL: "https://example.com", Icon: "fas fa-link"},
			desc: "link item",
		},
		{
			name: "CopyText item",
			item: Item{Title: "Username", CopyText: "testuser", Icon: "fas fa-user"},
			desc: "copyable text item",
		},
		{
			name: "Text item",
			item: Item{Text: "© 2025"},
			desc: "plain text item",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.item.Title == "" && tt.item.Text == "" {
				t.Error("item should have title or text")
			}
			hasContent := tt.item.URL != "" || tt.item.CopyText != "" || tt.item.Text != ""
			if !hasContent {
				t.Error("item should have at least one content field")
			}
		})
	}
}
