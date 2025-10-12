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

	tests := []struct {
		name string
		got  interface{}
		want interface{}
	}{
		{"name", cfg.Name, "Ada Lovelace"},
		{"links count", len(cfg.Links), 2},
		{"first link URL", cfg.Links[0].URL, "https://ada.blog"},
	}

	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s: got %v, want %v", tt.name, tt.got, tt.want)
		}
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
	if got, want := cfg.LinksCount(), 2; got != want {
		t.Errorf("LinksCount() = %d, want %d", got, want)
	}
}

func TestFooterBlocks(t *testing.T) {
	tests := []struct {
		name   string
		cfg    Config
		want   int
	}{
		{
			name: "with footer blocks",
			cfg: Config{
				Name: "Test",
				Footer: []FooterBlock{
					{Text: "© 2025 Test"},
					{Text: "Made with LinkBeam"},
				},
			},
			want: 2,
		},
		{
			name:   "empty footer",
			cfg:    Config{Name: "Test", Footer: []FooterBlock{}},
			want:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := len(tt.cfg.Footer); got != tt.want {
				t.Errorf("Footer length = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestFooterBlockText(t *testing.T) {
	block := FooterBlock{Text: "Test Footer"}
	want := "Test Footer"

	if block.Text != want {
		t.Errorf("FooterBlock.Text = %q, want %q", block.Text, want)
	}
}
