/*
	package main

import (

	"fmt"
	"net/http"

	"backend/pkg"
	db "backend/pkg/db/sqlite"

)

	func main() {
		// Init Database
		db.DatabaseInit()
		defer db.DB.Close()

		// Serve static files from the current directory
		http.Handle("/", http.FileServer(http.Dir("../public")))
		http.Handle("/avatars/", http.StripPrefix("/avatars", http.FileServer(http.Dir("./avatars"))))
		http.Handle("/postsImages/", http.StripPrefix("/postsImages", http.FileServer(http.Dir("./postsImages"))))
		http.Handle("/commentsImages/", http.StripPrefix("/commentsImages", http.FileServer(http.Dir("./commentsImages"))))

		//auth
		http.HandleFunc("/login", pkg.LoginHandler)
		http.HandleFunc("/register", pkg.RegisterHandler)
		http.HandleFunc("/session", pkg.SessionHandler)

		//profile
		http.HandleFunc("/privacy", func(w http.ResponseWriter, r *http.Request) {
			pkg.PrivacyHandler(w, r)
			pkg.SignalChan <- "privacy_updated"
		})
		http.HandleFunc("/user", pkg.GetUserInfoHandler)
		http.HandleFunc("/editProfile", pkg.ProfileEditorHandler)
		http.HandleFunc("/uploadImage", pkg.UpdateImageHandler)

		//notifications
		http.HandleFunc("/notifications", pkg.NotificationsHandler)
		http.HandleFunc("/markAsSeen", pkg.NotifMarkAsSeenHandler)

		//posts
		http.HandleFunc("/posts", pkg.PostsHandler)
		http.HandleFunc("/newPost", pkg.NewPostHandler)
		http.HandleFunc("/newComment", pkg.NewCommentHandler)
		http.HandleFunc("/comment", pkg.CommentHandler)

		//search
		http.HandleFunc("/allusers", pkg.GetAllUsersHandler)

		//followers
		http.HandleFunc("/api/followers", func(w http.ResponseWriter, r *http.Request) {
			pkg.FollowHandler(w, r)
			pkg.SignalChan <- "followers_updated"
		})
		http.HandleFunc("/messages", pkg.GetMessages)

		//groups
		http.HandleFunc("/groups", pkg.GetAllGroupsHandler)
		http.HandleFunc("/newGroup", pkg.NewGroupHandler)
		http.HandleFunc("/joinGroup", pkg.JoinGroupHandler)
		http.HandleFunc("/getGroup", pkg.GetGroupHandler)
		http.HandleFunc("/leaveGroup", pkg.LeaveGroupHandler)
		http.HandleFunc("/deleteGroup", pkg.DeleteGroupHandler)

		http.HandleFunc("/newEvent", pkg.NewEventHandler)
		http.HandleFunc("/events", pkg.GetEventsHandler)
		http.HandleFunc("/sendRSVP", pkg.SendRSVPHandler)

		// websocket
		http.HandleFunc("/ws", pkg.WsHandler)

		// Start the server on port 8080
		fmt.Println("SN is running on http://localhost:8080/")
		http.ListenAndServe(":8080", nil)
	}
*/
package main

import (
	"backend/pkg"
	db "backend/pkg/db/sqlite"
	"fmt"
	"net/http"
)

// Middleware to enable CORS
func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func main() {
	// Init Database
	db.DatabaseInit()
	defer db.DB.Close()
	// Create a new ServeMux
	mux := http.NewServeMux()
	// Serve static files from the current directory
	mux.Handle("/avatars/", http.StripPrefix("/avatars", http.FileServer(http.Dir("./avatars"))))
	mux.Handle("/postsImages/", http.StripPrefix("/postsImages", http.FileServer(http.Dir("./postsImages"))))
	mux.Handle("/commentsImages/", http.StripPrefix("/commentsImages", http.FileServer(http.Dir("./commentsImages"))))
	// auth
	mux.HandleFunc("/api/login", pkg.LoginHandler)
	mux.HandleFunc("/api/register", pkg.RegisterHandler)
	mux.HandleFunc("/api/session", pkg.SessionHandler)
	// profile
	mux.HandleFunc("/api/privacy", func(w http.ResponseWriter, r *http.Request) {
		pkg.PrivacyHandler(w, r)
		pkg.SignalChan <- "privacy_updated"
	})
	mux.HandleFunc("/api/user", pkg.GetUserInfoHandler)
	mux.HandleFunc("/api/editProfile", pkg.ProfileEditorHandler)
	mux.HandleFunc("/api/uploadImage", pkg.UpdateImageHandler)

	//notifications
	mux.HandleFunc("/api/notifications", pkg.NotificationsHandler)
	mux.HandleFunc("/api/markAsSeen", pkg.NotifMarkAsSeenHandler)

	// posts
	mux.HandleFunc("/api/posts", pkg.PostsHandler)
	mux.HandleFunc("/api/newPost", pkg.NewPostHandler)
	mux.HandleFunc("/api/newComment", pkg.NewCommentHandler)
	mux.HandleFunc("/api/comment", pkg.CommentHandler)
	// search
	mux.HandleFunc("/api/allusers", pkg.GetAllUsersHandler)
	// followers
	mux.HandleFunc("/api/followers", func(w http.ResponseWriter, r *http.Request) {
		pkg.FollowHandler(w, r)
		pkg.SignalChan <- "followers_updated"
	})
	mux.HandleFunc("/api/messages", pkg.GetMessages)
	// groups
	mux.HandleFunc("/api/groups", pkg.GetAllGroupsHandler)
	mux.HandleFunc("/api/newGroup", pkg.NewGroupHandler)
	mux.HandleFunc("/api/joinGroup", pkg.JoinGroupHandler)
	mux.HandleFunc("/api/getGroup", pkg.GetGroupHandler)
	mux.HandleFunc("/api/leaveGroup", pkg.LeaveGroupHandler)
	mux.HandleFunc("/api/deleteGroup", pkg.DeleteGroupHandler)
	mux.HandleFunc("/api/newEvent", pkg.NewEventHandler)
	mux.HandleFunc("/api/events", pkg.GetEventsHandler)
	mux.HandleFunc("/api/sendRSVP", pkg.SendRSVPHandler)
	// websocket
	mux.HandleFunc("/ws", pkg.WsHandler)
	// Wrap the mux with the CORS middleware
	handlerWithCors := enableCors(mux)
	// Start the server on port 8080
	fmt.Println("SN is running on http://localhost:8080/")
	http.ListenAndServe(":8080", handlerWithCors)
}
