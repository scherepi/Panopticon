package main

import (
	"fmt"
	"net/http"
	"os"
)

/*
	Routes and associated methods:
	/			GET								// Used only to serve a static placeholder
	/notifs		GET, POST, PUT, DELETE			// Used to get, create, and update pending notifs
	/devices	GET, POST, DELETE				// Used to get the devices within a watchgroup
	/enroll		GET								// Used to get keys for enrolling new devices to a watchgroup
*/

// main function to start the server and init REST routes, returns relevant errors if things go wrong
func StartREST(listeningPort int) error {
	mux := http.NewServeMux()

	// Set up the HTTP routes

	mux.HandleFunc("GET /", serveRoot) //
	mux.HandleFunc("GET /notifs", serveNotifs)

	http.ListenAndServe(":"+string(listeningPort), mux)
	return nil
}

// Route handling functions

func serveRoot(w http.ResponseWriter, req *http.Request) {
	// read the static file into a variable real quick
	body, err := os.ReadFile("static/index.html")
	if err != nil {
		fmt.Println("Couldn't find static file to serve on root")
		return
	}
	_, err = fmt.Fprintf(w, string(body)) // write the static file to the connection just as-is
	if err != nil {
		fmt.Println("Error writing response body when serving static file")
		return
	}
	fmt.Println("Served static root to ", req.RemoteAddr) // log the interaction for the fuck of it
}

func serveNotifs(w http.ResponseWriter, req *http.Request) {

}
