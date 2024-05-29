package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	db "backend/pkg/db/sqlite"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// KOIK LOGIN/REGISTER/AUTH HANDLERID SIIA
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var cred Credentials
		err := decoder.Decode(&cred)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		err = storeNewUser(cred.Username, cred.Password)
		if err != nil {
			fmt.Println("LOGINHANDLER:", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		fmt.Println("su password ja user on: ", cred)

	}
}

func validateUser(email, password string) (bool, error) {
	stmt := "SELECT * FROM users WHERE LOWER(email) = LOWER(?) AND password = (?)"

	var cred Credentials
	err := db.DB.QueryRow(stmt, email, password).Scan(&cred)
	if err != nil {
		fmt.Println("VALIDATE_USER: ", err)
		return false, err
	}
	return true, nil
}

func storeNewUser(username, password string) error {
	stmt, err := db.DB.Prepare("INSERT INTO users (username, password) VALUES (?, ?)")
	if err != nil {
		fmt.Println("Error preparing statement in storeNewUser1:", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, password)
	if err != nil {
		fmt.Println("Error executing statement in storeNewUser2:", err)
		return err
	}

	return nil
}

// stolen from RTF not usable yet
func InsertUser(nickname, firstname, lastname, password, email string, age int) (err error) {
	// turn password into hash
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	var count int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE LOWER(nickname) = LOWER(?)", nickname).Scan(&count)
	if err != nil {
		fmt.Println("Error checking for existing user in InsertUser:", err)
		return errors.New("error checking username")
	}
	if count > 0 {
		return errors.New("username already exists")
	}

	err = db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE LOWER(email) = LOWER(?)", email).Scan(&count)
	if err != nil {
		fmt.Println("Error checking for existing user in InsertUser:", err)
		return errors.New("error checking email")
	}
	if count > 0 {
		return errors.New("email already exists")
	}

	stmt, err := db.DB.Prepare("INSERT INTO users (nickname, hash, email, age, firstname, lastname, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println("Error preparing statement in InsertUser:", err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(nickname, []byte(hash), email, age, firstname, lastname, time.Now())
	if err != nil {
		fmt.Println("Error executing statement in InsertUser:", err)
		return err
	}
	return
}
