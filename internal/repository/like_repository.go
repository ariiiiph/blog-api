package repository

import "database/sql"

type LikeRepository struct {
	DB *sql.DB
}

func (r *LikeRepository) LikePost(userID, postID int64) error {
	_, err := r.DB.Exec(
		`INSERT INTO likes (user_id, post_id) VALUES (?, ?)`,
		userID, postID,
	)
	return err
}

func (r *LikeRepository) UnlikePost(userID, postID int64) error {
	result, err := r.DB.Exec(`DELETE FROM likes WHERE user_id = ? AND post_id = ?`, userID, postID)

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
func (r *LikeRepository) CountLikes(postID int64) (int64, error) {
	var count int64
	err := r.DB.QueryRow(
		`SELECT COUNT(*) FROM likes WHERE post_id = ?`,
		postID,
	).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil

}
func (r *LikeRepository) HasUserLikedPost(userID, postID int64) (bool, error) {
	var exists bool
	err := r.DB.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = ? AND post_id = ?)`,
		userID, postID,
	).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil

}
