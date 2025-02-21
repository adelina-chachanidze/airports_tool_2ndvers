package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

type Airport struct {
	Name         string
	ISOCountry   string
	Municipality string
	ICAOCode     string
	IATACode     string
	Coordinates  string
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: go run . <input.txt> <output.txt> <airports_lookup.csv>")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]
	airportsFile := os.Args[3]

	// Check if input files exist
	if !fileExists(inputFile) {
		fmt.Printf("Error: Input file %s does not exist\n", inputFile)
		return
	}

	if !fileExists(airportsFile) {
		fmt.Printf("Error: Airports lookup file %s does not exist\n", airportsFile)
		return
	}

	// Load and validate airports data
	airports, err := loadAirports(airportsFile)
	if err != nil {
		fmt.Printf("Error loading airports file: %v\n", err)
		return
	}

	// Process input file
	err = processFile(inputFile, outputFile, airports)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		return
	}
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func loadAirports(filename string) (map[string]Airport, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) < 1 {
		return nil, fmt.Errorf("empty CSV file")
	}

	// Validate header
	requiredHeaders := []string{"name", "iso_country", "municipality", "icao_code", "iata_code", "coordinates"}
	headers := records[0]
	headerMap := make(map[string]int)

	for i, header := range headers {
		headerMap[strings.ToLower(header)] = i
	}

	for _, required := range requiredHeaders {
		if _, exists := headerMap[required]; !exists {
			return nil, fmt.Errorf("missing required column: %s", required)
		}
	}

	airports := make(map[string]Airport)
	for _, record := range records[1:] {
		// Check for empty cells
		for _, field := range record {
			if field == "" {
				return nil, fmt.Errorf("found empty cell in airports data")
			}
		}

		airport := Airport{
			Name:         record[headerMap["name"]],
			ISOCountry:   record[headerMap["iso_country"]],
			Municipality: record[headerMap["municipality"]],
			ICAOCode:     record[headerMap["icao_code"]],
			IATACode:     record[headerMap["iata_code"]],
			Coordinates:  record[headerMap["coordinates"]],
		}

		airports[airport.IATACode] = airport
		airports[airport.ICAOCode] = airport
	}

	return airports, nil
}
