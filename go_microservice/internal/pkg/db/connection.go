package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func ConnectToDB(psqlInfo string) (*sql.DB, error) {
	var db *sql.DB
	var err error
	for i := 0; i < 5; i++ {
		log.Printf("Trying to connect postgres [%v/5]", i)
		db, err = sql.Open("postgres", psqlInfo)
		if err == nil {
			err = db.Ping()
			if err == nil {
				return db, nil
			}
		}
		time.Sleep(5 * time.Second)
	}
	return nil, err
}
