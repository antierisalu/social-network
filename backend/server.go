package main

import (
	"fmt"
	"net/http"

	"backend/pkg"
	db "backend/pkg/db/sqlite"
)

func main() {
	// Init Database
	db.DatabaseInit()
	defer db.DB.Close()

	// Serve static files from the current directory
	http.Handle("/", http.FileServer(http.Dir("../public")))

	http.HandleFunc("/login", pkg.LoginHandler)
	http.HandleFunc("/register", pkg.RegisterHandler)
	http.HandleFunc("/session", pkg.SessionHandler)

	// Start the server on port 8080
	fmt.Println("SN is running on http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}
