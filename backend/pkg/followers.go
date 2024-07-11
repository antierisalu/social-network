package pkg

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	db "backend/pkg/db/sqlite"
)

// FollowHandler handles the follow requests
func FollowHandler(w http.ResponseWriter, r *http.Request) {

	// PUT method for handling follow request accept and decline
	if r.Method != "PUT" && r.Method != "POST" {
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
	var response struct {
		User         SearchData `json:"user"`
		FollowStatus int        `json:"followStatus"`
	}

	response.User, err = userSearchData(userID)
	if err != nil {
		fmt.Println("FollowHandler: error ", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	if r.Method == "PUT" {
		temp := userID
		userID = requestBody.Target
		requestBody.Target = temp
	}

	response.FollowStatus = requestBody.Action

	if requestBody.Action < 0 {
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
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		if exists {
			err = UpdateRelationship(userID, requestBody.Target, requestBody.Action)
			if err != nil {
				log.Printf("error updating relationship %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
			}/*  else {
				// response = "updated relationship in database: sender " + strconv.Itoa(userID) + ", receiver " + strconv.Itoa(requestBody.Target)
			} */
		} else {
			err = InsertRelationship(userID, requestBody.Target, requestBody.Action)
			if err != nil {
				log.Printf("error inserting relationship %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
			} /* else {
				// response = "inserted relationship to database: sender " + strconv.Itoa(userID) + ", receiver " + strconv.Itoa(requestBody.Target)
			} */
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
	fmt.Println(action, "from: ", userID, " to: ", targetID)
	query := `UPDATE followers SET isFollowing = ? WHERE user_id = ? AND follower_id = ?`
	_, err := db.DB.Exec(query, action, targetID, userID)
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
			 WHERE user_id = ? AND isFollowing = 1`

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

func IsFollowing(targetID, clientID int) (int, error) {
	var relationship int
	query := `SELECT
    COALESCE(
        (SELECT isFollowing
         FROM followers
         WHERE user_id = ? AND follower_id = ?),
        -1
    ) AS isFollowing
`
	err := db.DB.QueryRow(query, targetID, clientID).Scan(&relationship)
	if err != nil {
		return -1, err
	}
	return relationship, nil
}

// Checks if the same notification already exists and inserts a new one if it doesn't.
func InsertNotification(userID int, content, link string) (Notification, error) {
	insertQuery := `INSERT OR REPLACE INTO notifications (user_id, content, link) VALUES (?, ?, ?)`
	res, err := db.DB.Exec(insertQuery, userID, content, link)
	if err != nil {
		log.Printf("error in inserting notification: %v", err)
		return Notification{}, err
	}

	notifID, err := res.LastInsertId()
	if err != nil {
		log.Printf("error in getting last insert id: %v", err)
		return Notification{}, err
	}

	// Fetch the newly inserted notification
	query := `SELECT id, content, link, seen, created_at
	          FROM notifications
	          WHERE id = ?`
	row := db.DB.QueryRow(query, notifID)

	var notification Notification
	err = row.Scan(&notification.ID, &notification.Content, &notification.Link, &notification.Seen, &notification.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return Notification{}, fmt.Errorf("no notification found with id: %d", notifID)
		}
		return Notification{}, err
	}

	linkElements := strings.Split(notification.Link, "_")
	notification.Type = linkElements[0] // Get notification type from link

	id, err := strconv.Atoi(linkElements[1])
	if err != nil {
		log.Printf("Error getting notification FromID: %v", err)
		return Notification{}, err
	}
	notification.FromID = id

	return notification, nil
}

func userSearchData(userID int) (SearchData, error) {
	var user SearchData
	query := `SELECT id, firstname, lastname, avatar FROM users WHERE id = ?`

	err := db.DB.QueryRow(query, userID).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Avatar)
	if err != nil {
		return user, err
	}

	return user, nil
}
