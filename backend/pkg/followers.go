package pkg

import (
	db "backend/pkg/db/sqlite"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// FollowHandler handles the follow requests
func FollowHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method bad", http.StatusMethodNotAllowed)
		return
	}

	userID, err := CheckAuth(r)
	if err != nil {
		log.Printf("Authentication error %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var requestBody struct {
		Action int `json:"action"`
		Target int `json:"target"`
	}
	err = decoder.Decode(&requestBody)
	if err != nil {
		fmt.Println("FollowHandler: badRequest ", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	response, err := userSearchData(userID)
	if err != nil {
		fmt.Println("FollowHandler: error ", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	if requestBody.Action == -1 {
		err = RemoveRelationship(userID, requestBody.Target)
		if err != nil {
			log.Printf("error removing relationship %v", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

	} else {
		exists, err := CheckUserRelationship(userID, requestBody.Target)
		if err != nil {
			log.Printf("error in check user relationship %v", err)
			return
		}

		if exists {
			err = UpdateRelationship(userID, requestBody.Target, requestBody.Action)
			if err != nil {
				log.Printf("error updating relationship %v", err)
			} else {
				//response = "updated relationship in database: sender " + strconv.Itoa(userID) + ", receiver " + strconv.Itoa(requestBody.Target)
			}
		} else {
			err = InsertRelationship(userID, requestBody.Target, requestBody.Action)
			if err != nil {
				log.Printf("error inserting relationship %v", err)
			} else {
				//response = "inserted relationship to database: sender " + strconv.Itoa(userID) + ", receiver " + strconv.Itoa(requestBody.Target)
			}
		}
		}
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			fmt.Println("FollowHandler: error ", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			}
		w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonResponse))
}

// CheckUserRelationship checks if a relationship between two users exists
func CheckUserRelationship(userID, targetID int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM followers WHERE user_id = ? AND follower_id = ?)`
	err := db.DB.QueryRow(query, targetID, userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// RemoveRelationship removes the relationship between two users
func RemoveRelationship(userID, targetID int) error {
	query := `DELETE FROM followers WHERE user_id = ? AND follower_id = ?`
	_, err := db.DB.Exec(query, targetID, userID)
	if err != nil {
		log.Printf("error removing relationship: %v", err)
	}
	return nil
}

// InsertRelationship inserts a new relationship(follow relationship) between two users
func InsertRelationship(userID, targetID, action int) error {
	query := `INSERT INTO followers (user_id, follower_id, isFollowing) VALUES (?, ?, ?)`
	_, err := db.DB.Exec(query, targetID, userID, action)
	if err != nil {
		log.Printf("error inserting relationship: %v", err)
		return err
	}
	return nil
}

// UpdateRelationship updates the relationship between two users
func UpdateRelationship(userID, targetID, action int) error {
	query := `UPDATE followers SET isFollowing = ? WHERE user_id = ? AND follower_id = ?`
	_, err := db.DB.Exec(query, targetID, userID, action)
	if err != nil {
		log.Printf("eroor updating relationship: %v", err)
	}
	return nil
}

// GetAllFollowers returns an array of user structs, that are followers for the given userID
func GetAllFollowers(userID int) ([]SearchData, error) {
	query := `SELECT follower_id, firstname, lastname, avatar
			 FROM followers
			 INNER JOIN users
			 ON followers.follower_id = users.id
			 WHERE user_id = ?`

	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []SearchData
	for rows.Next() {
		var follower SearchData
		if err := rows.Scan(&follower.ID, &follower.FirstName, &follower.LastName, &follower.Avatar); err != nil {
			return nil, err
		}
		followers = append(followers, follower)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return followers, nil
}

// GetAllFollowing returns an array of user structs, that the given userID is following
func GetAllFollowing(userID int) ([]SearchData, error) {
	query := `SELECT user_id, firstname, lastname, avatar
			 FROM followers
			 INNER JOIN users
			 ON followers.user_id = users.id
			 WHERE follower_id = ? AND isFollowing = 1`

	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []SearchData
	for rows.Next() {
		var follower SearchData
		if err := rows.Scan(&follower.ID, &follower.FirstName, &follower.LastName, &follower.Avatar); err != nil {
			return nil, err
		}
		followers = append(followers, follower)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return followers, nil
}

func IsFollowing(userID, targetID int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM followers WHERE user_id = ? AND follower_id = ? AND isFollowing = 1)`
	err := db.DB.QueryRow(query, targetID, userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func userSearchData(userID int)(SearchData, error){
	var user SearchData
	query:= `SELECT id, firstname, lastname, avatar FROM users WHERE id = ?`

	err := db.DB.QueryRow(query, userID).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Avatar)
	if err != nil {
		return user, err
	}

	return user, nil
}