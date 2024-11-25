package osUtil

import (
	"fmt"
	"os"
)

// WriteFile writes a new file with the provided content.
func WriteFile(fileName string, content string) error {
	err := os.WriteFile(fileName, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	return nil
}
