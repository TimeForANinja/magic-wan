package frr

import (
	"bufio"
	"fmt"
	"os"
)

// GenericFileProcessor reads a file, applies a modification function to each line, and writes it back to the file.
func GenericFileProcessor(filePath string, modifyFunc func(string) string) error {
	// Open the file for reading
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Read all lines into a slice
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Check for reading errors
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Open the file for writing
	file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Apply the modification function to each line and write it back
	for _, line := range lines {
		modifiedLine := modifyFunc(line)
		if _, err := writer.WriteString(modifiedLine + "\n"); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}

	// Flush the buffer
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush buffer to file: %w", err)
	}

	return nil
}
