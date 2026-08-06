package repository

import (
	"blogapi/internal/models"
	"database/sql"
)

type CategoryRepository struct {
	DB *sql.DB
}

func (r *CategoryRepository) CreateCategory(name, slug string) (*models.Category, error) {
	res, err := r.DB.Exec("INSERT INTO categories (name, slug) VALUES (?, ?)", name, slug)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.GetCategoryByID(id)
}

func (r *CategoryRepository) GetCategories() ([]*models.Category, error) {
	rows, err := r.DB.Query(`SELECT id, name, slug, created_at FROM categories`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*models.Category

	for rows.Next() {
		c := &models.Category{}
		err = rows.Scan(&c.ID, &c.Name, &c.Slug, &c.CreatedAt)
		if err != nil {
			return nil, err
		}

		categories = append(categories, c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoryRepository) GetCategoryByID(id int64) (*models.Category, error) {
	c := &models.Category{}

	rows := r.DB.QueryRow(
		`SELECT id, name, slug, created_at FROM categories WHERE id = ?`, id,
	)
	err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return c, nil

}

func (r *CategoryRepository) UpdateCategory(categoryID int64, name, slug string) error {
	result, err := r.DB.Exec(`UPDATE categories SET name = ?, slug = ? WHERE id = ?`, name, slug, categoryID)
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

func (r *CategoryRepository) DeleteCategory(id int64) error {
	result, err := r.DB.Exec(
		`DELETE FROM categories WHERE id = ?`, id,
	)
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
