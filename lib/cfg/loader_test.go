package main

import (
	"testing"
)

func TestLoadConfig_ExampleFile(t *testing.T) {
	_, err := LoadConfig("example.shared.yml")
	if err != nil {
		t.Fatalf("Failed to load example.shared.yml: %v", err)
	}
}
