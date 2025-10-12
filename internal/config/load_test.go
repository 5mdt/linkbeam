// load_test.go

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

func TestLoad_Validate(t *testing.T) {
	tmpDir := t.TempDir()
	badConfigPath := filepath.Join(tmpDir, "bad.yaml")

	// Write invalid config (empty name, invalid theme)
	content := []byte(`
name: ""
theme: "unknown"
`)
	if err := os.WriteFile(badConfigPath, content, 0644); err != nil {
		t.Fatalf("failed to write bad config: %v", err)
	}

	_, err := Load(badConfigPath)
	if err == nil {
		t.Error("expected error for invalid config, got nil")
	}
}
