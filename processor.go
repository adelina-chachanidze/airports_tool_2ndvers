package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Save the collection into the txt file and convert the values into binary
func ProcessFlightData(entries *[]string, database *[]string, fields AirportFields) ([]string, []string) {

	// Initialize slices to store processed and displayable lines

	processed := make([]string, 0)
	displayable := make([]string, 0)

	for _, v := range *entries {

		splits := strings.Split(v, " ")

		processedLine := ""
		displayableLine := ""
		for j, b := range splits {
			if j != 0 {
				processedLine += " "
				displayableLine += " "
			}

			// City search if there's a * before # or ##
			if strings.Contains(b, "*##") {

				pattern := regexp.MustCompile(`\*##[A-Z]{4}\b`)
				match := pattern.FindString(b)
				value := ""

				// ICAO Codes are 4 letters long + 3 from prefix. If it's not then it's not a valid ICAO code
				if len(match) == 7 {
					value = FetchCityName(*database, match, fields.icao_code, fields)
				} else {
					value = "false"
				}

				if value == "false" {
					processedLine += b
					displayableLine += b
				} else {
					processedLine += strings.Replace(b, match, value, -1)
					displayableLine += strings.Replace(b, match, value, -1)
				}

			} else if strings.Contains(b, "*#") {

				// Extract the IATA code from the line
				pattern := regexp.MustCompile(`\*#[A-Z]{3}\b`)
				match := pattern.FindString(b)
				value := ""

				// IATA Codes are 3 letters long + 2 from prefix. If it's not then it's not a valid ICAO code.
				if len(match) == 5 {
					value = FetchCityName(*database, match, fields.iata_code, fields)
				} else {
					value = "false"
				}

				if value == "false" {
					processedLine += b
					displayableLine += b
				} else {
					processedLine += strings.Replace(b, match, value, -1)
					displayableLine += strings.Replace(b, match, value, -1)
				}

			} else if strings.Contains(b, "##") {

				pattern := regexp.MustCompile(`\##[A-Z]{4}\b`)
				match := pattern.FindString(b)
				value := ""

				// ICAO Codes are 4 letters long + 2 from prefix. If it's not then it's not a valid ICAO code
				if len(match) == 6 {
					value = FetchAirportName(*database, match, fields.icao_code, fields)
				} else {
					value = "false"
				}

				if value == "false" {
					processedLine += b
					displayableLine += b
				} else {
					processedLine += strings.Replace(b, match, value, -1)
					displayableLine += strings.Replace(b, match, value, -1)
				}

			} else if strings.Contains(b, "#") {

				pattern := regexp.MustCompile(`\#[A-Z]{3}\b`)
				match := pattern.FindString(b)
				value := ""

				// IATA Codes are 3 letters long + 1 from prefix. If it's not then it's not a valid IATA Code
				if len(match) == 4 {
					value = FetchAirportName(*database, match, fields.iata_code, fields)
				} else {
					value = "false"
				}

				if value == "false" {
					processedLine += b
					displayableLine += b
				} else {
					processedLine += strings.Replace(b, match, value, -1)
					displayableLine += strings.Replace(b, match, value, -1)
				}
				//fmt.Println("##: " + value)

			} else if strings.Contains(b, "T24") {

				pattern := regexp.MustCompile(`T24\((.*?)\)`)
				match := pattern.FindString(b)

				formattedDate := FormatDateTime(match, "T24")
				if formattedDate == "false" {
					processedLine += b
					displayableLine += b
				} else {
					processedLine += strings.Replace(b, match, formattedDate, -1)
					displayableLine += strings.Replace(b, match, formattedDate, -1)
				}

			} else if strings.Contains(b, "T12") {

				pattern := regexp.MustCompile(`T12\((.*?)\)`)
				match := pattern.FindString(b)

				formattedDate := FormatDateTime(match, "T12")
				if formattedDate == "false" {
					processedLine += b
					displayableLine += b
				} else {
					processedLine += strings.Replace(b, match, formattedDate, -1)
					displayableLine += strings.Replace(b, match, formattedDate, -1)
				}

			} else if strings.Contains(b, "D") {

				pattern := regexp.MustCompile(`D\((.*?)\)`)
				match := pattern.FindString(b)

				formattedDate := FormatDateTime(match, "D")
				if formattedDate == "false" {
					processedLine += b
					displayableLine += b
				} else {
					processedLine += strings.Replace(b, match, formattedDate, -1)
					displayableLine += strings.Replace(b, match, formattedDate, -1)
				}

			} else {
				processedLine += b
				displayableLine += b
			}
		}

		// Trim newlines (vertical whitespaces)
		if len(processed) > 0 {
			if len(processed[len(processed)-1]) == 0 && len(processedLine) == 0 {
				// Do nothing
			} else {
				processed = append(processed, processedLine)
				displayable = append(displayable, displayableLine)
			}
		} else {
			processed = append(processed, processedLine)
			displayable = append(displayable, displayableLine)
		}

	}
	return processed, displayable
}

