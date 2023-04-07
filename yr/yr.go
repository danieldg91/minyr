package yr

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/danieldg91/funtemps/conv"
)

// -------------------------------------------------------------------- JANIS KODE UNDER DENNE LINJEN
func CelsiusToFahrenheitString(celsius string) (string, error) {
	var fahrFloat float64
	var err error
	if celsiusFloat, err := strconv.ParseFloat(celsius, 64); err == nil {
		fahrFloat = conv.CelsiusToFahrenheit(celsiusFloat)
	}
	fahrString := fmt.Sprintf("%.1f", fahrFloat)
	return fahrString, err
}

// Forutsetter at vi kjenner strukturen i filen og denne implementasjon
// er kun for filer som inneholder linjer hvor det fjerde element
// på linjen er verdien for temperaturaaling i grader celsius
func CelsiusToFahrenheitLine(line string) (string, error) {

	dividedString := strings.Split(line, ";")
	var err error

	if len(dividedString) == 4 {
		dividedString[3], err = CelsiusToFahrenheitString(dividedString[3])
		if err != nil {
			return "", err
		}
	} else {
		return "", errors.New("linje har ikke forventet format")
	}
	return strings.Join(dividedString, ";"), nil

	// return "Kjevik;SN39040;18.03.2022 01:50;42.8", err
}

// -------------------------------------------------------------------- JANIS' KODE OVER DENNE LINJEN

// HOVEDFUNKSJON FOR MAIN

func KonverterFil() error {

	// LAGE NY SKANNER
	scanner := bufio.NewScanner(os.Stdin)

	// TILLATE INPUT FRA BRUKER
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

	inputFile, err := os.Open("kjevik-temp-celsius-20220318-20230318.csv") // ÅPNE INPUTFILEN
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

	// BEHOLDE FØRSTE LINJE FRA INPUT-FILEN ...
	firstLine, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	//... OG SKRIVE DET UT TIL INPUT (writer) SOM DEN ER.
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

	// GJENNOMSNITTSTEMPERATUR (C&F)
	fmt.Println("Also, enter 'C' or 'F' to see the average temperature in Celsius or Fahrenheit:")
	scanner.Scan()
	tempType := scanner.Text()

	var sum float64
	var count int

	inputFile, err = os.Open("kjevik-temp-celsius-20220318-20230318.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()

	reader = bufio.NewReader(inputFile)

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

		if tempType == "F" {
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

	return nil
}
