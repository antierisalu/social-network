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
	http.Handle("/avatars/", http.StripPrefix("/avatars", http.FileServer(http.Dir("./avatars"))))

	http.HandleFunc("/login", pkg.LoginHandler)
	http.HandleFunc("/register", pkg.RegisterHandler)
	http.HandleFunc("/session", pkg.SessionHandler)
	http.HandleFunc("/privacy", pkg.PrivacyHandler)
	http.HandleFunc("/user", pkg.GetUserInfoHandler)
	http.HandleFunc("/allusers", pkg.GetAllUsersHandler)
	http.HandleFunc("/editProfile", pkg.ProfileEditorHandler)

	// websocket
	http.HandleFunc("/ws", pkg.WsHandler)
	
	// Start the server on port 8080
	fmt.Println("SN is running on http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}
