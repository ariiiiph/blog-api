package repository

import (
	"blogapi/internal/models"
	"database/sql"
)

type CommentRepository struct {
	DB *sql.DB
}

func (r *CommentRepository) SendComment(postID, authorID int64, content string) (*models.Comment, error) {
	res, err := r.DB.Exec(`
	INSERT INTO comments (post_id, author_id, content) VALUES ( ?, ?, ?)`, postID, authorID, content,
	)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &models.Comment{
		ID:       id,
		PostID:   postID,
		AuthorID: authorID,
		Content:  content,
	}, nil
}

func (r *CommentRepository) CountComments(postId int64) (int, error) {
	var count int

	err := r.DB.QueryRow("SELECT COUNT(*) FROM comments WHERE post_id = ?", postId).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil

}
func (r *CommentRepository) GetComments(postId int64, limit, offset int) ([]*models.Comment, error) {
	rows, err := r.DB.Query(
		`SELECT id, post_id, author_id, content, created_at
		FROM comments
		WHERE post_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`, postId, limit, offset,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []*models.Comment

	for rows.Next() {
		c := &models.Comment{}

		err = rows.Scan(&c.ID, &c.PostID, &c.AuthorID, &c.Content, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *CommentRepository) UpdateComment(commentID, authorID int64, content string) error {
	result, err := r.DB.Exec(`UPDATE comments SET content = ? WHERE id = ? AND author_id = ?`, content, commentID, authorID)
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
func (r *CommentRepository) DeleteComment(commentID, authorID int64, isAdmin bool) error {
	var result sql.Result
	var err error

	if isAdmin {
		result, err = r.DB.Exec(`DELETE FROM comments WHERE id = ?`, commentID)
	} else {
		result, err = r.DB.Exec(`DELETE FROM comments WHERE id = ? AND author_id = ?`, commentID, authorID)
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
