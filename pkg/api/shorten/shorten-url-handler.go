package shorten

import (
	"net/http"
)

func ShortenURLHandler(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")

	if url == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
