package database

import (
	"database/sql"
	"fmt"
	"forum/models"
	"log"
)

func GetCommentLikes(db *sql.DB, commentId int) ([]models.CommentLikes, error) {
	var likes []models.CommentLikes
	q := `SELECT * FROM CommentLikes WHERE comment_id = ?;`
	rows, err := db.Query(q, commentId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, commentId, userId, likeDislike int
		err = rows.Scan(&id, &commentId, &userId, &likeDislike)
		if err != nil {
			fmt.Println("getCommentLikes: scan failed")
		}
		likes = append(likes, models.CommentLikes{
			Id:          id,
			CommentId:   commentId,
			AuthorId:    userId,
			LikeDislike: likeDislike,
		})
	}
	return likes, nil
}

func InsertCommentLike(db *sql.DB, commentLikes models.CommentLikes) error {
	likesSQL := "INSERT INTO CommentLikes(comment_id, user_id, like_dislike) VALUES (?,?,?)"
	statement, err := db.Prepare(likesSQL)
	if err != nil {
		return err
	}
	_, err = statement.Exec(commentLikes.CommentId, commentLikes.AuthorId, commentLikes.LikeDislike)
	if err != nil {
		return err
	}
	return nil
}

func DeleteCommentLike(db *sql.DB, commentLikes models.CommentLikes) error {
	_, err := db.Exec("DELETE FROM CommentLikes WHERE id = ?", commentLikes.Id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateCommentLike(db *sql.DB, commentLikes models.CommentLikes) error {
	_, err := db.Exec("UPDATE CommentLikes SET like_dislike = ? WHERE id = ?", commentLikes.LikeDislike, commentLikes.Id)
	if err != nil {
		return err
	}
	return nil
}

func CreateCommentLikesTable(db *sql.DB) error {
	commentLikesTable := `CREATE TABLE CommentLikes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,	
		comment_id INTEGER references postsTable(id),
		user_id INTEGER references User(id),
		like_dislike INTEGER
	  );`

	log.Println("Creating CommentLikes table...")
	statement, err := db.Prepare(commentLikesTable)
	if err != nil {
		return err
	}
	statement.Exec()
	log.Println("Table created")
	return err
}

func CreatePostLikesTable(db *sql.DB) error {
	postLikesTable := `CREATE TABLE PostLikes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,	
		posts_id INTEGER references postsTable(id),
		user_id INTEGER references User(id),
		like_dislike INTEGER
	  );`

	log.Println("Creating PostLikes table...")
	statement, err := db.Prepare(postLikesTable)
	if err != nil {
		return err
	}
	statement.Exec()
	log.Println("Table created")
	return nil
}

func GetPostLikes(db *sql.DB, postId int) ([]models.PostLikes, error) {
	var likes []models.PostLikes

	q := `SELECT * FROM PostLikes WHERE posts_id = ?;`
	rows, err := db.Query(q, postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, postId, userId, likeDislike int
		err = rows.Scan(&id, &postId, &userId, &likeDislike)
		if err != nil {
			fmt.Println("getPostLikes: scan failed")
			return nil, err
		}
		likes = append(likes, models.PostLikes{
			Id:          id,
			PostId:      postId,
			AuthorId:    userId,
			LikeDislike: likeDislike,
		})
	}
	return likes, nil
}

func InsertPostLike(db *sql.DB, postLikes models.PostLikes) error {
	likesSQL := "INSERT INTO PostLikes(posts_id, user_id, like_dislike) VALUES (?,?,?)"
	statement, err := db.Prepare(likesSQL)
	if err != nil {
		return err
	}
	_, err = statement.Exec(postLikes.PostId, postLikes.AuthorId, postLikes.LikeDislike)
	if err != nil {
		return err
	}
	return nil
}

func DeletePostLike(db *sql.DB, postLikes models.PostLikes) error {
	_, err := db.Exec("DELETE FROM PostLikes WHERE id = ?", postLikes.Id)
	if err != nil {
		return err
	}
	return nil
}

func UpdatePostLike(db *sql.DB, postLikes models.PostLikes) error {
	_, err := db.Exec("UPDATE PostLikes SET like_dislike = ? WHERE id = ?", postLikes.LikeDislike, postLikes.Id)
	if err != nil {
		return err
	}
	return nil
}
