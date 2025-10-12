// config.go

/*
 * Copyright (c) - All Rights Reserved.
 *
 * See the LICENSE file for more information.
 */

package config

import (
	"errors"
	"os"
	"strings"
)

type Config struct {
	Name           string        `yaml:"name"`
	Bio            string        `yaml:"bio"`
	Avatar         string        `yaml:"avatar"`
	Theme          string        `yaml:"theme"`
	FontAwesomeCDN string        `yaml:"font_awesome_cdn"`
	Footer         []FooterBlock `yaml:"footer"`
	Links          []Link        `yaml:"links"`
	Socials        []Social      `yaml:"socials"`
}

type Link struct {
	Title string `yaml:"title"`
	URL   string `yaml:"url"`
	Icon  string `yaml:"icon"`
}

type Social struct {
	Platform string `yaml:"platform"`
	URL      string `yaml:"url"`
	Icon     string `yaml:"icon"`
}

type FooterBlock struct {
	Text string `yaml:"text"`
}

func (c *Config) LinksCount() int {
	return len(c.Links)
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
