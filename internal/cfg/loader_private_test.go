package cfg

import (
	"testing"
)

func TestLoadPrivateConfig_ExampleFile(t *testing.T) {
	file := "example.private.yml"
	_, err := LoadPrivateConfig(file)
	if err != nil {
		t.Fatalf("Failed to load %s: %v", file, err)
	}
}
