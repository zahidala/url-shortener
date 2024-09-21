package shorten

import (
	"database/sql"
	"log"
	"net/http"
	"url-shortener/pkg/db"
)

func ShortenGetOriginalURLHandler(w http.ResponseWriter, r *http.Request) {
	shortUrl := r.PathValue("shortUrl")

	if shortUrl == "" {
		http.Error(w, "Short URL is required", http.StatusBadRequest)
		return
	}

	query := `SELECT originalUrl FROM urls WHERE shortUrl = $1`

	var originalUrl string

	shortUrlStmt, shortUrlErr := db.GetDB().Prepare(query)
	if shortUrlErr != nil {
		log.Println(shortUrlErr)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer shortUrlStmt.Close()

	findShortUrlErr := shortUrlStmt.QueryRow(shortUrl).Scan(&originalUrl)
	switch {
	case findShortUrlErr == sql.ErrNoRows:
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	case findShortUrlErr != nil:
		log.Println(findShortUrlErr)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, originalUrl, http.StatusSeeOther)
}
