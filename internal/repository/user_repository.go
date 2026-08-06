package repository

import (
	"blogapi/internal/models"
	"database/sql"
)

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) CreateUser(email, hashedPassword string) (*models.User, error) {
	res, err := r.DB.Exec(
		"INSERT INTO users (email, password_hash) VALUES (?, ?)",
		email,
		hashedPassword,
	)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.GetUserByID(id)
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	u := &models.User{}

	row := r.DB.QueryRow(
		`SELECT id, email, password_hash, role, created_at 
        FROM users
        WHERE email = ?`,
		email,
	)

	err := row.Scan(&u.ID, &u.Email, &u.Password, &u.Role, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) GetUserByID(id int64) (*models.User, error) {
	u := &models.User{}

	row := r.DB.QueryRow(
		`SELECT id, email, password_hash, role, created_at 
        FROM users
        WHERE id = ?`,
		id,
	)

	err := row.Scan(&u.ID, &u.Email, &u.Password, &u.Role, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return u, nil
}
