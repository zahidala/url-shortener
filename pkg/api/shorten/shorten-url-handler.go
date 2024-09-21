package shorten

import (
	"net/http"
	"url-shortener/pkg/db"
	"url-shortener/pkg/utils"

	UrlVerifier "github.com/davidmytton/url-verifier"
)

func ShortenURLHandler(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")

	if url == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	verifier := UrlVerifier.NewVerifier()
	ret, err := verifier.Verify(url)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !ret.IsURL {
		w.WriteHeader(http.StatusBadRequest)
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
