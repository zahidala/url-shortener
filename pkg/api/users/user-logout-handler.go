package users

import (
	"net/http"
	"time"
	"url-shortener/pkg/db"
)

func UserLogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, cookieErr := r.Cookie("sessionId")
	if cookieErr != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionId := cookie.Value

	sessionDeleteQuery := "DELETE FROM sessions WHERE id = ?"

	sessionDeleteExecErr := db.PrepareAndExecute(sessionDeleteQuery, sessionId)
	if sessionDeleteExecErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "sessionId",
		Value:   "",
		Expires: time.Now(),
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
