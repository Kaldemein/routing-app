package emailtokens

import (
	"database/sql"
)

type PostgresEmailRepository struct {
	db *sql.DB
}

func NewPostgresEmailRepository(db *sql.DB) *PostgresEmailRepository {
	return &PostgresEmailRepository{db}
}

func (repo *PostgresEmailRepository) Save(email, token string) error {
	query := `INSERT INTO email_tokens (email, token) VALUES ($1, $2)`
	_, err := repo.db.Exec(query, email, token)
	return err
}

func (repo *PostgresEmailRepository) FindByToken(token string) (string, error) {
	var email string
	query := `SELECT email FROM email_tokens WHERE token=$1`
	err := repo.db.QueryRow(query, token).Scan(&email)
	return email, err
}
