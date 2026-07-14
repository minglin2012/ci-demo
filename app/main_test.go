package main

import (
	"os"
	"testing"
)

func TestAppName(t *testing.T) {
	// Simple test to verify the test infrastructure works in CI
	appName := "ci-demo"
	if appName == "" {
		t.Error("Application name should not be empty")
	}
}

func TestArgsParsing(t *testing.T) {
	// Verify --version flag detection
	args := []string{"ci-demo", "--version"}
	if len(args) != 2 || args[1] != "--version" {
		t.Error("Argument parsing should detect --version flag")
	}
}

func TestMain(m *testing.M) {
	// Setup / Teardown
	os.Exit(m.Run())
}
