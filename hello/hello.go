package main

import (
	"fmt"
	"log"

	"example.com/greetings"
)

func main() {
	// Get a greeting message and print it.
	log.SetPrefix("greetings: ")
	log.SetFlags(0)
	// Request a greeting message.
	message, err := greetings.Hello("")
	// If an error was returned, print it to the console and
	// exit the program.
	if err != nil {
		log.Fatal(err)
	}
	// If no error was returned, print the returned message
	// to the console.
	fmt.Println(message)
}
