package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"

	db "backend/pkg/db/sqlite"
)

// only some info for searching purposes
type SearchData struct {
	ID        int
	FirstName string
	LastName  string
	Avatar    string
}

// get basic info of every single user
func FetchAllUsers() ([]SearchData, error) {
	rows, err := db.DB.Query(`SELECT id, firstname, lastname, avatar FROM users `)
	if err != nil {
		return nil, err
	}

	var userArr []SearchData
	for rows.Next() {
		var user SearchData
		err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Avatar)
		if err != nil {
			fmt.Println("fetchAllUsers: unable to scan user", err)
			continue
		}
		userArr = append(userArr, user)
	}
	return userArr, nil
}

// return every user's basic info to frontend
func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	_, err := CheckAuth(r)
	if err != nil {
		// http.Error(w, "paha poiss", http.StatusUnauthorized)
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}

	if r.Method == "GET" {
		userArr, err := FetchAllUsers()
		if err != nil {
			http.Error(w, "TRA", http.StatusBadRequest)
			return
		}
		jsonResponse, err := json.Marshal(userArr)
		if err != nil {
			http.Error(w, "TRA2", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)

	}
}
