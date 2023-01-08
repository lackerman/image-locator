package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

func main() {
	var imageDir string
	var csvFile string

	var rootCmd = &cobra.Command{
		Use:   "image-relocator",
		Short: "Move images according to their GPS location and creation date",
		Long:  "Extract GPS metadata from images and move them to directories according to date and location",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Parse the CSV file")
			cities, err := parseCSV(csvFile)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Println("Set up the SQLite database")
			db, err := setupDatabase(cities)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			defer db.Close()

			fmt.Println("Read the directory and process each file")
			err = readDirectory(imageDir, func(filename string, coords *coord) error {
				if coords == nil {
					fmt.Printf("Move %s to default dir\n", filename)
					return nil
				}

				// Get the nearest town or city for the coordinates
				location, err := getLocation(coords, db)
				if err != nil {
					return err
				}

				fmt.Printf("%s (%v) should go to the '%s' directory\n", filename, coords, location)

				return nil
			})
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	rootCmd.Flags().StringVarP(&imageDir, "imageDir", "d", "", "Directory to scan for images (required)")
	if err := rootCmd.MarkFlagRequired("imageDir"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rootCmd.Flags().StringVarP(&csvFile, "csvFile", "g", "geo.db", "SQLite database file for storing location data")
	if err := rootCmd.MarkFlagRequired("csvFile"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
