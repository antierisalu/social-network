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
	"time"

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

// Gets User's profile from db and returns all info
// IF PRIVACY == 0, then profile is public
func fetchUserByID(id int) (*User, error) {
	row := db.DB.QueryRow(`SELECT
            id,
            firstname,
            lastname,
            date_of_birth,
            avatar,
            nickname,
            about,
            privacy,
			email
        FROM users
        WHERE id = ?`, id)

	var user User

	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.DateOfBirth, &user.Avatar, &user.NickName, &user.AboutMe, &user.Privacy, &user.Email)
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
		fmt.Println("Error getting followers", err)
	}

	user.AreFollowing, err = IsFollowing(userID, clientID)
	if err != nil {
		fmt.Println("Couldnt retrieve relationship data, something went wrong in userhandler")
	}
	user.IsFollowing, err = IsFollowing(clientID, userID)
	if err != nil {
		fmt.Println("Couldnt retrieve relationship data, something went wrong in userhandler")
	}

	/* // this means that the user has requested
	fmt.Println(hasRelationship, user.IsFollowing)
	if hasRelationship && !user.IsFollowing {
		user.HasRequested = true
		fmt.Println(user.FirstName, clientID)

	} */

	// if you shouldnt be able to see the profile, clear About me and date of birth

	if user.Privacy == 1 && user.AreFollowing < 1 && clientID != userID {
		fmt.Println("NOT SEEING THE ABOUTME LOL")
		user.AboutMe = sql.NullString{}
		user.DateOfBirth = sql.NullString{}
	}

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

	user.Posts, err = GetPostsForProfile(userID, clientID)
	if err != nil {
		fmt.Println("GetuserInfo: error with getPostsForProfile")
	}

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
		imgPath = "./commentsImages/" + commentID + ext
	} else {
		log.Println("Error: Invalid 'from' value")
		return
	}

	dst, err := os.Create(imgPath) // Overwrite the existing file if it's present
	if err != nil {
		fmt.Println("imageUpload os.Create error: %w", err)
		http.Error(w, "Could not create the file", http.StatusInternalServerError)

		return
	}

	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil { // Copy the uploaded file to the destination file
		fmt.Println("imageUpload io.Copy error: %w", err)
		http.Error(w, "Could not copy the file", http.StatusInternalServerError)
		return
	}

	switch from {
	case "changeAvatarImage":
		query := `Update users Set avatar = ? Where id = ?`
		_, err = db.DB.Exec(query, imgPath, userID)
		fmt.Println(imgPath, userID, " IMGPATH AND USERID")

		if err != nil {
			fmt.Println("imageUpload avatar db update error: %w", err)
		}
	case "postImage":
		query := `Update posts Set media = ? Where id = ?`
		_, err = db.DB.Exec(query, imgPath, postID)
		if err != nil {
			fmt.Println("imageUpload post db update error: %w", err)
		}
	case "commentImage":
		query := `Update comments Set media = ? Where id = ?`
		_, err = db.DB.Exec(query, imgPath, commentID)
		if err != nil {
			fmt.Println("imageUpload comment db update error: %w", err)
		}
	default:
		log.Println("Error: Invalid 'from' value")
		return
	}

	w.Write([]byte(imgPath))
}

func GetPostsForProfile(userID, clientID int) ([]Post, error) {
	query := `SELECT DISTINCT p.id, p.user_id, media, content, p.created_at
				FROM posts p
				LEFT JOIN followers ON followers.user_id = p.user_id
				LEFT JOIN post_custom_privacy on p.id = post_custom_privacy.post_id
				WHERE (p.user_id = ? AND p.privacy = 0)
				OR (p.privacy = 1 AND p.user_id = ? AND followers.follower_id = ? AND followers.isFollowing = 1)
				OR (p.privacy = 2 AND p.user_id = ? AND post_custom_privacy.user_id = ? AND followers.isFollowing = 1)
				OR p.user_id = ? AND p.user_id = ?
				ORDER BY p.created_at DESC;`

	postRows, err := db.DB.Query(query, userID, userID, clientID, userID, clientID, clientID, userID)
	if err != nil {
		return nil, err
	}
	defer postRows.Close()

	var posts []Post
	for postRows.Next() {
		var post Post
		err = postRows.Scan(&post.ID, &post.UserID, &post.Img, &post.Content, &post.CreatedAt)
		if err != nil {
			fmt.Println("GetPostsForProfile: error scan post: ", post.ID)
			continue
		}

		commentQuery := `select c.id, c.user_id, c.post_id, c.content, c.media, c.created_at,
			u.FirstName, u.LastName, u.Avatar from comments c
						left join users u
						on c.user_id = u.id
						where post_id = ?`
		commentRows, err := db.DB.Query(commentQuery, post.ID)
		if err != nil {
			fmt.Println("GetPostsForProfile: error querying post comments: ", post.ID, err)
			continue
		}

		for commentRows.Next() {
			var comment Comment
			err = commentRows.Scan(&comment.ID, &comment.UserID, &comment.PostID, &comment.Content, &comment.Img, &comment.CreatedAt,
				&comment.User.FirstName, &comment.User.LastName, &comment.User.Avatar)
			if err != nil {
				fmt.Println("GetPostsForProfile: error querying comment: ", comment.ID, err)
				continue
			}

			parsedTime, err := time.Parse(time.RFC3339, comment.CreatedAt)
			if err != nil {
				log.Println("GetPrivatePosts, parsedTime section:", err)
				continue
			}

			comment.CreatedAt = parsedTime.Format("2006-01-02 15:04:05")

			post.Comments = append(post.Comments, comment)
		}

		parsedTime, err := time.Parse(time.RFC3339, post.CreatedAt)
		if err != nil {
			log.Println("GetPrivatePosts, parsedTime section:", err)
			continue
		}

		post.CreatedAt = parsedTime.Format("2006-01-02 15:04:05")

		posts = append(posts, post)
	}

	return posts, nil
}
