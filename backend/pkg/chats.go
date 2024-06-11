package pkg

import (
	"fmt"
	"time"

	db "backend/pkg/db/sqlite"
)

// Creates a new database chat entry for user1 & user2
func InsertNewChat(user1, user2 int) error {
	stmt, err := db.DB.Prepare("INSERT INTO chats (user1, user2, created_at) VALUES (?, ?, ?)")
	if err != nil {
		fmt.Println("Error preparing statement in InsertNewChat:", err)
		return err
	}
	defer stmt.Close()
	if user1 > user2 {
		user1, user2 = user2, user1
	}
	_, err = stmt.Exec(user1, user2, time.Now())
	if err != nil {
		fmt.Println("Error executing statement in InsertNewChat:", err)
		return err
	}
	return nil
}

// Takes in user1 ID and user2ID to check if there is a database entry for their chat,
// IF entry exists returns the chatID if not returns -1, error
func GetChatID(userID1, userID2 int) (int, error) {
	stmt := "SELECT * FROM chats"
	rows, err := db.DB.Query(stmt)
	if err != nil {
		return -1, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			chatID      int
			user1       int
			user2       int
			createdate  string
			lastmessage string
		)
		if err := rows.Scan(&chatID, &user1, &user2, &lastmessage, &createdate); err != nil {
			fmt.Println("Error getting row values from db:", err)
			return -1, err
		}
		if (user1 == userID1 && user2 == userID2) || (user1 == userID2 && user2 == userID1) {
			//Return chatid
			return chatID, nil
		}
		if err := rows.Err(); err != nil {
			fmt.Println("Error iterating over rows:", err)
			return -1, err
		}
	}
	return -1, err
}
