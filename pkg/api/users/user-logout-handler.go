package users

import (
	"log"
	"net/http"
	"time"
	"url-shortener/pkg/db"
)

func UserLogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, cookieErr := r.Cookie("sessionId")
	if cookieErr != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sessionId := cookie.Value

	sessionDeleteQuery := "DELETE FROM sessions WHERE id = $1"

	sessionDeleteExecErr := db.PrepareAndExecute(sessionDeleteQuery, sessionId)
	if sessionDeleteExecErr != nil {
		log.Println(sessionDeleteExecErr)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "sessionId",
		Value:   "",
		Expires: time.Now(),
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
