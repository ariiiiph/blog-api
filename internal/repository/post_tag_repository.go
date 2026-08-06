package repository

import (
	"blogapi/internal/models"
	"database/sql"
)

type PostTagRepository struct {
	DB *sql.DB
}

func (r *PostTagRepository) AddTagToPost(postID, tagID int64) error {
	_, err := r.DB.Exec(`INSERT INTO post_tags (post_id, tag_id) VALUES (?, ?)`, postID, tagID)
	return err
}

func (r *PostTagRepository) RemoveTagFromPost(postID, tagID int64) error {
	result, err := r.DB.Exec(`DELETE FROM post_tags WHERE post_id = ? AND tag_id = ?`, postID, tagID)

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
func (r *PostTagRepository) GetTagsByPost(postID int64) ([]*models.Tag, error) {
	rows, err := r.DB.Query(
		`SELECT tags.id, tags.name, tags.slug, tags.created_at
		FROM tags
		JOIN post_tags ON tags.id = post_tags.tag_id
		WHERE post_tags.post_id = ?`,
		postID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*models.Tag
	for rows.Next() {
		t := &models.Tag{}
		err = rows.Scan(&t.ID, &t.Name, &t.Slug, &t.CreatedAt)
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

func (r *PostTagRepository) GetPostsByTag(tagID int64) ([]*models.Post, error) {
	rows, err := r.DB.Query(
		`SELECT posts.id, posts.author_id,posts.category_id,posts.title, posts.slug, posts.content, posts.published, posts.created_at, posts.updated_at
		FROM posts
		JOIN post_tags ON posts.id = post_tags.post_id
		WHERE post_tags.tag_id = ?`,
		tagID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		p := &models.Post{}
		err = rows.Scan(&p.ID, &p.AuthorID, &p.CategoryID, &p.Title, &p.Slug, &p.Content, &p.Published, &p.CreatedAt, &p.UpdatedAt)
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
