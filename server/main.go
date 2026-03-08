package main

import (
	"errors"
	"fmt"
	"log"
	"os"
)

func startup() error {
	defer func() {
		// set up a process for handling panics in startup
		if r := recover(); r != nil {
			log.Fatal("Ran into a critical error during startup:", r)
		}
	}()
	// get the local host for logging real quick
	hostname, err := os.Hostname()
	if err != nil {
		return errors.New("Couldn't get hostname")
	}
	fmt.Println("Starting up Panopticon on host", hostname)

	// create a channel to asynchronously collect errors (any kind of error) from the webserver
	webErrors := make(chan any)

	go StartREST(4001, webErrors) // start our REST API web server with the async channel to collect errors from
	go func() {
		webError := <-webErrors
		if webError != nil {
			fmt.Println("Panic in webserver:", webError)
		}
	}()
	fmt.Println("Startup executed without errors")
	return nil // able to startup without any errors
}

func main() {

	startupErr := startup()
	if startupErr != nil {
		fmt.Println("Fatal error during startup", startupErr)
		os.Exit(1)
	}
	fmt.Println("Locating and pinging watchgroup")
}
