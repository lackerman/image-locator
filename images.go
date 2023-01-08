package main

import (
	"os/exec"
	"strings"
)

type coord struct {
	Lat float64
	Lon float64
}

// getGPSCoordinates extracts the GPS coordinates from the EXIF data of an image
func getGPSCoordinates(filename string) (*coord, error) {
	// Use exiftool to extract the EXIF data
	out, err := exec.Command("exiftool", "-c", "%.6f", "-GPSPosition", filename).Output()
	if err != nil {
		return nil, err
	}

	// Parse the EXIF data
	return parseExif(string(out)), nil
}

// parseExif parses the output of exiftool and returns the GPS coordinates
func parseExif(exif string) *coord {
	// Split the EXIF data into lines
	lines := strings.Split(exif, "\n")

	// Parse the latitude and longitude
	for _, line := range lines {
		if strings.HasPrefix(line, "GPS Position") {
			coords := strings.Split(line, ":")[1]
			coords = strings.TrimSpace(coords)
			latlon := strings.Split(coords, ",")
			latitude := parseCoordinate(latlon[0])
			longitude := parseCoordinate(latlon[1])
			return &coord{latitude, longitude}
		}
	}

	return nil
}

// parseCoordinate parses a coordinate value from a string
func parseCoordinate(s string) float64 {
	// Split the string into fields
	fields := strings.Fields(s)

	// Parse coordinates
	value := toFloat64(fields[0])
	if fields[1] == "S" || fields[1] == "W" {
		value *= -1
	}

	return value
}