// FetchAirportName returns an airport name for a given IATA/ICAO code from the database
func FetchAirportName(database []string, code string, codeType int, fields AirportFields) string {

	code = strings.Trim(code, "#()*,.")
	for _, v := range database {
		splits := strings.Split(v, ",")

		if len(splits) >= 6 {
			if splits[codeType] == code {
				return splits[fields.name]
			}
		} else {
			fmt.Println("Airport Lookup Malformed. Aborting!")
			break
		}
	}
	return "false"
}

// FetchCityName returns a city name for a given IATA/ICAO code from the database
func FetchCityName(database []string, code string, codeType int, fields AirportFields) string {

	code = strings.Trim(code, "#()*,.")
	for _, v := range database {
		splits := strings.Split(v, ",")

		if len(splits) >= 6 {
			if splits[codeType] == code {
				return splits[fields.municipality]
			}
		} else {
			fmt.Println("Airport Lookup Malformed. Aborting!")
			break
		}
	}
	return "false"
}

// FormatDateTime converts timestamps between different time and date formats
func FormatDateTime(timestamp string, format string) string {

	timestamp = strings.Replace(timestamp, "Z", "+00:00", -1)
	timestamp = strings.Replace(timestamp, "T", "-", -1)
	timestamp = strings.Replace(timestamp, "+", "-+", -1)

	if format == "D" {
		timestamp = strings.Trim(timestamp, "D()")
		dateSplits := strings.Split(timestamp, "-")

		if len(dateSplits) < 3 {
			return "false"
		}
		month := ConvertMonthToAbbrev(dateSplits[1])
		if month == "false" {
			return month
		}
		formattedDate := dateSplits[2] + " " + month + " " + dateSplits[0]

		return formattedDate

	} else if format == "T12" {
		// Split the data into slices
		dateSplits := strings.Split(timestamp, "-")

		// If there are not at least 3 slices then the data is not in the T12(2080-05-04T14:54Z) format
		if len(dateSplits) < 3 {
			return "false"
		}
		temporatyData := strings.Trim(dateSplits[len(dateSplits)-2], "+:()")

		// Remove the :-character from between the numbers and parse it into an integer for easier handling
		temporatyDataReplaced := strings.Replace(temporatyData, ":", "", -1)

		if len(temporatyDataReplaced) != 4 {
			return "false"
		}
		parsedInt, err := strconv.Atoi(temporatyDataReplaced)

		if err != nil {
			return "false"
		}

		// SuffixCheck. See if it's AM or PM and transform the PM time into 12h format
		timeSuffix := ""
		if parsedInt >= 1200 {
			if parsedInt > 1200 {
				parsedInt = parsedInt - 1200
			}
			timeSuffix = "PM"
		} else {
			timeSuffix = "AM"
			if parsedInt == 0000 {
				parsedInt = 1200
			}
		}

		// String conversion and adding of 0 in front in case it is not in the 0000 format
		parsedTime := strconv.Itoa(parsedInt)
		if len(parsedTime) < 4 {
			parsedTime = "0" + parsedTime
		}

		// Divide the 0000 time into 00:00 format and return it
		timeFormatted := ""
		for i := 0; i < len(parsedTime); i++ {
			if i == 2 {
				timeFormatted += ":"
			}
			timeFormatted += string(parsedTime[i])
		}
		if !strings.Contains(dateSplits[len(dateSplits)-1], "+") {
			dateSplits[len(dateSplits)-1] = "-" + dateSplits[len(dateSplits)-1]
		}
		formattedDate := timeFormatted + timeSuffix + " (" + dateSplits[len(dateSplits)-1]
		return formattedDate

	} else if format == "T24" {
		dateSplits := strings.Split(timestamp, "-")
		if len(dateSplits) < 3 {
			return "false"
		}
		if !strings.Contains(dateSplits[len(dateSplits)-1], "+") {
			dateSplits[len(dateSplits)-1] = "-" + dateSplits[len(dateSplits)-1]
		}
		formattedDate := dateSplits[len(dateSplits)-2] + " (" + dateSplits[len(dateSplits)-1]
		return formattedDate

	} else {
		return "false"
	}
}

// ConvertMonthToAbbrev converts numeric months (01-12) to abbreviated names (Jan-Dec)
func ConvertMonthToAbbrev(month string) string {
	switch month {
	case "01":
		return "Jan"
	case "02":
		return "Feb"
	case "03":
		return "Mar"
	case "04":
		return "Apr"
	case "05":
		return "May"
	case "06":
		return "Jun"
	case "07":
		return "Jul"
	case "08":
		return "Aug"
	case "09":
		return "Sep"
	case "10":
		return "Oct"
	case "11":
		return "Nov"
	case "12":
		return "Dec"
	default:
		return "false"
	}
}
