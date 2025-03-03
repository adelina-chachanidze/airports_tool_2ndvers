package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

// Open the file and give the reference to the file to the main
func InitializeFile(path string) *os.File {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Failed to open file ", err)
	}
	return file
}

// Close the file before exiting the program
func ShutdownFile(file *os.File) {
	file.Close()
}

// Parse the information from the txt file into a Collection struct.
func LoadFileContent(file *os.File, content *[]string) {

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Replace various line endings with newline
		r := strings.NewReplacer("\r", "\n", "\v", "\n", "\f", "\n")
		line = r.Replace(line)
		// Split by newlines
		splits := strings.Split(line, "\n")
		// Add each split directly without normalizing whitespace
		for _, split := range splits {
			*content = append(*content, split)
		}
	}
}

func userErrors() error {
	content, _ := os.ReadFile("output.txt")

	scanner := bufio.NewScanner(bytes.NewReader(content))
	lineNumber := 0
	var errorLines []int

	for scanner.Scan() {
		line := scanner.Text()
		lineNumber++

		// Check for '#' or non-ASCII characters
		if strings.Contains(line, "#") || containsNonASCII(line) {
			errorLines = append(errorLines, lineNumber)
		}
	}

	if len(errorLines) > 0 {
		numbers := make([]string, len(errorLines))
		for i, line := range errorLines {
			numbers[i] = fmt.Sprintf("%d", line)
		}
		fmt.Println("Output file created successfully!\n")
		fmt.Printf("\033[33mPossible errors were detected on line(s) %s in the output file. Please check if formatting is correct in the input file.\033[0m\n",
			strings.Join(numbers, ","))
	}

	return nil
}

// Save the collection into the txt file and convert the values into binary
func SaveFileContent(path string, content []string) error {
	// Open the file for reading and writing, creating it if it doesn't exist
	file := InitializeFile(path)
	// Ensure the file is closed when the function exits
	defer ShutdownFile(file)

	// Move the file pointer to the beginning and clear the file's contents
	file.Seek(0, 0)
	file.Truncate(0)

	// Write each string in the content slice to the file
	for i, v := range content {
		if i == len(content)-1 {
			// Write the last element without a newline
			fmt.Fprintf(file, v)
		} else {
			// Write each element followed by a newline
			fmt.Fprintln(file, v)
		}
	}

	// Check for potential errors in the file content
	if err := userErrors(); err != nil {
		return err
	}
	return nil
}

func containsNonASCII(s string) bool {
	for _, r := range s {
		if r > 127 {
			return true
		}
	}
	return false
}
