package pkg

import (
	db "backend/pkg/db/sqlite"
	"encoding/json"
	"log"
	"net/http"
)

func NotificationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "method bad", http.StatusMethodNotAllowed)
		return
	}

	userID, err := CheckAuth(r)
	if err != nil {
		log.Printf("Authentication error %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var response struct {
		UserID        int            `json:"userID"`
		Notifications []Notification `json:"notifications"`
	}

	response.UserID = userID

	exists, err := CheckNotification(userID)
	if err != nil {
		http.Error(w, "Error checking for notification", http.StatusInternalServerError)
	}

	if exists {
		response.Notifications, err = GetNotifications(userID)
		if err != nil {
			http.Error(w, "Error retrieving notifications", http.StatusInternalServerError)
			return
		}
	}
	response.Notifications = nil
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error marshaling notification response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonResponse))
}

func GetNotifications(userID int) ([]Notification, error) {
	query := `SELECT id, content, link, seen, created_at
			 FROM notifications
			 WHERE user_id = ?`
	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var notification Notification
		if err := rows.Scan(&notification.ID, &notification.Content, &notification.Link, &notification.Seen, &notification.CreatedAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notifications, nil
}

func CheckNotification(userID int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM notifications WHERE user_id = ?)`
	err := db.DB.QueryRow(query, userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// CREATE TABLE "notifications" (
// 	"id"	INTEGER,
// 	"user_id"	INTEGER NOT NULL,
// 	"content"	TEXT NOT NULL,
// 	"link"	TEXT NOT NULL,
// 	"seen"	BOOLEAN NOT NULL DEFAULT 0,
// 	"created_at"	DATE NOT NULL DEFAULT CURRENT_DATE,
// 	FOREIGN KEY("user_id") REFERENCES "users"("id") ON DELETE CASCADE,
// 	PRIMARY KEY("id" AUTOINCREMENT)
// );
