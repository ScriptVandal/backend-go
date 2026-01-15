package repositories

import (
	"database/sql"
	"time"

	"github.com/ScriptVandal/backend-go/internal/models"
)

type RefreshTokenRepository interface {
	Create(token *models.RefreshToken) error
	GetByJTI(jti string) (*models.RefreshToken, error)
	Revoke(jti string) error
	DeleteExpired() error
}

type PGRefreshTokenRepository struct {
	db *sql.DB
}

func NewPGRefreshTokenRepository(db *sql.DB) *PGRefreshTokenRepository {
	return &PGRefreshTokenRepository{db: db}
}

func (r *PGRefreshTokenRepository) Create(token *models.RefreshToken) error {
	query := `INSERT INTO refresh_tokens (jti, user_id, expires_at, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, token.JTI, token.UserID, token.ExpiresAt, token.CreatedAt)
	return err
}

func (r *PGRefreshTokenRepository) GetByJTI(jti string) (*models.RefreshToken, error) {
	query := `SELECT jti, user_id, expires_at, revoked_at, created_at FROM refresh_tokens WHERE jti = $1`
	var token models.RefreshToken
	err := r.db.QueryRow(query, jti).Scan(&token.JTI, &token.UserID, &token.ExpiresAt, &token.RevokedAt, &token.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *PGRefreshTokenRepository) Revoke(jti string) error {
	query := `UPDATE refresh_tokens SET revoked_at = $1 WHERE jti = $2`
	_, err := r.db.Exec(query, time.Now(), jti)
	return err
}

func (r *PGRefreshTokenRepository) DeleteExpired() error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < $1`
	_, err := r.db.Exec(query, time.Now())
	return err
}
