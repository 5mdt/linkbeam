// generator_test.go

/*
 * Copyright (c) - All Rights Reserved.
 *
 * See the LICENCE file for more information.
 */

package generator

import (
	"os"
	"strings"
	"testing"

	"linkbeam/internal/config"
)

func TestGenerateSite(t *testing.T) {
	cfg, err := config.Load("../config/testdata/config.yaml")
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	tmpDir := t.TempDir()
	outPath := tmpDir + "/index.html"

	templatePath := "../../templates/base.html"
	if err = GenerateSite(cfg, templatePath, outPath); err != nil {
		t.Fatalf("failed to generate site: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, cfg.Name) {
		t.Errorf("output does not contain name %q", cfg.Name)
	}
	if !strings.Contains(content, cfg.Bio) {
		t.Errorf("output does not contain bio %q", cfg.Bio)
	}
}

func TestGenerateSiteWithFooter(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := tmpDir + "/index.html"

	footerText := "© 2025 Test User"
	cfg := &config.Config{
		Name:  "Test User",
		Bio:   "Test bio",
		Theme: "auto",
		Footer: []config.FooterBlock{
			{Text: footerText},
			{Text: "Made with LinkBeam"},
		},
	}

	templatePath := "../../templates/base.html"
	if err := GenerateSite(cfg, templatePath, outPath); err != nil {
		t.Fatalf("failed to generate site: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, footerText) {
		t.Errorf("output does not contain footer text %q", footerText)
	}
	if !strings.Contains(content, "Made with LinkBeam") {
		t.Errorf("output does not contain footer text %q", "Made with LinkBeam")
	}
}

func TestRenderUserPage_EmptyConfig(t *testing.T) {
	cfg := &config.Config{}
	output := RenderUserPage(cfg)

	for _, want := range []string{"Name:", "Bio:", "Links:"} {
		if !strings.Contains(output, want) {
			t.Errorf("output missing %q", want)
		}
	}
}

func TestCopyAssets(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a temporary themes directory with test files
	tmpThemesDir := tmpDir + "/themes_src"
	if err := os.MkdirAll(tmpThemesDir, 0755); err != nil {
		t.Fatalf("failed to create temp themes dir: %v", err)
	}

	// Create test CSS files
	testCSS := "body { color: red; }"
	if err := os.WriteFile(tmpThemesDir+"/test.css", []byte(testCSS), 0644); err != nil {
		t.Fatalf("failed to write test CSS: %v", err)
	}

	// Create output directory
	distDir := tmpDir + "/dist"

	// Copy assets
	if err := CopyAssets(distDir, tmpThemesDir); err != nil {
		t.Fatalf("CopyAssets failed: %v", err)
	}

	// Verify the file was copied
	copiedFile := distDir + "/themes/test.css"
	data, err := os.ReadFile(copiedFile)
	if err != nil {
		t.Fatalf("failed to read copied file: %v", err)
	}

	if string(data) != testCSS {
		t.Errorf("copied file content mismatch: got %q, want %q", string(data), testCSS)
	}
}

func TestCopyStaticFiles(t *testing.T) {
	tests := []struct {
		name        string
		setupStatic bool
		files       []string
		subdirs     []string
		wantErr     bool
		description string
	}{
		{
			name:        "copy single file",
			setupStatic: true,
			files:       []string{"logo.png"},
			subdirs:     []string{},
			wantErr:     false,
			description: "should copy a single static file",
		},
		{
			name:        "copy multiple files",
			setupStatic: true,
			files:       []string{"logo.png", "favicon.ico", "style.css"},
			subdirs:     []string{},
			wantErr:     false,
			description: "should copy multiple static files",
		},
		{
			name:        "copy files with subdirectories",
			setupStatic: true,
			files:       []string{"logo.png", "images/banner.jpg", "css/custom.css"},
			subdirs:     []string{"images", "css"},
			wantErr:     false,
			description: "should copy files preserving subdirectory structure",
		},
		{
			name:        "no static directory",
			setupStatic: false,
			files:       []string{},
			subdirs:     []string{},
			wantErr:     false,
			description: "should not fail when static directory doesn't exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			// Change to temp directory to simulate project root
			oldWd, err := os.Getwd()
			if err != nil {
				t.Fatalf("failed to get working directory: %v", err)
			}
			defer func() {
				if err := os.Chdir(oldWd); err != nil {
					t.Errorf("failed to restore working directory: %v", err)
				}
			}()

			if err := os.Chdir(tmpDir); err != nil {
				t.Fatalf("failed to change to temp directory: %v", err)
			}

			// Setup static directory and files if needed
			if tt.setupStatic {
				staticDir := tmpDir + "/static"
				if err := os.MkdirAll(staticDir, 0755); err != nil {
					t.Fatalf("failed to create static dir: %v", err)
				}

				// Create subdirectories
				for _, subdir := range tt.subdirs {
					subdirPath := staticDir + "/" + subdir
					if err := os.MkdirAll(subdirPath, 0755); err != nil {
						t.Fatalf("failed to create subdir %s: %v", subdir, err)
					}
				}

				// Create test files
				for _, file := range tt.files {
					filePath := staticDir + "/" + file
					testContent := []byte("test content for " + file)
					if err := os.WriteFile(filePath, testContent, 0644); err != nil {
						t.Fatalf("failed to write test file %s: %v", file, err)
					}
				}
			}

			distDir := tmpDir + "/dist"

			// Execute CopyStaticFiles
			err = CopyStaticFiles(distDir)

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("%s: CopyStaticFiles() error = %v, wantErr %v", tt.description, err, tt.wantErr)
				return
			}

			// Verify copied files
			if tt.setupStatic {
				for _, file := range tt.files {
					copiedPath := distDir + "/static/" + file
					if _, err := os.Stat(copiedPath); os.IsNotExist(err) {
						t.Errorf("%s: expected file at %s but it doesn't exist", tt.description, copiedPath)
					} else {
						// Verify content
						data, err := os.ReadFile(copiedPath)
						if err != nil {
							t.Errorf("%s: failed to read copied file %s: %v", tt.description, file, err)
						}
						expectedContent := "test content for " + file
						if string(data) != expectedContent {
							t.Errorf("%s: file %s content mismatch: got %q, want %q",
								tt.description, file, string(data), expectedContent)
						}
					}
				}
			}
		})
	}
}

