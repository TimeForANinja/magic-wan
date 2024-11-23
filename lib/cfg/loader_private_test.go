package cfg

import (
	"testing"
)

func TestLoadPrivateConfig_ExampleFile(t *testing.T) {
	_, err := LoadPrivateConfig("example.private.yml")
	if err != nil {
		t.Fatalf("Failed to load example.private.yml: %v", err)
	}
}
