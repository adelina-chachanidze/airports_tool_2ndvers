package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
)

// FlightData holds the input, output and database content for flight processing
type FlightData struct {
	Input    []string
	Output   []string
	Database []string
}

// defines what data needs to be extracted from the CSV file
type AirportFields struct {
	name         int
	iso_country  int
	municipality int
	icao_code    int
	iata_code    int
	coordinates  int
}

const (
	Yellow  = "\033[33m"
	DarkRed = "\033[31;2m"
)

func main() {
	fmt.Println("Welcome to the Flight Itinerary Program!")

	//creates a new FlightData struct named flightData.
	//This struct will hold flight-related data, and start as empty slices of strings.
	var flightData FlightData
	var airportFields AirportFields

	helperFlag := flag.Bool("h", false, "Display the usage.")
	flag.Parse()
	if *helperFlag {
		fmt.Println("itinerary usage: go run . ./input.txt ./output.txt ./airport-lookup.csv")
		return
	}

	if len(os.Args) > 4 {
		os.Args = append(os.Args[:1], os.Args[2:]...)
	}

	// Check that the program was initialized with enough arguments
	if len(os.Args) < 4 {
		fmt.Println(DarkRed + "Too few arguments. \nUsage: go run . ./input.txt ./output.txt ./airport-lookup.csv")
		return
	}

	// Check that input and lookup exists
	if _, err := os.Stat(os.Args[1]); os.IsNotExist(err) {
		fmt.Println(DarkRed + "Input file not found. " + err.Error())
		return
	}
	if _, err := os.Stat(os.Args[3]); os.IsNotExist(err) {
		fmt.Println(DarkRed + "Lookup file not found. " + err.Error())
		return
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]
	lookupPath := os.Args[3]

	input := InitializeFile(inputPath)
	lookup := InitializeFile(lookupPath)

	// Loads the content of the input file into flightData.Input
	LoadFileContent(input, &flightData.Input)
	// Loads the content of the lookup file into flightData.Database
	LoadFileContent(lookup, &flightData.Database)

	if !ValidateLookup(flightData.Database, &airportFields) {
		fmt.Println(DarkRed + "Error: Airport Lookup malformed")
		return
	}

	data, _ := ProcessFlightData(&flightData.Input, &flightData.Database, airportFields)

	SaveFileContent(outputPath, data)
}

// getUserInput reads and returns a line of input from the user
func getUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// ValidateLookup checks if the airport database has valid formatting
func ValidateLookup(lookup []string, fields *AirportFields) bool {
	// Ensure there is at least one line for the header
	if len(lookup) == 0 {
		return false
	}

	// Expected number of fields (not commas)
	expectedFields := 6

	for i, line := range lookup {
		// Parse the CSV line properly to handle quoted fields
		r := csv.NewReader(strings.NewReader(line))
		splits, err := r.Read()

		// Check if parsing failed or field count is wrong
		if err != nil || len(splits) != expectedFields {
			return false
		}

		// Get the order of the fields
		if i == 0 {
			SetLookupFields(splits, fields)
		}

		// Check for empty fields
		for _, b := range splits {
			if len(b) == 0 {
				return false
			}
		}
	}
	return true
}

// sets up the struct based on the actual data layout in the CSV file
func SetLookupFields(lookup []string, fields *AirportFields) bool {

	for i, v := range lookup {
		switch v {
		case "name":
			fields.name = i
		case "iso_country":
			fields.iso_country = i
		case "municipality":
			fields.municipality = i
		case "icao_code":
			fields.icao_code = i
		case "iata_code":
			fields.iata_code = i
		case "coordinates":
			fields.coordinates = i
		default:
			return false
		}
	}
	return true
}
