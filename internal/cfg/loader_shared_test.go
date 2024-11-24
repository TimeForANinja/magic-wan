package cfg

import (
	"testing"
)

func TestLoadSharedConfig_ExampleFile(t *testing.T) {
	file := "example.shared.yml"
	_, err := LoadSharedConfig(file)
	if err != nil {
		t.Fatalf("Failed to load %s: %v", file, err)
	}
}
