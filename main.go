package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Information struct {
	Input  []string
	Output []string
	lookup []string
}

type LookupFields struct {
	name         int
	iso_country  int
	municipality int
	icao_code    int
	iata_code    int
	coordinates  int
}

const (
	ClearScreen = "\033[H\033[2J"
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
)

func main() {

	var information Information
	var lookupFields LookupFields

	helperFlag := flag.Bool("h", false, "Display the usage.")
	printFlag := flag.Bool("p", false, "Print flight information on screen.")
	flag.Parse()
	// If a helper flag is used return with instructions
	if *helperFlag {
		fmt.Println("itinerary usage: go run . ./input.txt ./output.txt ./airport-lookup.csv")
		return
	}

	if len(os.Args) > 4 {
		os.Args = append(os.Args[:1], os.Args[2:]...)
	}

	// Check that the program was initialized with enough arguments. If not, then return with a instruction message
	if len(os.Args) < 4 {
		fmt.Println("Too few arguments. \nUsage: go run . ./input.txt ./output.txt ./airport-lookup.csv")
		return
	}

	// Check that input and lookup exists
	if _, err := os.Stat(os.Args[1]); os.IsNotExist(err) {
		fmt.Println("Input file not found. ", err)
		return
	}
	if _, err := os.Stat(os.Args[3]); os.IsNotExist(err) {
		fmt.Println("Lookup file not found. ", err)
		return
	}

	// Save the argument and welcome the user
	inputPath := os.Args[1]
	outputPath := os.Args[2]
	lookupPath := os.Args[3]

	input := OpenFile(inputPath)
	//output := OpenFile(outputPath)
	lookup := OpenFile(lookupPath)

	ReadFromFile(input, &information.Input)
	ReadFromFile(lookup, &information.lookup)

	if !ValidateLookup(information.lookup, &lookupFields) {
		fmt.Println("Error: Airport Lookup malformed")
		return
	}

	// Parse data
	data, print := ParseInput(&information.Input, &information.lookup, lookupFields)

	// Print if user has designated a print flag
	if *printFlag {
		for _, v := range print {
			fmt.Println(v)
		}
	}

	WriteToFile(outputPath, data)

}

func getUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func ValidateLookup(lookup []string, fields *LookupFields) bool {

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

func SetLookupFields(lookup []string, fields *LookupFields) bool {

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
