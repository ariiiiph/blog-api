package repository

import (
	"blogapi/internal/models"
	"database/sql"
	"time"
)

type RefreshTokenRepository struct {
	DB *sql.DB
}

func (r *RefreshTokenRepository) CreateRefreshToken(userID int64, token string, expiresAt time.Time) error {
	_, err := r.DB.Exec(
		`INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES (?, ?, ?)`,
		userID,
		token,
		expiresAt.Format(time.RFC3339),
	)
	return err
}

func (r *RefreshTokenRepository) GetRefreshToken(token string) (*models.RefreshToken, error) {
	rt := &models.RefreshToken{}
	err := r.DB.QueryRow(
		`SELECT id, user_id, token, expires_at, created_at FROM refresh_tokens WHERE token = ?`,
		token,
	).Scan(&rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt, &rt.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return rt, nil
}

func (r *RefreshTokenRepository) DeleteRefreshToken(token string) error {
	result, err := r.DB.Exec(`DELETE FROM refresh_tokens WHERE token = ?`, token)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *RefreshTokenRepository) DeleteUserRefreshTokens(userID int64) error {
	_, err := r.DB.Exec(
		`DELETE FROM refresh_tokens WHERE user_id = ?`,
		userID,
	)
	return err
}
