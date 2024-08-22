package emailtokens

import (
	"database/sql"
	"service/internal/pkg/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	dbConn := getTestDb()
	repo := NewPostgresEmailRepository(getTestDb())

	expEmail := "ruzel@email.com"
	expToken := "token"
	dbConn.Exec("TRUNCATE TABLE email_tokens;")

	err := repo.Save(expEmail, expToken)
	assert.NoError(t, err)

	email, err := repo.FindByToken(expToken)
	assert.NoError(t, err)
	assert.Equal(t, expEmail, email)
}

func getTestDb() *sql.DB {
	// Connection to postgres
	psqlInfo := "host=go_db port=5432 user=postgres password=postgres dbname=email_verification_test sslmode=disable"

	// Database connection
	dbConn, _ := db.ConnectToDB(psqlInfo)
	return dbConn
}
