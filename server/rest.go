package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"github.com/gin-gonic/gin" // this is a web framework for Go for building web servers and APIs
)

//using gin over net/http (WHICH IS BS BTW) bc its genuinely easier to use + powerful 

/*
	Routes and associated methods:
	/			GET								// Used only to serve a static placeholder
	/notifs		GET, POST, PUT, DELETE			// Used to get, create, and update pending notifs
	/devices	GET, POST, DELETE				// Used to get the devices within a watchgroup
	/enroll		GET								// Used to get keys for enrolling new devices to a watchgroup
*/

// main function to start the server and init REST routes, returns relevant errors if things go wrong
func StartREST(listeningPort int, e chan any, db *sql.DB) {
	router := gin.Default()

	router.GET("/", serveRoot)
	router.GET("/notifs", func(context *gin.Context) { getNotifs(context, db) })
	
	router.Run(":" + strconv.Itoa(listeningPort))
}

func serveRoot(context *gin.Context) {
	body, err := os.ReadFile("static/index.html")
	if err != nil {
		fmt.Println("Couldn't find static file to serve on root")
		context.String(http.StatusInternalServerError, "Internal server error")
		return
	}
	context.Data(http.StatusOK, "text/html", body)
	fmt.Println("Served static root to", context.Request.RemoteAddr)
}

//pull incomplete notifs from db and return them as json
func getNotifs(context *gin.Context, db *sql.DB) {
	rows, _ := db.Query("SELECT id, header, description, status, created_at FROM notifications WHERE status = 'pending' ORDER BY created_at DESC")
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
	context.JSON(http.StatusOK, gin.H{"data": notifs})
}