// main.go

/*
 * Copyright (c) - All Rights Reserved.
 *
 * See the LICENSE file for more information.
 */
package main

import (
	"fmt"
	"log"
	"path/filepath"

	flag "github.com/spf13/pflag"

	"linkbeam/internal/config"
	"linkbeam/internal/generator"
)

var (
	// Version is set during build
	Version = "dev"
	// Commit is set during build
	Commit = "unknown"

	fmtPrintf = fmt.Printf
)

func run(loader func(string) (*config.Config, error), configPath string) error {
	cfg, err := loader(configPath)
	if err != nil {
		return err
	}
	_, _ = fmtPrintf("Hello, %s!\n", cfg.Name)
	return nil
}

func main() {
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		_, _ = fmt.Fprintf(out, "LinkBeam - Static site generator for link-in-bio pages\n\n")
		_, _ = fmt.Fprintf(out, "Usage:\n")
		_, _ = fmt.Fprintf(out, "  linkbeam [options]\n\n")
		_, _ = fmt.Fprintf(out, "Options:\n")
		flag.PrintDefaults()
		_, _ = fmt.Fprintf(out, "\nExamples:\n")
		_, _ = fmt.Fprintf(out, "  linkbeam --config myconfig.yaml\n")
		_, _ = fmt.Fprintf(out, "  linkbeam -c config.yaml -o public/index.html\n")
		_, _ = fmt.Fprintf(out, "  linkbeam --version\n")
	}

	configPath := flag.StringP("config", "c", "config.yaml", "path to config YAML file")
	templatePath := flag.StringP("template", "t", "templates/base.html", "path to HTML template")
	outputPath := flag.StringP("output", "o", "dist/index.html", "path to output HTML file")
	showVersion := flag.BoolP("version", "v", false, "print version information")
	flag.Parse()

	if *showVersion {
		fmt.Printf("LinkBeam %s (commit: %s)\n", Version, Commit)
		return
	}

	if err := runMain(*configPath, *templatePath, *outputPath); err != nil {
		log.Fatal(err)
	}
}

func runMain(configPath, templatePath, outputPath string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("config load failed: %w", err)
	}

	// Validate theme against available themes
	if err := cfg.ValidateWithThemes("themes"); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	if err := generator.GenerateSite(cfg, templatePath, outputPath); err != nil {
		return fmt.Errorf("site generation failed: %w", err)
	}

	distDir := filepath.Dir(outputPath)
	if err := generator.CopyAssets(distDir, "themes"); err != nil {
		return fmt.Errorf("asset copy failed: %w", err)
	}

	if err := generator.CopyStaticFiles(distDir); err != nil {
		return fmt.Errorf("static files copy failed: %w", err)
	}

	if err := generator.CopyAvatar(cfg, distDir); err != nil {
		return fmt.Errorf("avatar copy failed: %w", err)
	}

	fmt.Printf("Site generated at %s\n", outputPath)
	return nil
}
