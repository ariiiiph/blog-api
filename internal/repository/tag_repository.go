package repository

import (
	"blogapi/internal/models"
	"database/sql"
)

type TagRepository struct {
	DB *sql.DB
}

func (r *TagRepository) CreateTag(name, slug string) (*models.Tag, error) {
	res, err := r.DB.Exec(
		`INSERT INTO tags (name, slug) VALUES (?, ?)`,
		name,
		slug,
	)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.GetTagByID(id)
}

func (r *TagRepository) GetTags() ([]*models.Tag, error) {
	rows, err := r.DB.Query(`SELECT id, name, slug, created_at FROM tags`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*models.Tag

	for rows.Next() {
		t := &models.Tag{}
		err := rows.Scan(&t.ID, &t.Name, &t.Slug, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *TagRepository) GetTagByID(id int64) (*models.Tag, error) {
	t := &models.Tag{}

	row := r.DB.QueryRow(`SELECT id, name, slug, created_at FROM tags WHERE id = ?`, id)

	err := row.Scan(&t.ID, &t.Name, &t.Slug, &t.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return t, nil
}

func (r *TagRepository) UpdateTag(id int64, name, slug string) error {
	result, err := r.DB.Exec(`UPDATE tags SET name = ?, slug = ? WHERE id = ?`, name, slug, id)
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

func (r *TagRepository) DeleteTag(id int64) error {
	result, err := r.DB.Exec(
		`DELETE FROM tags WHERE id = ?`, id,
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
