package shorten

import (
	"log"
	"net/http"
	"strings"
	"url-shortener/pkg/db"
	"url-shortener/pkg/utils"

	UrlVerifier "github.com/davidmytton/url-verifier"
)

func ShortenURLHandler(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")

	if url == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	// Ensure the URL starts with http:// or https://
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	verifier := UrlVerifier.NewVerifier()
	ret, err := verifier.Verify(url)

	if err != nil {
		log.Println(err)
		http.Error(w, "Error verifying URL", http.StatusBadRequest)
		return
	}

	if !ret.IsURL {
		log.Println("Not a valid URL")
		http.Error(w, "Not a valid URL", http.StatusBadRequest)
		return
	}

	shortenURL := utils.GenerateShortenedURL()

	user := utils.GetUserInfoBySession(w, r)

	query := "INSERT INTO urls (shortUrl, originalUrl, userId) VALUES ($1, $2, $3)"

	shortUrlAddExecErr := db.PrepareAndExecute(query,
		shortenURL,
		url,
		user.ID,
	)
	if shortUrlAddExecErr != nil {
		log.Println(shortUrlAddExecErr)
		http.Error(w, "Error saving shortened URL", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
