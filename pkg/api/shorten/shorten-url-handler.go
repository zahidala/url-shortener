package shorten

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"url-shortener/pkg/db"
	Types "url-shortener/pkg/types"
	"url-shortener/pkg/utils"

	UrlVerifier "github.com/davidmytton/url-verifier"
)

type ShortenURLBody struct {
	URL string `json:"url"`
}

func ShortenURLHandler(w http.ResponseWriter, r *http.Request) {
	var body ShortenURLBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println("Error decoding body:", err)
		return
	}

	if body.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	// Ensure the URL starts with http:// or https://
	if !strings.HasPrefix(body.URL, "http://") && !strings.HasPrefix(body.URL, "https://") {
		body.URL = "http://" + body.URL
	}

	verifier := UrlVerifier.NewVerifier()
	ret, err := verifier.Verify(body.URL)

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

	user := utils.GetUserInfoBySession(r)

	var query string
	var params []interface{}

	if (user != Types.User{}) {
		query = "INSERT INTO urls (shortUrl, originalUrl, userId) VALUES ($1, $2, $3)"
		params = []interface{}{shortenURL, body.URL, user.ID}
	} else {
		query = "INSERT INTO urls (shortUrl, originalUrl) VALUES ($1, $2)"
		params = []interface{}{shortenURL, body.URL}
	}

	shortUrlAddExecErr := db.PrepareAndExecute(query,
		params...,
	)
	if shortUrlAddExecErr != nil {
		log.Println(shortUrlAddExecErr)
		http.Error(w, "Error saving shortened URL", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": map[string]string{
			"shortUrl": shortenURL,
		},
	})
}
