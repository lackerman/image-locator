package main

import (
	"database/sql"
	"errors"
	"math"
	"sort"
)

// City represents a city in the CSV file
type City struct {
	City       string  `csv:"city"`
	CityAscii  string  `csv:"city_ascii"`
	Lat        float64 `csv:"lat"`
	Lng        float64 `csv:"lng"`
	Country    string  `csv:"country"`
	Iso2       string  `csv:"iso2"`
	Iso3       string  `csv:"iso3"`
	AdminName  string  `csv:"admin_name"`
	Capital    string  `csv:"capital"`
	Population int     `csv:"population"`
	ID         int     `csv:"id"`
}

// setupDatabase sets up a SQLite database and populates a cities table
func setupDatabase(cities []City) (*sql.DB, error) {
	// Open the SQLite database
	db, err := sql.Open("sqlite3", "app.db")
	if err != nil {
		return nil, err
	}

	// Create the cities table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS cities (
			city TEXT NOT NULL,
			city_ascii TEXT NOT NULL,
			lat REAL NOT NULL,
			lng REAL NOT NULL,
			country TEXT NOT NULL,
			iso2 TEXT NOT NULL,
			iso3 TEXT NOT NULL,
			admin_name TEXT NOT NULL,
			capital TEXT NOT NULL,
			population INTEGER NOT NULL,
			id INTEGER NOT NULL
		)
	`)
	if err != nil {
		return nil, err
	}

	var count int
	// Check if the cities table is already populated
	err = db.QueryRow("SELECT COUNT(*) FROM cities").Scan(&count)
	if err != nil {
		return nil, err
	}
	if count == 42906 {
		return db, nil
	}

	// Insert the cities into the table
	for _, city := range cities {
		_, err := db.Exec("INSERT INTO cities (city, city_ascii, lat, lng, country, iso2, iso3, admin_name, capital, population, id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", city.City, city.CityAscii, city.Lat, city.Lng, city.Country, city.Iso2, city.Iso3, city.AdminName, city.Capital, city.Population, city.ID)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

// getLocation looks up the nearest city for the given coordinates in the database
func getLocation(coord *coord, db *sql.DB) (string, error) {
	rows, err := db.Query(`
	SELECT city, lat, lng
	FROM cities
	WHERE
		(@lat - 0.4) <= lat AND lat <= (@lat + 0.4)
		AND 
		(@lng - 0.4) <= lng AND lng <= (@lng + 0.4)
	`,
		sql.Named("lat", coord.Lat),
		sql.Named("lng", coord.Lon))
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var cities []City
	// Execute the SQL statement
	for rows.Next() {
		var city City
		if err := rows.Scan(&city.City, &city.Lat, &city.Lng); err != nil {
			return "", err
		}
		cities = append(cities, city)
	}

	if len(cities) == 0 {
		return "", errors.New("failed to find any cities in the vicinity of the coordinates")
	}

	sort.Slice(cities, func(i, j int) bool {
		distI := distCalculator(coord.Lat, coord.Lon, cities[i].Lat, cities[i].Lng)
		distJ := distCalculator(coord.Lat, coord.Lon, cities[j].Lat, cities[j].Lng)
		return distI < distJ
	})

	return cities[0].City, nil
}

func distCalculator(centerLat, centerLon, lat, lon float64) float64 {
	return 6371 * math.Acos(
		math.Sin(lat*0.0175)*math.Sin(centerLat*0.0175)+
			math.Cos(lat*0.0175)*math.Cos(centerLat*0.0175)*math.Cos((centerLon*0.0175)-(lon*0.0175)))
}
