package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/danieldg91/funtemps/conv"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("")
	fmt.Println("Running...")
	fmt.Println("Enter 'q' or 'exit' to close the program")
	fmt.Println("Enter 'convert' to convert the file to output.csv in fahrenheit values")

	for scanner.Scan() {
		input := scanner.Text()
		if input == "q" || input == "exit" {

			os.Exit(0)
		}
		if input == "convert" {

			// Open the CSV file
			file, err := os.Open("kjevik-temp-celsius-20220318-20230318.csv")
			if err != nil {
				fmt.Println("Error opening file:", err)
				return
			}
			defer file.Close()

			// Creates a CSV reader
			reader := csv.NewReader(file)
			reader.Comma = ';'

			// Reads the CSV records
			records, err := reader.ReadAll()
			if err != nil {
				fmt.Println("Error reading CSV file:", err)
				return
			}

			outFile, err := os.Create("output.csv")
			if err != nil {
				fmt.Println("Error creating output file:", err)
				return
			}
			defer outFile.Close()

			// Create a CSV writer
			writer := csv.NewWriter(outFile)
			writer.Comma = ';'

			// Loop through the records and convert temperatures to Fahrenheit
			var sum float64
			var count int
			for _, record := range records {
				// Skip the rows without any temperature data
				if record[0] == "Navn" || record[0] == "Data er gyldig per" {
					continue
				}

				// Extract the temperature value and convert it to float
				temp, err := strconv.ParseFloat(record[3], 64)
				if err != nil {
					fmt.Println("Error parsing temperature:", err)
					continue
				}

				// Convert temperature from Celsius to Fahrenheit
				fahrenheit := conv.CelsiusToFahrenheit(temp)

				record[3] = strconv.FormatFloat(fahrenheit, 'f', 1, 64)

				// Write the updated record to the output CSV file
				err = writer.Write(record)
				if err != nil {
					fmt.Println("Error writing record to output file:", err)
					continue
				}

				// Add the temperature to the sum and increment the coexunt
				sum += temp
				count++

				//Ensures that the data is written at the right moment
				writer.Flush()
			}

			// Calculates the average temperature
			if count > 0 {
				avg := sum / float64(count)
				fahravg := conv.CelsiusToFahrenheit(avg)
				fmt.Printf("Average temperature: %.1f°C, or %.1f°F\n", avg, fahravg)
			}

			fmt.Println("The result of the conversion has been written to the file: 'output.csv'")
		}
	}

}
