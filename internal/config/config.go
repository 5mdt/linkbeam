// config.go

/*
 * Copyright (c) - All Rights Reserved.
 *
 * See the LICENSE file for more information.
 */

package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Name    string         `yaml:"name"`
	Bio     string         `yaml:"bio"`
	Avatar  string         `yaml:"avatar"`
	Theme   string         `yaml:"theme"`
	Content []ContentBlock `yaml:"content"`
}

// ContentBlock represents a container section with a specific display type
// that can hold any number of named collections of items
type ContentBlock struct {
	Type string `yaml:"type"`
	// Collections holds dynamic named groups of items (e.g., "links", "gaming", "socials")
	// We need to manually unmarshal this to handle arbitrary keys
	Collections map[string][]Item `yaml:",inline"`
}

// Item represents a universal content item that can be a link, copyable text, or plain text
type Item struct {
	Title    string `yaml:"title,omitempty"`
	URL      string `yaml:"url,omitempty"`      // Clickable link
	CopyText string `yaml:"copy-text,omitempty"` // Click to copy text
	Text     string `yaml:"text,omitempty"`     // Plain display text
	Icon     string `yaml:"icon,omitempty"`
	Image    string `yaml:"image,omitempty"`    // Image URL or path
	AltText  string `yaml:"alt-text,omitempty"` // Accessibility text for images
}

func (c *Config) LinksCount() int {
	count := 0
	for _, block := range c.Content {
		for _, items := range block.Collections {
			for _, item := range items {
				if item.URL != "" {
					count++
				}
			}
		}
	}
	return count
}

// GetAvailableThemes discovers available themes from the themes directory.
// It returns a list of theme names (without the .css extension).
func GetAvailableThemes(themesDir string) ([]string, error) {
	entries, err := os.ReadDir(themesDir)
	if err != nil {
		// If themes directory doesn't exist, return empty list
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var themes []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".css") {
			themeName := strings.TrimSuffix(name, ".css")
			themes = append(themes, themeName)
		}
	}
	return themes, nil
}

func (c *Config) Validate() error {
	return c.ValidateWithThemes("")
}

// ValidateWithThemes validates the config with theme discovery.
// If themesDir is empty, theme validation is skipped.
func (c *Config) ValidateWithThemes(themesDir string) error {
	if c.Name == "" {
		return errors.New("name cannot be empty")
	}

	// Validate content blocks
	validBlockTypes := map[string]bool{
		"vertical-list-text":    true,
		"horizontal-list-icons": true,
		"vertical-list-images":  true,
		"footer":                true,
	}

	for i, block := range c.Content {
		if !validBlockTypes[block.Type] {
			return fmt.Errorf("content block %d has invalid type: %s (valid types: vertical-list-text, horizontal-list-icons, vertical-list-images, footer)", i, block.Type)
		}

		// Validate items within collections
		for collectionName, items := range block.Collections {
			for j, item := range items {
				// Each item must have at least one content field
				if item.URL == "" && item.CopyText == "" && item.Text == "" && item.Image == "" {
					return fmt.Errorf("content block %d, collection '%s', item %d must have at least one of: url, copy-text, text, or image", i, collectionName, j)
				}
			}
		}
	}

	// Skip theme validation if no themes directory provided
	if themesDir == "" {
		return nil
	}

	availableThemes, err := GetAvailableThemes(themesDir)
	if err != nil {
		return err
	}

	// If no themes found, skip validation
	if len(availableThemes) == 0 {
		return nil
	}

	// Check if the configured theme is available
	validThemes := make(map[string]bool)
	for _, theme := range availableThemes {
		validThemes[theme] = true
	}

	if !validThemes[c.Theme] {
		return errors.New("theme '" + c.Theme + "' not found. Available themes: " + strings.Join(availableThemes, ", "))
	}

	return nil
}
