package yr

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/danieldg91/funtemps/conv"
)

// HOVEDFUNKSJON FOR MAIN

func KonverterFil() error {

	// LAGE NY SKANNER
	scanner := bufio.NewScanner(os.Stdin)

	// FØRSTE INPUT FRA BRUKER
	fmt.Print("Enter output file name: ")
	scanner.Scan()
	outputFileName := scanner.Text()

	// DERSOM FIL(NAVNET) FINNES ALLEREDE, OVERSKRIVE?
	_, err := os.Stat(outputFileName)
	if err == nil {
		fmt.Print("File already exists. Overwrite it? (y/n): ")
		scanner.Scan()
		answer := scanner.Text()
		if answer != "y" {
			return nil
		}
	}

	inputFile, err := os.Open("/minyr/kjevik-temp-celsius-20220318-20230318.csv") // ÅPNE INPUTFILEN
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputFileName) // LAGER OUTPUTFILEN
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	// BUFFERE FOR INPUT (reader) og OUTPUT (writer)
	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()
	reader := bufio.NewReader(inputFile)

	// LESE FØRSTE LINJE FRA INPUT-FILEN OG SKRIVE DET UT TIL OUTPUT (writer) SOM DEN ER.
	firstLine, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	_, err = writer.WriteString(firstLine)
	if err != nil {
		log.Fatal(err)
	}

	for {
		line, err := reader.ReadString('\n') // LESER INPUTFILEN LINJE FOR LINJE
		if err != nil {
			if err.Error() != "EOF" {
				log.Fatal(err)
			}
			break
		}

		line = strings.TrimSpace(line)
		parts := strings.Split(line, ";")                                // Dele linjene opp iht ";"
		temperatureC, err := strconv.ParseFloat(parts[len(parts)-1], 64) // Konverterer data på siste del av linjene
		if err != nil {
			continue // Hopper over linje som ikke er float
		}
		temperatureF := conv.CelsiusToFahrenheit(temperatureC)

		parts[len(parts)-1] = fmt.Sprintln(math.Round(temperatureF*100) / 100) // Oppdaterer siste data på linje før det skrives til ny fil

		newLine := strings.Join(parts, ";") // samle delene 'parts' til en string separert med ";" igjen.

		_, err = writer.WriteString(newLine) // Skrive den nye linjen til fil
		if err != nil {
			log.Fatal(err)
		}
	}

	// ENDRE SISTE LINJE I OUTPUTFIL TIL...
	lastLine := "Data er basert på gyldig data (per 18.03.2023) (CC BY 4.0) fra Meteorologisk institutt (MET);endringen er gjort av Daniel D. Gray\n"
	_, err = writer.WriteString(lastLine)
	if err != nil {
		log.Fatal(err)
	}
	return nil

}

func SeGjennomsnitt() error {

	var sum float64
	var count int

	// LAGE NY SKANNER
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Please enter 'C' or 'F' to see the average temperature in Celsius or Fahrenheit:")

	inputFile, err := os.Open("/minyr/kjevik-temp-celsius-20220318-20230318.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()

	scanner.Scan()
	reader := bufio.NewReader(inputFile)

	tempType := scanner.Text()

	// LESER LINJE FOR LINJE, TRIMMER MELLOMROMMENE, LINJENE DELES OPP ETTER ";" OG SER PÅ SISTE DEL FOR Å KONVERTERE TEMP TIL FAHR.
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() != "EOF" { // ignore EOF error
				log.Fatal(err)
			}
			break
		}

		line = strings.TrimSpace(line)
		parts := strings.Split(line, ";")
		temperature, err := strconv.ParseFloat(parts[len(parts)-1], 64)
		if err != nil {
			continue // hopper over linjer uten float.
		}

		switch tempType {
		case "F":
			temperature = conv.CelsiusToFahrenheit(temperature)
		}

		// Legger sammen tallene for å regne ut gj.snitt.
		sum += temperature
		count++
	}

	// kalkulere/skrive ut gjennomsnittstemp.
	average := sum / float64(count)
	if tempType == "C" {
		fmt.Printf("The average temperature in Celsius is %.2f\n", average)
	} else if tempType == "F" {
		fmt.Printf("The average temperature in Fahrenheit is %.2f\n", average)
	} else {
		fmt.Println("Invalid temperature type entered.")
	}
	os.Exit(0)
	return nil
}
