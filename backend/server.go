
package main

import (
	app "backend/pkg/app"
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
	mux.Handle("/images/", http.StripPrefix("/images", http.FileServer(http.Dir("./images"))))
	/* mux.Handle("/avatars/", http.StripPrefix("/avatars", http.FileServer(http.Dir("./avatars"))))
	mux.Handle("/postsImages/", http.StripPrefix("/postsImages", http.FileServer(http.Dir("./postsImages"))))
	mux.Handle("/commentsImages/", http.StripPrefix("/commentsImages", http.FileServer(http.Dir("./commentsImages")))) */
	// auth
	mux.HandleFunc("/api/login", app.LoginHandler)
	mux.HandleFunc("/api/register", app.RegisterHandler)
	mux.HandleFunc("/api/session", app.SessionHandler)
	// profile
	mux.HandleFunc("/api/privacy", func(w http.ResponseWriter, r *http.Request) {
		app.PrivacyHandler(w, r)
		app.SignalChan <- "privacy_updated"
	})
	mux.HandleFunc("/api/user", app.GetUserInfoHandler)
	mux.HandleFunc("/api/editProfile", app.ProfileEditorHandler)
	mux.HandleFunc("/api/uploadImage", app.UpdateImageHandler)

	//notifications
	mux.HandleFunc("/api/notifications", app.NotificationsHandler)
	mux.HandleFunc("/api/markAsSeen", app.NotifMarkAsSeenHandler)

	// posts
	mux.HandleFunc("/api/posts", app.PostsHandler)
	mux.HandleFunc("/api/newPost", app.NewPostHandler)
	mux.HandleFunc("/api/newComment", app.NewCommentHandler)
	mux.HandleFunc("/api/comment", app.CommentHandler)
	// search
	mux.HandleFunc("/api/allusers", app.GetAllUsersHandler)
	// followers
	mux.HandleFunc("/api/followers", func(w http.ResponseWriter, r *http.Request) {
		app.FollowHandler(w, r)
		app.SignalChan <- "followers_updated"
	})
	mux.HandleFunc("/api/messages", app.GetMessages)
	// groups
	mux.HandleFunc("/api/groups", app.GetAllGroupsHandler)
	mux.HandleFunc("/api/newGroup", app.NewGroupHandler)
	mux.HandleFunc("/api/joinGroup", app.JoinGroupHandler)
	mux.HandleFunc("/api/getGroup", app.GetGroupHandler)
	mux.HandleFunc("/api/leaveGroup", app.LeaveGroupHandler)
	mux.HandleFunc("/api/deleteGroup", app.DeleteGroupHandler)
	mux.HandleFunc("/api/newEvent", app.NewEventHandler)
	mux.HandleFunc("/api/events", app.GetEventsHandler)
	mux.HandleFunc("/api/sendRSVP", app.SendRSVPHandler)
	mux.HandleFunc("/api/getGroupMembers", app.GetGroupMembersHandler)
	// websocket
	mux.HandleFunc("/ws", app.WsHandler)
	// Wrap the mux with the CORS middleware
	handlerWithCors := enableCors(mux)
	// Start the server on port 8080
	fmt.Println("SN is running on http://localhost:8080/")
	http.ListenAndServe(":8080", handlerWithCors)
}
