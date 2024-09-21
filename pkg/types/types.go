package types

import (
	"database/sql"
	"sync"
	"time"
)

type Database struct {
	Conn *sql.DB
	Mu   sync.Mutex
}

type User struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	ProfilePicture string `json:"profilePicture"`
}

type Session struct {
	ID        string    `json:"id"`
	UserID    int       `json:"userId"`
	Data      string    `json:"data"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type URL struct {
	ID          int       `json:"id"`
	ShortURL    string    `json:"shortUrl"`
	OriginalURL string    `json:"originalUrl"`
	UserID      int       `json:"userId"`
	CreatedAt   time.Time `json:"createdAt"`
	ExpiresAt   time.Time `json:"expiresAt"`
	IsActive    bool      `json:"isActive"`
}
