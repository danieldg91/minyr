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
	fmt.Println("Enter 'q' or 'exit' to close the program, or")
	fmt.Println("Enter 'convert' to run the program:")

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
			return
		}
	}
}
