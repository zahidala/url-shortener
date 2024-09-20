package db

import (
	"database/sql"
	"log"
	"os"
	"sync"
	Types "url-shortener/pkg/types"

	_ "github.com/mattn/go-sqlite3"
)

var instance *Types.Database
var once sync.Once

// Init initializes the database connection
func Init() {
	once.Do(func() {
		dbFile := "./database.sqlite3"
		_, err := os.Stat(dbFile)
		dbExists := !os.IsNotExist(err)

		conn, err := sql.Open("sqlite3", "./database.sqlite3")
		if err != nil {
			log.Fatalf("Error opening the database: %s", err)
			return
		}

		if err := conn.Ping(); err != nil {
			log.Fatalf("Error connecting to the database: %s", err)
			return
		}

		if dbExists {
			log.Println("Connected to the database")
		}

		instance = &Types.Database{
			Conn: conn,
		}

		if !dbExists {
			log.Println("Database file does not exist. Creating a new file and seeding database...")
			seedDB()
		}
	})
}

// seedDB seeds the newly created database file with initial data if it does not exist
func seedDB() {
	createTables := []string{
		`CREATE TABLE Users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			username TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			profilePicture TEXT
		);`,

		`CREATE TABLE Sessions (
			id TEXT PRIMARY KEY,
			userId INTEGER NOT NULL,
			data TEXT,
			createdAt DATETIME NOT NULL,
			expiresAt DATETIME NOT NULL,
			FOREIGN KEY (userId) REFERENCES Users(id)
		);`,

		`CREATE TABLE Urls (
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			shortUrl TEXT NOT NULL, 
			originalUrl TEXT NOT NULL, 
			createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, 
			expiresAt DATETIME, 
			clicks INTEGER DEFAULT 0, 
			userId INTEGER, 
			isActive BOOLEAN DEFAULT TRUE
		);`,
	}

	// Create tables
	for _, query := range createTables {
		if err := PrepareAndExecute(query); err != nil {
			log.Println(query)
			log.Fatalf("Error creating table: %s", err)
			return
		}
	}

	initialData := []string{
		`INSERT INTO Users (name, username, email, password, profilePicture) VALUES ('John Doe', 'johndoe', 'johndoe@gmail.com', '$2a$10$M9APgO1pJZgsfMdj9SmZEORF94WYnS5RkXrIaVA7ZG6bXgzSB5lEa', 'https://iili.io/dW44kLG.jpg');`,
		`INSERT INTO Users (name, username, email, password, profilePicture) VALUES ('Jane Doe', 'janedoe', 'janedoe@gmail.com', '$2a$10$M9APgO1pJZgsfMdj9SmZEORF94WYnS5RkXrIaVA7ZG6bXgzSB5lEa', 'https://iili.io/dW44kLG.jpg');`,

		`INSERT INTO urls (shortUrl, originalUrl, userId) VALUES ('http://localhost:8080/abc123', 'https://www.google.com', 1);`,
		`INSERT INTO urls (shortUrl, originalUrl, userId) VALUES ('http://localhost:8080/xyz123', 'https://www.facebook.com', 2);`,
	}

	// // Insert initial data
	for _, query := range initialData {
		if err := PrepareAndExecute(query); err != nil {
			log.Println(query)
			log.Fatalf("Error inserting initial data: %s", err)
			return
		}
	}

	log.Println("Database seeded successfully")
	log.Println("Connected to the database")
}

// GetDB returns the database instance
func GetDB() *sql.DB {
	if instance == nil {
		log.Fatal("Database not initialized. Call Init() first.")
	}
	return instance.Conn
}

// CloseDB closes the database connection
func CloseDB() {
	if instance != nil {
		instance.Mu.Lock()
		defer instance.Mu.Unlock()
		if err := instance.Conn.Close(); err != nil {
			log.Printf("Error closing the database: %s", err)
		}
	}
}

// PrepareAndExecute prepares and executes a query. It returns an error if the query fails.
// May be expanded to return the result of the query in the future.
func PrepareAndExecute(query string, args ...interface{}) error {
	stmt, stmtErr := GetDB().Prepare(query)
	if stmtErr != nil {
		log.Printf("Error preparing query: %s", stmtErr)
		return stmtErr
	}

	defer stmt.Close()

	_, execErr := stmt.Exec(args...)
	if execErr != nil {
		log.Printf("Error executing query: %s", execErr)
		return execErr
	}

	return nil
}
