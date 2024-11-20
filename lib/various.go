package lib

import (
	"bufio"
	"fmt"
	"os"
)

// Helper function to compare slices
func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// GenericFileProcessor reads a file, applies a modification function to each line, and writes it back to the file.
func GenericFileProcessor(filePath string, modifyFunc func(string) string) error {
	// Open the file for reading
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
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
		return fmt.Errorf("failed to read file: %v", err)
	}

	// Open the file for writing
	file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Apply the modification function to each line and write it back
	for _, line := range lines {
		modifiedLine := modifyFunc(line)
		if _, err := writer.WriteString(modifiedLine + "\n"); err != nil {
			return fmt.Errorf("failed to write to file: %v", err)
		}
	}

	// Flush the buffer
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush buffer to file: %v", err)
	}

	return nil
}
