package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Open the file and give the reference to the file to the main
func OpenFile(arg string) *os.File {
	file, err := os.OpenFile(arg, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Failed to open file ", err)
	}
	return file
}

// Close the file before exiting the program
func CloseFile(file *os.File) {
	file.Close()
}

// Parse the information from the txt file into a Collection struct.
func ReadFromFile(file *os.File, data *[]string) {

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		text := scanner.Text()

		r := strings.NewReplacer("\r", "\n", "\v", "\n", "\f", "\n")
		text = r.Replace(text)

		splits := strings.Split(text, "\n")

		*data = append(*data, splits...)
	}
}

// Save the collection into the txt file and convert the values into binary
func WriteToFile(fileName string, data []string) {

	file := OpenFile((fileName))

	file.Seek(0, 0)
	file.Truncate(0)

	for i, v := range data {
		if i == len(data)-1 {
			fmt.Fprintf(file, v)
		} else {
			fmt.Fprintln(file, v)
		}
	}
}
