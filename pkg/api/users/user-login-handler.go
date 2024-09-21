package users

import (
	"database/sql"
	"net/http"
	"strings"
	"time"
	"url-shortener/pkg/db"
	"url-shortener/pkg/utils"
)

func UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimSpace(r.FormValue("username"))
	password := r.FormValue("password")

	if len(username) < 3 || len(username) > 20 || len(password) < 8 || len(password) > 128 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userQuery := "SELECT id, password FROM users WHERE username = $1"

	userStmt, userErr := db.GetDB().Prepare(userQuery)
	if userErr != nil {
		http.Error(w, "Error preparing query", http.StatusInternalServerError)
		return
	}
	defer userStmt.Close()

	var userId int
	var hashedPassword string

	findUserErr := userStmt.QueryRow(username).Scan(&userId, &hashedPassword)
	switch {
	case findUserErr == sql.ErrNoRows:
		w.WriteHeader(http.StatusNotFound)
		return
	case findUserErr != nil:
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}

	compareErr := utils.CompareHashAndPassword(hashedPassword, password)
	if compareErr != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Look for existing sessions and delete them
	sessionDeleteQuery := "DELETE FROM sessions WHERE userId = $1"

	sessionDeleteExecErr := db.PrepareAndExecute(sessionDeleteQuery, userId)
	if sessionDeleteExecErr != nil {
		http.Error(w, "Error deleting session", http.StatusInternalServerError)
		return
	}

	sessionId := utils.GenerateSessionID()

	createdAt := time.Now()
	expiresAt := createdAt.Add(24 * time.Hour)

	http.SetCookie(w, &http.Cookie{
		Name:    "sessionId",
		Value:   sessionId,
		Expires: expiresAt,
	})

	sessionsAddQuery := "INSERT INTO sessions (id, userId, createdAt, expiresAt) VALUES ($1, $2, $3, $4)"

	sessionAddExecErr := db.PrepareAndExecute(sessionsAddQuery, sessionId, userId, createdAt, expiresAt)
	if sessionAddExecErr != nil {
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
