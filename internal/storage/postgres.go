package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/kapi1023/kingGivefivecoins-Authentication/internal/models"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStorage{db: db}, nil
}

func (s *PostgresStorage) CreateUser(email, passwordHash string) error {
	_, err := s.db.Exec(
		"INSERT INTO users (email, password_hash) VALUES ($1, $2) ON CONFLICT (email) DO NOTHING",
		email, passwordHash,
	)
	if err != nil {
		return errors.New("could not insert user into database: " + err.Error())
	}
	return nil
}
func (s *PostgresStorage) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	row := s.db.QueryRow("SELECT id, email, password_hash, oauth_provider, oauth_id, created_at FROM users WHERE email = $1", email)
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.OAuthProvider, &user.OAuthID, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

func (s *PostgresStorage) CreateOAuthUser(email, provider, oauthID string) (*models.User, error) {
	_, err := s.db.Exec(
		"INSERT INTO users (email, oauth_provider, oauth_id) VALUES ($1, $2, $3) ON CONFLICT (email) DO NOTHING",
		email, provider, oauthID,
	)
	if err != nil {
		return nil, err
	}
	return s.GetUserByEmail(email)
}

func (s *PostgresStorage) GetOrCreateUserByOAuth(email, provider string) (*models.User, error) {
	user, err := s.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if user != nil {
		return user, nil
	}

	_, err = s.db.Exec(
		"INSERT INTO users (email, oauth_provider) VALUES ($1, $2)",
		email, provider,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user for OAuth: %w", err)
	}

	return s.GetUserByEmail(email)
}
