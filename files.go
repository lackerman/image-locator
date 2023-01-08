package main

import (
	"encoding/csv"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

func readDirectory(dir string, f func(string, *coord) error) error {
	return filepath.WalkDir(dir, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() {
			coords, err := getGPSCoordinates(path)
			if err != nil {
				return err
			}
			if err := f(path, coords); err != nil {
				return err
			}
		}
		return nil
	})
}

// parseCSV reads a CSV file and returns a slice of City structs
func parseCSV(csvFile string) ([]City, error) {
	// Open the CSV file
	file, err := os.Open(csvFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Read the rows
	var cities []City
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// Parse the row
		city := City{
			City:       row[0],
			CityAscii:  row[1],
			Lat:        toFloat64(row[2]),
			Lng:        toFloat64(row[3]),
			Country:    row[4],
			Iso2:       row[5],
			Iso3:       row[6],
			AdminName:  row[7],
			Capital:    row[8],
			Population: toInt(row[9]),
			ID:         toInt(row[10]),
		}
		cities = append(cities, city)
	}

	return cities, nil
}

// toFloat64 converts a string to a float64, or returns 0 if the string is empty
func toFloat64(s string) float64 {
	if s == "" {
		return 0
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

// toInt converts a string to an int, or returns 0 if the string is empty
func toInt(s string) int {
	if s == "" {
		return 0
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}
