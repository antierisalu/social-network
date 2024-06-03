package pkg

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	db "backend/pkg/db/sqlite"

	"github.com/gofrs/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// KOIK LOGIN/REGISTER/AUTH HANDLERID SIIA

// logs user in if credentials are valid
func LoginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var cred Credentials
		err := decoder.Decode(&cred)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		if cred.Password == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		valid, err := validateLogin(cred.Email, cred.Password)
		if err != nil || !valid {
			fmt.Println("LOGINHANDLER:", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := GenerateToken()
		session := Session{
			Token:   token,
			Expires: 24,
		}
		err = updateToken(token, cred.Email)
		if err != nil {
			http.Error(w, "Error updating token", http.StatusInternalServerError)
			return
		}

		jsonResponse, err := json.Marshal(session)
		if err != nil {
			http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
			return
		}
		w.Write(jsonResponse)

	}
}

// inserts new user into database if passwords match and email is not taken
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var userData UserData
		err := decoder.Decode(&userData)
		if err != nil {
			fmt.Println(userData, "bad request")
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		fmt.Println("RegisterHandler data recieved:\n", userData)
		token := GenerateToken()

		err = InsertUser(userData, token)
		if err != nil {
			fmt.Println("REGISTERHANDLER: ", err)

			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		session := Session{
			Token:   token,
			Expires: 24,
		}
		jsonResponse, err := json.Marshal(session)
		if err != nil {
			http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
			return
		}
		w.Write(jsonResponse)
	}
}

// automatically logs user in and returns user data if session token is valid
func SessionHandler(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("sessionToken")
	if err != nil {
		http.Error(w, "Unauthorized1", http.StatusUnauthorized)
		return
	}

	//check if cookie exists in database and return user data
	user, err := checkAuth(cookie.Value)
	if err != nil {
		http.Error(w, "Unauthorized2", http.StatusUnauthorized)
	}

	//send userdata to client
	jsonResponse, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)
}

// HELPER FUNCTIONS

// validateLogin validates the login credentials of a user.
//
// It takes in two parameters: email (string) and password (string).
// It returns a boolean value indicating whether the login is valid or not,
// and an error if any occurred during the validation process.
func validateLogin(email, password string) (bool, error) {

	stmt := "SELECT hash FROM users WHERE LOWER(email) = LOWER(?)"

	var hash string
	err := db.DB.QueryRow(stmt, email).Scan(&hash)
	if err != nil {
		/* 		if err == sql.ErrNoRows {
			return false, nil
		} */
		fmt.Println("VALIDATE_USER ERROR: ", err, string(hash))
		return false, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}

// InsertUser inserts a new user into the database with the given user data and token.
//
// Parameters:
// - userData: The user data containing the user's email, password, first name, last name, date of birth, avatar, nickname, and about me.
// - token: The session token for the user.
//
// Returns:
// - err: An error if there was a problem inserting the user into the database.
func InsertUser(userData UserData, token string) (err error) {

	var count int
	err = db.DB.QueryRow("SELECT id FROM users WHERE LOWER(email) = LOWER(?)", userData.Email).Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println("Error checking for existing user in InsertUser:", err)
		return errors.New("error checking email")
	}
	if count > 0 {
		return errors.New("email already exists")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(userData.Password), 12)
	stmt, err := db.DB.Prepare("INSERT INTO users (email, hash, firstname, lastname, date_of_birth, avatar, nickname, about, session) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println("Error preparing statement in InsertUser:", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(userData.Email, hash, userData.FirstName, userData.LastName, userData.DateOfBirth, userData.Avatar, userData.NickName, userData.AboutMe, token)
	if err != nil {
		fmt.Println("Error executing statement in InsertUser:", err)
		return err
	}
	return
}

func GenerateToken() string {
	newUUID, err := uuid.NewV4()
	if err != nil {
		fmt.Println(err)
		return fmt.Sprintf("GenerateToken ERROR: %s", err)
	}
	return newUUID.String()
}

func TokenExists(token string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE session = ?`
	var count int
	err := db.DB.QueryRow(query, token).Scan(&count)
	if err != nil {
		fmt.Println("TokenExists: error getting count:", err)
		return false, err
	}
	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func checkAuth(token string) (User, error) {
	user := User{}
	err := db.DB.QueryRow("SELECT id, email, firstname, lastname, date_of_birth, avatar, nickname, about, session FROM users WHERE session = ?", token).Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.DateOfBirth, &user.Avatar, &user.NickName, &user.AboutMe, &user.Session)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func updateToken(token, email string) error {
	stmt, err := db.DB.Prepare("UPDATE users SET session = ? WHERE email = ?")
	if err != nil {
		fmt.Println("Error preparing statement in UpdateToken:", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(token, email)
	if err != nil {
		fmt.Println("Error executing statement in UpdateToken:", err)
		return err
	}
	return nil
}
