package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	db "backend/pkg/db/sqlite"
)

func GetMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("getMessages, error getting req body:", err)
			return
		}
		var msgGet MessageGetter
		err = json.Unmarshal(body, &msgGet)
		if err != nil {
			fmt.Println("getMessages error unmarshaling,", err)
		}

		messages := GetTenMessages(msgGet.Date, msgGet.ID, msgGet.ChatID)

		jsonResponse, err := json.Marshal(messages)
		if err != nil {
			http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(jsonResponse)
		if err != nil {
			http.Error(w, "Failed to send response", http.StatusInternalServerError)
			return
		}
	}
}

func GetTenMessages(date time.Time, msgid, chatid int) []ChatMessage {
	var messages []ChatMessage

	// fmt.Printf("Getting ten messages for chatid: %v | last msgid:%v\n ", chatid, msgid)
	var query string
	if msgid == 0 { // initial load
		query = `SELECT id, content, user_id, created_at FROM chatmessages WHERE chat_id = ? and id > ? ORDER BY id DESC LIMIT 10;`
	} else { // subsequent loads
		query = `SELECT id, content, user_id, created_at FROM chatmessages WHERE chat_id = ? and id < ? ORDER BY id DESC LIMIT 10;`
	}

	rows, err := db.DB.Query(query, chatid, msgid)
	if err != nil {
		fmt.Println("GetTenMessages: error querying db: ", err)
		return []ChatMessage{}
	}

	for rows.Next() {
		var msg ChatMessage
		rows.Scan(&msg.ID, &msg.Content, &msg.User, &msg.Date)
		msg.SetUsername(db.DB)
		messages = append(messages, msg)
	}
	return messages
}

// // Inserts a private message to database 'chatmessages' and returns the createdAt, message_ID, nil on success
// // On error returns "ERROR", -1, err
func InsertPrivateMessage(userID, chatID int, message string, isGroup bool) (string, int, error) {
	stmt, err := db.DB.Prepare("INSERT INTO chatmessages (user_id, chat_id, content, is_group, created_at, seen) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println("Error preparing statement in InsertPrivateMessage:", err)
		return "ERROR", -1, err
	}
	defer stmt.Close()
	//https://pkg.go.dev/database/sql#Result
	now := time.Now()
	seen := false
	result, err := stmt.Exec(userID, chatID, message, isGroup, now, seen)
	if err != nil {
		fmt.Println("Error executing statement in InsertPrivateMessage:", err)
		return "ERROR", -1, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Error getting last message_ID InsertPrivateMessage:", err)
		return "ERROR", -1, err
	}
	// Update lastmessage date
	stmt, err = db.DB.Prepare(`UPDATE user_chats SET last_message = ? WHERE id = ?;`)
	if err != nil {
		fmt.Println("InsertPrivateMessage: Error Inserting LastMessage:", err)
	}
	_, err = stmt.Exec(now, chatID)
	if err != nil {
		fmt.Println("Error executing LastMessage statement in InsertPrivateMessage:", err)
	}

	return now.Format("2006-01-02 15:04:05.999999-07:00"), int(lastID), nil
}

// Creates a new database chat entry for user1 & user2
func InsertNewChat(user1, user2 int) error {
	stmt, err := db.DB.Prepare("INSERT INTO user_chats (user1, user2, last_message, created_at) VALUES (?, ?, ?, ?)")
	if err != nil {
		fmt.Println("Error preparing statement in InsertNewChat:", err)
		return err
	}
	defer stmt.Close()
	if user1 > user2 {
		user1, user2 = user2, user1
	}
	// Empty string for last_message
	_, err = stmt.Exec(user1, user2, "", time.Now())
	if err != nil {
		fmt.Println("Error executing statement in InsertNewChat:", err)
		return err
	}
	return nil
}

// Takes in user1 ID and user2ID to check if there is a database entry for their chat,
// IF entry exists returns the chatID if not returns -1, error
func GetChatID(userID1, userID2 int) (int, error) {
	stmt := "SELECT * FROM user_chats"
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

// Get Email From ID
func GetEmailFromID(id int) (string, error) {
	stmt := "SELECT email FROM users WHERE id = ?"
	var email string
	err := db.DB.QueryRow(stmt, id).Scan(&email)
	if err != nil {
		return "", err
	}
	return email, nil
}

// Get ID From Email
func GetIDFromEmail(email string) (int, error) {
	stmt := "SELECT id FROM users WHERE email = ?"
	var ID int
	err := db.DB.QueryRow(stmt, email).Scan(&ID)
	if err != nil {
		return -1, err
	}
	return ID, nil
}

// Generate last message map (store) for clientID
func GetLastMessageStore(clientID int) (map[int]string, error) {
	lastMsgMap := make(map[int]string)

	// Fetch all rows where clientID is included
	stmt := "SELECT CASE WHEN user1 = ? THEN user2 ELSE user1 END AS other_user, last_message FROM user_chats WHERE user1 = ? OR user2 = ?"

	rows, err := db.DB.Query(stmt, clientID, clientID, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Populate the map
	for rows.Next() {
		var otherUserID int
		var lastMessage string
		if err := rows.Scan(&otherUserID, &lastMessage); err != nil {
			return nil, err
		}
		lastMsgMap[otherUserID] = lastMessage
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return lastMsgMap, nil
}
