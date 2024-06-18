package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	db "backend/pkg/db/sqlite"
)

// // Handler

func GetMessages(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET MESSAGE RECIEVED")
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

		fmt.Println(msgGet)
		messages := GetTenMessages(msgGet.Date, msgGet.ID, msgGet.ChatID)
		fmt.Println("got these messages:")
		for _, v := range messages {
			fmt.Println(v.ID, v.User, v.Content, v.Date)
			fmt.Println(v.Username, ": ", v.Content)
		}
	}
	// 	jsonResponse, err := json.Marshal(messages)
	// 	if err != nil {
	// 		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
	// 		return
	// 	}
	// 	w.Header().Set("Content-Type", "application/json")
	// 	_, err = w.Write(jsonResponse)
	// 	if err != nil {
	// 		http.Error(w, "Failed to send response", http.StatusInternalServerError)
	// 		return
	// 	}
	// }
}

func GetTenMessages(date time.Time, msgid, chatid int) []ChatMessage {
	var messages []ChatMessage
	// query := `SELECT content, nickname, chatmessages.created_at, message_id FROM chatmessages
	// JOIN users ON chatmessages.user_id = users.user_id
	// WHERE chat_id = ? AND chatmessages.created_at < ?
	// AND chatmessages.message_id <> ?
	// ORDER BY chatmessages.created_at DESC
	// LIMIT 10;`
	//this long query gets content, create_at and message_id from the chatmessages table and the nickname from the joined users table.
	//it filters messages for only a specific chat_id and messages created before a certain date.
	// it also excludes the messageid of the last previously loaded message so certain messages dont get loaded multiple times.
	query := `SELECT id, content, user_id, created_at FROM chatmessages WHERE chat_id = ? and id > ? ORDER BY id DESC LIMIT 10;`

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
	stmt, err := db.DB.Prepare("INSERT INTO chatmessages (user_id, chat_id, content, is_group, created_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println("Error preparing statement in InsertPrivateMessage:", err)
		return "ERROR", -1, err
	}
	defer stmt.Close()
	//https://pkg.go.dev/database/sql#Result
	now := time.Now()
	result, err := stmt.Exec(userID, chatID, message, isGroup, now)
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
