# Flight Itinerary Program

## Overview
The Flight Itinerary Program processes flight itinerary data by replacing airport codes with their corresponding airport or city names using a lookup database. The program reads input data, processes flight information, and saves the results to an output file.

## Features
- Reads flight itinerary data from an input file.
- Matches airport codes (IATA/ICAO) to airport names and city names from a lookup database.
- Formats timestamps into human-readable formats.
- Saves the processed data into an output file.
- Provides error detection for potential formatting issues.

## Usage
### Running the Program
```
go run . ./input.txt ./output.txt ./airport-lookup.csv
```

### Command-Line Arguments
- `input.txt` - File containing flight data entries.
- `output.txt` - File where processed flight data will be saved.
- `airport-lookup.csv` - CSV file containing airport and city data.

### Optional Flags
- `-h` - Displays usage instructions.

## Input Format
The input file should contain flight-related text with airport codes and timestamps in the following formats:
- `#IATA` (e.g., `#JFK`) – Replaced with airport name.
- `##ICAO` (e.g., `##KJFK`) – Replaced with airport name.
- `*#IATA` – Replaced with city name.
- `*##ICAO` – Replaced with city name.
- `D(YYYY-MM-DD)` – Reformats the date.
- `T12(HH:MM)` – Converts time to 12-hour format, including zulu time.
- `T24(HH:MM)` – Converts time to 24-hour format, including zulu time.

## Output Format
The processed output will replace airport codes and timestamps with their respective names and formatted time values.

## Error Handling
- If the lookup file is malformed or missing data, the program will terminate with an error message.
- If the output file contains unusual characters or symbols, warnings will be displayed with line numbers.

## Example
### Input (`input.txt`)
```
Flight from #JFK to ##KJFK departs at T24(14:30Z) on D(2025-03-15).
```

### Lookup (`airport-lookup.csv`)
```
name,iso_country,municipality,icao_code,iata_code,coordinates
John F Kennedy International,US,New York,KJFK,JFK,40.6413,-73.7781
```

### Output (`output.txt`)
```
Flight from John F Kennedy International to John F Kennedy International departs at 14:30 (+00:00) on 15 Mar 2025.
```

## Dependencies
- Go 1.18 or later

