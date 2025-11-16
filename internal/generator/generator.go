// generator.go

/*
 * Copyright (c) - All Rights Reserved.
 *
 * See the LICENCE file for more information.
 */

package generator

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	"linkbeam/internal/config"
)

func GenerateSite(cfg *config.Config, templatePath, outPath string) error {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	outFile, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}
	defer func() { _ = outFile.Close() }()

	if err := tmpl.Execute(outFile, cfg); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}
	return nil
}

// RenderUserPage renders a plain text representation of the user's page.
func RenderUserPage(cfg *config.Config) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Name: %s\n", cfg.Name)
	fmt.Fprintf(&sb, "Bio: %s\n", cfg.Bio)

	sb.WriteString("\nContent Blocks:\n")
	for i, block := range cfg.Content {
		fmt.Fprintf(&sb, "\nBlock %d (type: %s):\n", i+1, block.Type)

		for collectionName, items := range block.Collections {
			fmt.Fprintf(&sb, "  Collection '%s':\n", collectionName)
			for _, item := range items {
				if item.Title != "" {
					fmt.Fprintf(&sb, "    - %s", item.Title)
				}
				if item.URL != "" {
					fmt.Fprintf(&sb, " (url: %s)", item.URL)
				}
				if item.CopyText != "" {
					fmt.Fprintf(&sb, " (copy: %s)", item.CopyText)
				}
				if item.Text != "" {
					fmt.Fprintf(&sb, " (text: %s)", item.Text)
				}
				if item.Icon != "" {
					fmt.Fprintf(&sb, " [%s]", item.Icon)
				}
				sb.WriteString("\n")
			}
		}
	}
	return sb.String()
}

// CopyAssets copies theme and static assets to the output directory.
func CopyAssets(distDir string, themeDirs ...string) error {
	themesDir := filepath.Join(distDir, "themes")
	if err := os.MkdirAll(themesDir, 0755); err != nil {
		return fmt.Errorf("create themes dir: %w", err)
	}

	for _, themeDir := range themeDirs {
		if err := copyThemeFiles(themeDir, themesDir); err != nil {
			return err
		}
	}
	return nil
}

func copyThemeFiles(srcDir, dstDir string) error {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("read theme dir %s: %w", srcDir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		src := filepath.Join(srcDir, entry.Name())
		dst := filepath.Join(dstDir, entry.Name())
		if err := copyFile(src, dst); err != nil {
			return fmt.Errorf("copy %s: %w", entry.Name(), err)
		}
	}
	return nil
}

// copyFile copies a single file from src to dst.
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = srcFile.Close() }()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = dstFile.Close() }()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return dstFile.Sync()
}

// CopyAvatar copies the avatar file to the dist directory if it's a local file.
// It skips URLs (http:// or https://) and doesn't fail if the file doesn't exist.
func CopyAvatar(cfg *config.Config, distDir string) error {
	if cfg.Avatar == "" {
		return nil
	}

	// Skip if avatar is a URL
	if strings.HasPrefix(cfg.Avatar, "http://") || strings.HasPrefix(cfg.Avatar, "https://") {
		return nil
	}

	// Check if the source file exists
	if _, err := os.Stat(cfg.Avatar); os.IsNotExist(err) {
		// Don't fail if file doesn't exist, just skip
		return nil
	}

	// Create the destination path, preserving directory structure
	dstPath := filepath.Join(distDir, cfg.Avatar)

	// Create destination directory if needed
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return fmt.Errorf("create avatar dir: %w", err)
	}

	// Copy the file
	if err := copyFile(cfg.Avatar, dstPath); err != nil {
		return fmt.Errorf("copy avatar: %w", err)
	}

	return nil
}

// CopyStaticFiles copies the static directory to the dist directory.
// It recursively copies all files and subdirectories from static/ to dist/static/.
func CopyStaticFiles(distDir string) error {
	staticSrc := "static"
	staticDst := filepath.Join(distDir, "static")

	// Check if static directory exists
	if _, err := os.Stat(staticSrc); os.IsNotExist(err) {
		// Don't fail if static directory doesn't exist, just skip
		return nil
	}

	// Create destination static directory
	if err := os.MkdirAll(staticDst, 0755); err != nil {
		return fmt.Errorf("create static dir: %w", err)
	}

	// Walk through the static directory and copy all files
	return filepath.Walk(staticSrc, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get relative path from static source
		relPath, err := filepath.Rel(staticSrc, path)
		if err != nil {
			return err
		}

		// Create destination path
		dstPath := filepath.Join(staticDst, relPath)

		// If it's a directory, create it
		if info.IsDir() {
			return os.MkdirAll(dstPath, 0755)
		}

		// If it's a file, copy it
		return copyFile(path, dstPath)
	})
}
