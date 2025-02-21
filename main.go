package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

type FlightData struct {
	Input    []string
	Output   []string
	Database []string // was lookup
}

type AirportFields struct { // was LookupFields
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

	var flightData FlightData
	var airportFields AirportFields

	helperFlag := flag.Bool("h", false, "Display the usage.")
	flag.Parse()
	// If a helper flag is used return with instructions
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

	// Save the argument and welcome the user
	inputPath := os.Args[1]
	outputPath := os.Args[2]
	lookupPath := os.Args[3]

	input := InitializeFile(inputPath)
	lookup := InitializeFile(lookupPath)

	LoadFileContent(input, &flightData.Input)
	LoadFileContent(lookup, &flightData.Database)

	if !ValidateLookup(flightData.Database, &airportFields) {
		fmt.Println(DarkRed + "Error: Airport Lookup malformed")
		return
	}

	// Parse data
	data, _ := ProcessFlightData(&flightData.Input, &flightData.Database, airportFields)

	SaveFileContent(outputPath, data)
}

func getUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func ValidateLookup(lookup []string, fields *AirportFields) bool {

	for i, v := range lookup {
		splits := strings.Split(v, ",")

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
