package cfg

import (
	"testing"
)

func TestLoadSharedConfig_ExampleFile(t *testing.T) {
	_, err := LoadSharedConfig("example.shared.yml")
	if err != nil {
		t.Fatalf("Failed to load example.shared.yml: %v", err)
	}
}
