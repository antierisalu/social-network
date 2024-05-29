package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
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
		fmt.Println("su password ja user on: ", cred)
	}
}
