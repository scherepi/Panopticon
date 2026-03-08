package main

import (
	"errors"
	"fmt"
	"os"
)

func startup() error {
	hostname, err := os.Hostname()
	if err != nil {
		return errors.New("Couldn't get hostname")
	}
	fmt.Println("Starting up Panopticon on host", hostname)
	StartREST(4001)
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
