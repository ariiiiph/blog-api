package repository

import (
	"blogapi/internal/models"
	"database/sql"
)

type PostRepository struct {
	DB *sql.DB
}

func (r *PostRepository) CreatePost(authorId int64, categoryID int64, title, slug, content string, published bool) (*models.Post, error) {

	res, err := r.DB.Exec(
		"INSERT INTO posts (author_id,category_id, title, slug, content, published) VALUES (? ,?, ?, ?, ?, ?)",
		authorId,
		categoryID,
		title,
		slug,
		content,
		published,
	)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.GetPostByID(id)
}

func (r *PostRepository) CountPosts(search string) (int, error) {
	var count int
	var err error

	if search == "" {
		err = r.DB.QueryRow(`SELECT COUNT(*) FROM posts`).Scan(&count)
	} else {
		like := "%" + search + "%"
		err = r.DB.QueryRow(`SELECT COUNT(*) FROM posts WHERE title LIKE ? OR content LIKE ?`, like, like).Scan(&count)
	}
	if err != nil {
		return 0, err
	}
	return count, nil

}

func (r *PostRepository) GetPosts(limit, offset int, search string) ([]*models.Post, error) {
	var rows *sql.Rows
	var err error

	if search == "" {
		rows, err = r.DB.Query(
			`SELECT id, author_id, category_id, title, slug, content, published, image_url, created_at, updated_at
            FROM posts
            ORDER BY created_at DESC
            LIMIT ? OFFSET ?`,
			limit, offset,
		)
	} else {
		like := "%" + search + "%"
		rows, err = r.DB.Query(
			`SELECT id, author_id, category_id, title, slug, content, published, image_url, created_at, updated_at
            FROM posts
            WHERE title LIKE ? OR content LIKE ?
            ORDER BY created_at DESC
            LIMIT ? OFFSET ?`,
			like, like, limit, offset,
		)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []*models.Post

	for rows.Next() {
		p := &models.Post{}

		err = rows.Scan(&p.ID, &p.AuthorID, &p.CategoryID, &p.Title, &p.Slug, &p.Content, &p.Published, &p.ImageURL, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *PostRepository) GetPostByID(id int64) (*models.Post, error) {
	p := &models.Post{}

	rows := r.DB.QueryRow(
		`SELECT id, author_id, category_id, title, slug, content, published, image_url, created_at, updated_at
        FROM posts
        WHERE id = ?
        `,
		id,
	)
	err := rows.Scan(&p.ID, &p.AuthorID, &p.CategoryID, &p.Title, &p.Slug, &p.Content, &p.Published, &p.ImageURL, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PostRepository) UpdatePost(postID int64, authorID int64, categoryID int64, title, slug, content string, published bool) error {
	result, err := r.DB.Exec(
		`UPDATE posts SET category_id = ?, title = ?, slug = ?, content = ?, published = ?, updated_at = datetime('now') WHERE id = ? AND author_id = ?`,
		categoryID,
		title,
		slug,
		content,
		published,
		postID,
		authorID,
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

func (r *PostRepository) UpdatePostImage(postID int64, imageURL string) error {
	result, err := r.DB.Exec(
		`UPDATE posts SET image_url = ?, updated_at = datetime('now') WHERE id = ?`,
		imageURL, postID,
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

func (r *PostRepository) DeletePost(postId, authorId int64, isAdmin bool) error {
	var result sql.Result
	var err error

	if isAdmin {
		result, err = r.DB.Exec(`DELETE FROM posts WHERE id = ?`, postId)
	} else {
		result, err = r.DB.Exec(`DELETE FROM posts WHERE id = ? AND author_id = ?`, postId, authorId)
	}
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
