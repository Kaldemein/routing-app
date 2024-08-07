package database

import (
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func TestConnectToDB(t *testing.T) {

	psqlInfo := "host=go_db port=5432 user=postgres " + "password=postgres dbname=email_verification sslmode=disable"

	retryCount := 3
	retryInterval := 10 * time.Second

	db, err := ConnectToDB(psqlInfo, retryCount, retryInterval)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}
}
