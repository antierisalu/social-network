package pkg

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	db "backend/pkg/db/sqlite"
)

// POST request updates privacy for own profile
func PrivacyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		userID, err := CheckAuth(r)
		if err != nil {
			fmt.Println("PrivacyHandler: Autherror ", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var requestBody struct {
			NewPrivacy bool `json:"newPrivacy"`
		}
		err = decoder.Decode(&requestBody)
		if err != nil {
			fmt.Println("PrivacyHandler: badRequest ", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		err = updatePrivacy(userID, requestBody.NewPrivacy)
		if err != nil {
			fmt.Println("PrivacyHandler: Unable to update privacy ", err)
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}

		jsonResponse, err := json.Marshal(requestBody)
		if err != nil {
			http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
			return
		}
		w.Write(jsonResponse)
	}
}

// Gets User's profile from db and returns relevant fields based on their privacy
// IF PRIVACY == 0, then profile is public
func fetchUserByID(id int) (*User, error) {
	row := db.DB.QueryRow(`SELECT
            id,
            firstname,
            lastname,
            CASE WHEN privacy = 0 THEN date_of_birth ELSE NULL END AS date_of_birth,
            avatar,
            nickname,
            CASE WHEN privacy = 0 THEN about ELSE NULL END AS about,
            privacy
        FROM users
        WHERE id = ?`, id)

	var user User

	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.DateOfBirth, &user.Avatar, &user.NickName, &user.AboutMe, &user.Privacy)
	if err == sql.ErrNoRows {
		return nil, nil // User not found
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

// return user info when clicking on name in profile search
func GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	_, err := CheckAuth(r)
	if err != nil {
		http.Error(w, "Session not found: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Get the user ID from query parameters
	userIDStr := r.URL.Query().Get("id")
	if userIDStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Convert user ID to integer
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	followers, err := GetAllFollowers(userID)
	if err != nil {
		fmt.Println("Error getting followers", err)
	}
	fmt.Println(followers, userID)
	user, err := fetchUserByID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// change privacy of userID to privacy
func updatePrivacy(userID int, privacy bool) error {
	query := `UPDATE users SET privacy = ? WHERE id = ?`

	result, err := db.DB.Exec(query, privacy, userID)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	fmt.Println(result)

	return nil
}
