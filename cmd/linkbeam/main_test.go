// main_test.go

/*
 * Copyright (c) - All Rights Reserved.
 *
 * See the LICENCE file for more information.
 */
package main

import (
	"bytes"
	"fmt"
	"testing"

	"linkbeam/internal/config"
)

func TestRun(t *testing.T) {
	mockLoader := func(path string) (*config.Config, error) {
		return &config.Config{Name: "Test User"}, nil
	}

	var buf bytes.Buffer
	fmtPrintf = func(format string, args ...interface{}) (int, error) {
		return fmt.Fprintf(&buf, format, args...)
	}

	err := run(mockLoader, "test-config.yaml")
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}

	got := buf.String()
	want := "Hello, Test User!\n"
	if got != want {
		t.Errorf("unexpected output:\ngot: %q\nwant: %q", got, want)
	}
}
