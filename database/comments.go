package database

import (
	"database/sql"
	"fmt"
	"forum/models"
	"log"
	"time"
)

func GetComments(db *sql.DB, post models.Post) ([]models.Comment, error) {
	var comments []models.Comment

	q := `SELECT id, content, userId, posts_id, created_at FROM Comments WHERE posts_id = ?;`

	var rows *sql.Rows
	if res, err := db.Query(q, post.Id); err != nil {
		return nil, err
	} else {
		rows = res
	}

	defer rows.Close()

	for rows.Next() {
		var id, posts_id int
		var content, createdAt, userId string

		if err := rows.Scan(&id, &content, &userId, &posts_id, &createdAt); err != nil {
			fmt.Println("comments.go - getComments, FAILED")
			return nil, err
		}

		comments = append(comments, models.Comment{
			Id:           id,
			Content:      content,
			Author:       userId,
			PostId:       posts_id,
			CreationTime: createdAt,
		})
	}
	return comments, nil
}

func CreateCommentsTable(db *sql.DB) error {
	commentsTable := `CREATE TABLE Comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,		
		content TEXT,
		userId TEXT,
		posts_id INTEGER references Posts(rowid),
		created_at TEXT
	  );`

	log.Println("Creating commentsTable table...")
	statement, err := db.Prepare(commentsTable)
	if err != nil {
		return err
	}
	statement.Exec()
	log.Println("Table created")
	return nil
}

func InsertComment(db *sql.DB, comment models.Comment) error {
	creationTime := time.Now().String()
	q := `INSERT INTO Comments(content, userId, posts_id, created_at) VALUES (?,?,?,?);`
	statement, err := db.Prepare(q)
	if err != nil {
		return err
	}
	_, err = statement.Exec(comment.Content, comment.Author, comment.PostId, creationTime)
	if err != nil {
		return err
	}
	return nil
}
