package main

import (
	"log"
	"net/http"
	"os"
	"url-shortener/pkg/api/shorten"
	"url-shortener/pkg/api/users"
	"url-shortener/pkg/db"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	db.Init()

	port := "8080"

	if len(os.Args) == 2 {
		port = os.Args[1]
	} else if len(os.Args) > 2 {
		log.Fatal("Usage: ./url-shortener [port]")
	}

	http.HandleFunc("POST /login", users.UserLoginHandler)
	http.HandleFunc("GET /logout", users.UserLogoutHandler)
	http.HandleFunc("POST /register", users.UserCreateHandler)

	http.HandleFunc("POST /shorten", shorten.ShortenURLHandler)

	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
