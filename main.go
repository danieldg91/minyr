package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/danieldg91/minyr/yr"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("")
	fmt.Println("Running...")
	fmt.Println("")
	fmt.Println("Enter 'q' or 'exit' to close the program.")
	fmt.Println("Enter 'convert' to run the program and convert the input file to Fahrenheit values.")
	fmt.Println("Enter 'average' to see the average temperature in Celcius or Fahrenheit.")

	for scanner.Scan() {
		input := scanner.Text()

		switch input {
		case "q", "exit":
			os.Exit(0)
		case "convert":
			err := yr.KonverterFil()
			if err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		case "average":
			err := yr.SeGjennomsnitt()
			if err != nil {
				log.Fatal(err)
			}

			return
		}
	}
}