func TestCopyAvatar(t *testing.T) {
	tests := []struct {
		name        string
		avatar      string
		setupFile   bool
		wantErr     bool
		expectCopy  bool
		description string
	}{
		{
			name:        "local file",
			avatar:      "static/avatar.png",
			setupFile:   true,
			wantErr:     false,
			expectCopy:  true,
			description: "should copy local avatar file",
		},
		{
			name:        "http URL",
			avatar:      "http://example.com/avatar.png",
			setupFile:   false,
			wantErr:     false,
			expectCopy:  false,
			description: "should skip http:// URLs",
		},
		{
			name:        "https URL",
			avatar:      "https://example.com/avatar.png",
			setupFile:   false,
			wantErr:     false,
			expectCopy:  false,
			description: "should skip https:// URLs",
		},
		{
			name:        "missing file",
			avatar:      "static/nonexistent.png",
			setupFile:   false,
			wantErr:     false,
			expectCopy:  false,
			description: "should not fail when file doesn't exist",
		},
		{
			name:        "empty avatar",
			avatar:      "",
			setupFile:   false,
			wantErr:     false,
			expectCopy:  false,
			description: "should handle empty avatar gracefully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			distDir := tmpDir + "/dist"

			// Setup test file if needed
			if tt.setupFile {
				avatarPath := tmpDir + "/" + tt.avatar
				if err := os.MkdirAll(strings.TrimSuffix(avatarPath, "/avatar.png"), 0755); err != nil {
					t.Fatalf("failed to create avatar dir: %v", err)
				}
				testContent := []byte("fake image data")
				if err := os.WriteFile(avatarPath, testContent, 0644); err != nil {
					t.Fatalf("failed to write test avatar: %v", err)
				}
				// Update avatar path to use the temp directory
				tt.avatar = avatarPath
			}

			cfg := &config.Config{
				Name:   "Test User",
				Theme:  "dark",
				Avatar: tt.avatar,
			}

			// Execute CopyAvatar
			err := CopyAvatar(cfg, distDir)

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("%s: CopyAvatar() error = %v, wantErr %v", tt.description, err, tt.wantErr)
				return
			}

			// Check if file was copied
			if tt.expectCopy {
				copiedPath := distDir + "/" + tt.avatar
				if _, err := os.Stat(copiedPath); os.IsNotExist(err) {
					t.Errorf("%s: expected file at %s but it doesn't exist", tt.description, copiedPath)
				}
			}
		})
	}
}
