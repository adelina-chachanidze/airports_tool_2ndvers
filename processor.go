package main

import (
	"bufio"
	"os"
	"regexp"
	"strings"
	"time"
)

func processFile(inputFile, outputFile string, airports map[string]Airport) error {
	input, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer input.Close()

	output, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer output.Close()

	scanner := bufio.NewScanner(input)
	writer := bufio.NewWriter(output)
	defer writer.Flush()

	// Compile regular expressions
	airportPattern := regexp.MustCompile(`(#[A-Z]{3}|##[A-Z]{4})`)
	cityPattern := regexp.MustCompile(`\*(#[A-Z]{3}|##[A-Z]{4})`)
	datePattern := regexp.MustCompile(`D\((\d{4}-\d{2}-\d{2})T\d{2}:\d{2}[−-]\d{2}:\d{2}\)`)
	time12Pattern := regexp.MustCompile(`T12\((\d{4}-\d{2}-\d{2}T\d{2}:\d{2}(?:[−-]\d{2}:\d{2}|Z))\)`)
	time24Pattern := regexp.MustCompile(`T24\((\d{4}-\d{2}-\d{2}T\d{2}:\d{2}(?:[−-]\d{2}:\d{2}|Z))\)`)

	for scanner.Scan() {
		line := scanner.Text()

		// Process dates
		line = datePattern.ReplaceAllStringFunc(line, convertDate)

		// Process times
		line = time12Pattern.ReplaceAllStringFunc(line, func(s string) string {
			return convertTime(s, true)
		})
		line = time24Pattern.ReplaceAllStringFunc(line, func(s string) string {
			return convertTime(s, false)
		})

		// Process city patterns first (they have *)
		line = cityPattern.ReplaceAllStringFunc(line, func(s string) string {
			code := strings.TrimPrefix(s, "*")
			if airport, ok := airports[strings.Trim(code, "#")]; ok {
				return airport.Municipality
			}
			return s
		})

		// Process airport patterns
		line = airportPattern.ReplaceAllStringFunc(line, func(s string) string {
			if airport, ok := airports[strings.Trim(s, "#")]; ok {
				return airport.Name
			}
			return s
		})

		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return scanner.Err()
}

func convertDate(s string) string {
	dateStr := s[2 : len(s)-1] // Remove D() wrapper
	t, err := time.Parse("2006-01-02T15:04-07:00", dateStr)
	if err != nil {
		return s
	}
	return t.Format("02 Jan 2006")
}

func convertTime(s string, is12Hour bool) string {
	// Extract the time string
	timeStr := s[4 : len(s)-1] // Remove T12() or T24() wrapper

	var t time.Time
	var err error

	// Handle both Zulu time and offset time
	if strings.HasSuffix(timeStr, "Z") {
		t, err = time.Parse("2006-01-02T15:04Z", timeStr)
		if err != nil {
			return s
		}
		if is12Hour {
			return t.Format("03:04PM (+00:00)")
		}
		return t.Format("15:04 (+00:00)")
	}

	t, err = time.Parse("2006-01-02T15:04-07:00", timeStr)
	if err != nil {
		return s
	}

	offset := t.Format("-07:00")
	if is12Hour {
		return t.Format("03:04PM") + " (" + offset + ")"
	}
	return t.Format("15:04") + " (" + offset + ")"
}
