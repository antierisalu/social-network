package pkg

import (
	db "backend/pkg/db/sqlite"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
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
	} else {
		response.Notifications = nil

	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error marshaling notification response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonResponse))
}

func NotifMarkAsSeenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method bad", http.StatusMethodNotAllowed)
		return
	}
	_, err := CheckAuth(r)
	if err != nil {
		log.Printf("Authentication error %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	var requestBody struct {
		NotificationID int `json:"notificationID"`
	}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&requestBody)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	err = markNotificationAsSeen(requestBody.NotificationID)
	if err != nil {
		http.Error(w, "Error marking notification as seen", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

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
		linkElements := strings.Split(notification.Link, "_")

		notification.Type = linkElements[0] // Get notification type from link
		id, err := strconv.Atoi(linkElements[1])
		if err != nil {
			fmt.Println("Error getting notification fromID", err)
		}
		notification.FromID = id
		notifications = append(notifications, notification)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notifications, nil
}

func GetNotificationBasedOnLink(link string) (Notification, error) {
	query := `SELECT id, content, link, seen, created_at
	FROM notifications
	WHERE link = ?`
	row := db.DB.QueryRow(query, link)

	var notification Notification
	err := row.Scan(&notification.ID, &notification.Content, &notification.Link, &notification.Seen, &notification.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return Notification{}, nil // or return a specific error indicating no rows found
		}
		return Notification{}, err
	}

	linkElements := strings.Split(notification.Link, "_")
	notification.Type = linkElements[0] // Get notification type from link

	id, err := strconv.Atoi(linkElements[1])
	if err != nil {
		fmt.Println("Error getting notification fromID", err)
		return Notification{}, err
	}
	notification.FromID = id

	return notification, nil
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

func CheckDuplicateNotification(userID int, link string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM notifications WHERE user_id = ? AND link = ?)`
	err := db.DB.QueryRow(query, userID, link).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func ClearUserNotifications(userID int) error {
	query := `
        DELETE FROM notifications
        WHERE id IN (
            SELECT id
            FROM notifications
            WHERE user_id = ?
            ORDER BY created_at DESC
            LIMIT 10
        )
    `
	_, err := db.DB.Exec(query, userID)
	if err != nil {
		return err
	}
	return nil
}

func clearNotification(userID int) {

	err := ClearUserNotifications(userID)
	if err != nil {
		fmt.Println("Error clearing notifications")
		return
	}
}

func clearSingleNotification(notificationID int) error {

	query := `DELETE FROM notifications WHERE id = ?`
	_, err := db.DB.Exec(query, notificationID)
	if err != nil {
		fmt.Println("Error clearing notification")
		return err
	}
	// fmt.Println("Notification cleared:", notificationID)
	return nil
}

func markNotificationAsSeen(notificationID int) error {

	query := `UPDATE notifications SET seen = 1 WHERE id = ?`
	_, err := db.DB.Exec(query, notificationID)
	if err != nil {
		fmt.Println("Error marking notification as seen")
		return err
	}
	// fmt.Println("Notification marked as seen:", notificationID)
	return nil
}
