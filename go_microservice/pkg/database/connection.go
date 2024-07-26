package database

import (
	"database/sql"
	"fmt"
	"time"
)

func ConnectToDB(psqlInfo string, retryCount int, retryInterval time.Duration) (*sql.DB, error) {
	var db *sql.DB
	var err error
	for i := 0; i < retryCount; i++ {
		db, err = sql.Open("postgres", psqlInfo)
		if err == nil {
			err = db.Ping()
			if err == nil {
				return db, nil
			}
		}
		fmt.Printf("Failed to connect to database. Retrying in %v seconds...\n", retryInterval.Seconds())
		time.Sleep(retryInterval)
	}
	return nil, err
}
