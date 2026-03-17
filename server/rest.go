package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

/*
	Routes and associated methods:
	/			GET								// Used only to serve a static placeholder
	/notifs		GET, POST, PUT, DELETE			// Used to get, create, and update pending notifs
	/devices	GET, POST, DELETE				// Used to get the devices within a watchgroup
	/enroll		GET								// Used to get keys for enrolling new devices to a watchgroup
*/

// main function to start the server and init REST routes, returns relevant errors if things go wrong
func StartREST(listeningPort int, e chan any, db *sql.DB) {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", recoverHandler(e, serveRoot))
	mux.HandleFunc("GET /notifs", recoverHandler(e, func(w http.ResponseWriter, req *http.Request) {
		getNotifs(w, req, db)
	}))
	
	//run server (makes sure it blocks so server keeps running)
	addr := ":" + strconv.Itoa(listeningPort)
	if err := http.ListenAndServe(addr, mux); err != nil {e <- err}
}

// wraps handler to catch panics and push them to the error channel - it wont work the prev way you had it 
func recoverHandler(e chan any, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("PANIC in webserver, pushing upstream")
			e <- r
		}
	}()
		next(w, req)
	}
}


//serve the static placeholder
func serveRoot(w http.ResponseWriter, req *http.Request) {
	//i think http.ServeFile would be better here but it won't have as explicit logging but idk if that's a big deal
		// if you do wanna use it just uncomment the line below and comment out line 53-65
	//http.ServeFile(w, req, "static/index.html")
	body, err := os.ReadFile("static/index.html")
	if err != nil {
		fmt.Println("Couldn't find static file to serve on root")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(body)
	if err != nil {
		fmt.Println("Error writing response body when serving static file")
		return
	}
	fmt.Println("Served static root to", req.RemoteAddr)
}

//pulls incomplete notifs from db and returns them as json
func getNotifs(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	rows, err := db.Query("SELECT id, header, description, status, created_at FROM notifications WHERE status = 'pending' ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var notifs []Notification
	for rows.Next() {
		var notif struct {
			ID          int
			Header      string
			Description string
			Status      string
			CreatedAt   string
		}
		rows.Scan(&notif.ID, &notif.Header, &notif.Description, &notif.Status, &notif.CreatedAt)
		notifs = append(notifs, Notification{Header: notif.Header, Description: notif.Description})
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": notifs})
	fmt.Println("Served notifs to", req.RemoteAddr)
}

//helper to write json
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
