package users

import (
	"log"
	"net/http"
	"url-shortener/pkg/db"
	"url-shortener/pkg/utils"
)

func UserCreateHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	hashedPassword, hashErr := utils.HashPassword(password)
	if hashErr != nil {
		log.Println(hashErr)
		http.Error(w, "Error generating password hash", http.StatusInternalServerError)
		return
	}

	query := "INSERT INTO users (name, username, email, password, profilePicture) VALUES (?, ?, ?, ?, 'https://iili.io/dW44kLG.jpg')"

	userAddExecErr := db.PrepareAndExecute(query,
		name,
		username,
		email,
		hashedPassword,
	)
	if userAddExecErr != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
