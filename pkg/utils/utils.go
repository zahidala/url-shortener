package utils

import (
	"log"
	"math/rand"
	"net/http"
	"time"
	"url-shortener/pkg/db"
	Types "url-shortener/pkg/types"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

// GenerateSessionID generates a new UUID for session ID
func GenerateSessionID() string {
	uuid, _ := uuid.NewV4()
	return uuid.String()
}

// GenerateShortenedURL generates a random string of specified length
func GenerateShortenedURL() string {
	var (
		randomChars  = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321")
		stringLength = 8 // Length of the shortened URL
	)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	shortURL := make([]rune, stringLength)
	for i := range shortURL {
		shortURL[i] = randomChars[r.Intn(len(randomChars))]
	}
	return string(shortURL)
}

// HashPassword generates a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// CompareHashAndPassword compares a bcrypt hashed password with its possible plaintext equivalent.
// Returns nil on success, or an error on failure.
func CompareHashAndPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// Returns the user object based on the session ID cookie
func GetUserInfoBySession(w http.ResponseWriter, r *http.Request) Types.User {
	cookie, cookieErr := r.Cookie("sessionId")

	if cookieErr != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return Types.User{}
	}

	sessionId := cookie.Value

	sessionsQuery := "SELECT userId FROM sessions WHERE id = ?"

	sessionStmt, sessionErr := db.GetDB().Prepare(sessionsQuery)
	if sessionErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(sessionErr)
		return Types.User{}
	}

	var userId int

	sessionRowErr := sessionStmt.QueryRow(sessionId).Scan(&userId)

	if sessionRowErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return Types.User{}
	}

	userQuery := "SELECT * FROM users WHERE id = ?"

	userStmt, userErr := db.GetDB().Prepare(userQuery)

	if userErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(userErr)
		return Types.User{}
	}

	var user Types.User

	userRowErr := userStmt.QueryRow(userId).Scan(&user.ID, &user.Name, &user.Username, &user.Email, &user.Password, &user.ProfilePicture)

	if userRowErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(userRowErr)
		return Types.User{}
	}

	return user
}
