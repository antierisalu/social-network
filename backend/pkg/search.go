package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"

	db "backend/pkg/db/sqlite"
)

// get basic info of every single user, userID is the ID of the client, will get following relationships.
// isFollowing will be an integer. -1 == not following, 0 == requested, 1 == following.
func FetchAllUsers(userID int) ([]SearchData, error) {
	rows, err := db.DB.Query(`SELECT 
    u.id, 
    u.firstname, 
    u.lastname, 
    u.avatar, 
	u.privacy,
    COALESCE(f1.isFollowing, -1) AS areFollowing, 
    COALESCE(f2.isFollowing, -1) AS isFollowing
FROM 
    users u
LEFT JOIN 
    followers f1 ON u.id = f1.user_id AND f1.follower_id = ?
LEFT JOIN 
    followers f2 ON u.id = f2.follower_id AND f2.user_id = ?
GROUP BY 
    u.id, 
    u.firstname, 
    u.lastname, 
    u.avatar;
`, userID, userID)
	if err != nil {
		return nil, err
	}

	var userArr []SearchData
	for rows.Next() {
		var user SearchData
		err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Avatar, &user.Privacy, &user.AreFollowing, &user.IsFollowing)
		if err != nil {
			fmt.Println("fetchAllUsers: unable to scan user", err, userID)
			continue
		}
		userArr = append(userArr, user)
	}
	return userArr, nil
}

// return every user's basic info to frontend
func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := CheckAuth(r)
	if err != nil {
		// http.Error(w, "paha poiss", http.StatusUnauthorized)
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}

	if r.Method == "GET" {
		userArr, err := FetchAllUsers(userID)
		if err != nil {
			http.Error(w, "Unable to fetch all users", http.StatusBadRequest)
			return
		}
		jsonResponse, err := json.Marshal(userArr)
		if err != nil {
			http.Error(w, "Unable to marshal response", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)

	}
}
