package main

import (
	"bufio"
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
		// Process each split to remove extra spaces
		for _, split := range splits {
			// Fields splits by whitespace and Join reconstructs with single spaces
			normalized := strings.Join(strings.Fields(split), " ")
			*content = append(*content, normalized)
		}
	}
}

// Save the collection into the txt file and convert the values into binary
func SaveFileContent(path string, content []string) {
	file := InitializeFile(path)

	file.Seek(0, 0)
	file.Truncate(0)

	for i, v := range content {
		if i == len(content)-1 {
			fmt.Fprintf(file, v)
		} else {
			fmt.Fprintln(file, v)
		}
	}
}
