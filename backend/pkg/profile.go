package pkg

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	db "backend/pkg/db/sqlite"
)

// POST request updates privacy for own profile
func PrivacyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		userID, err := CheckAuth(r)
		if err != nil {
			fmt.Println("PrivacyHandler: Autherror ", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var requestBody struct {
			NewPrivacy bool `json:"newPrivacy"`
		}
		err = decoder.Decode(&requestBody)
		if err != nil {
			fmt.Println("PrivacyHandler: badRequest ", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		err = updatePrivacy(userID, requestBody.NewPrivacy)
		if err != nil {
			fmt.Println("PrivacyHandler: Unable to update privacy ", err)
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}

		jsonResponse, err := json.Marshal(requestBody)
		if err != nil {
			http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
			return
		}
		w.Write(jsonResponse)
	}
}

func ProfileEditorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		userID, err := CheckAuth(r)
		if err != nil {
			fmt.Println("ProfileEditorHandler: Autherror ", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var requestBody struct {
			NewNickName   string `json:"nickname"`
			NewAboutMe    string `json:"aboutMe"`
			NewAvatarPath string `json:"avatar"`
		}
		err = decoder.Decode(&requestBody)
		if err != nil {
			fmt.Println("ProfileEditorHandler: badRequest ", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		fmt.Println(requestBody.NewAvatarPath)

		err = updateProfileData(userID, requestBody.NewAboutMe, requestBody.NewNickName, requestBody.NewAvatarPath)
		if err != nil {
			fmt.Println("ProfileEditorHandler: Unable to update profile ", err)
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}

		jsonResponse, err := json.Marshal(requestBody)
		if err != nil {
			http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
			return
		}
		w.Write(jsonResponse)
	}
}

// Gets User's profile from db and returns relevant fields based on their privacy
// IF PRIVACY == 0, then profile is public
func fetchUserByID(id int) (*User, error) {
	row := db.DB.QueryRow(`SELECT
            id,
            firstname,
            lastname,
            CASE WHEN privacy = 0 THEN date_of_birth ELSE NULL END AS date_of_birth,
            avatar,
            nickname,
            CASE WHEN privacy = 0 THEN about ELSE NULL END AS about,
            privacy
        FROM users
        WHERE id = ?`, id)

	var user User

	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.DateOfBirth, &user.Avatar, &user.NickName, &user.AboutMe, &user.Privacy)
	if err == sql.ErrNoRows {
		return nil, nil // User not found
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

// return user info when clicking on name in profile search
func GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	clientID, err := CheckAuth(r)
	if err != nil {
		http.Error(w, "Session not found: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Get the user ID from query parameters
	userIDStr := r.URL.Query().Get("id")
	if userIDStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Convert user ID to integer
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	user, err := fetchUserByID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	user.Followers, err = GetAllFollowers(userID)
	if err != nil {
		fmt.Println("Error getting followers", err)
	}
	user.Following, err = GetAllFollowing(userID)
	if err != nil {
		fmt.Println("Error getting followers", err)
	}
	user.IsFollowing, err = IsFollowing(clientID, userID)
	if err != nil {
		fmt.Printf("Couldnt retrieve is following, probably doesnt exists")
		user.IsFollowing = false
	}
	fmt.Println("USERINFO ", user)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// change privacy of userID to privacy
func updatePrivacy(userID int, privacy bool) error {
	query := `UPDATE users SET privacy = ? WHERE id = ?`

	result, err := db.DB.Exec(query, privacy, userID)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	fmt.Println(result)

	return nil
}

// Avatar path only removed, image still in directory
func updateProfileData(userID int, aboutMe, nickName, emptyAvatarPath string) error {
	query := `Update users Set about = ?, nickname = ?, avatar = ? Where id = ?`

	_, err := db.DB.Exec(query, aboutMe, nickName, emptyAvatarPath, userID)
	if err != nil {
		return fmt.Errorf("updateProfileData error: %w", err)
	}
	return nil
}

func UpdateImageHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := CheckAuth(r)
	if err != nil {
		fmt.Println("UpdateImageHandler: Autherror ", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Parse the multipart form in the request
	err = r.ParseMultipartForm(550 << 10) // limit your maxMultipartMemory to 550KB
	if err != nil {
		http.Error(w, "Could not parse multipart form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image") // Retrieve the file from form data "image" is the key of the form data
	if err != nil {
		http.Error(w, "Could not get the file", http.StatusBadRequest)
		return
	}

	userString := strconv.Itoa(userID)
	from := r.FormValue("from")
	postID := r.FormValue("postID")
	commentID := r.FormValue("commentID")
	imgPath := ""

	fmt.Println(from, postID, commentID)

	defer file.Close()

	ext := filepath.Ext(header.Filename)

	if from == "changeAvatarImage" {
		imgPath = "./avatars/" + userString + ext
	} else if from == "postImage" {
		imgPath = "./postsImages/" + postID + ext
	} else if from == "commentImage" {
		imgPath = "./commentImages/" + commentID + ext
	} else {
		log.Println("Error: Invalid 'from' value")
		return
	}

	dst, err := os.Create(imgPath) // Overwrite the existing file if it's present
	if err != nil {
		http.Error(w, "Could not create the file", http.StatusInternalServerError)
		return
	}

	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil { // Copy the uploaded file to the destination file
		http.Error(w, "Could not copy the file", http.StatusInternalServerError)
		return
	}

	if from == "changeAvatarImage" {
		query := `Update users Set avatar = ? Where id = ?`
		_, err = db.DB.Exec(query, imgPath, userID)
		if err != nil {
			fmt.Println("imageUpload db update error: %w", err)
		}
	}

	log.Println("Successfully uploaded the image:", imgPath)
}
