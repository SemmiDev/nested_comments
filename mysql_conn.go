package nested_comments

import (
	"context"
	"database/sql"
	"log"
	"time"
)

func LoadMysqlConnection(connectionStr string) (*sql.DB, error) {
	db, err := sql.Open("mysql", connectionStr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(timeoutCtx)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to database")

	return db, nil
}

func CloseMysqlConnection(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Println("Failed to close database connection")
	}
}
