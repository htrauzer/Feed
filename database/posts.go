package database

import (
	"database/sql"
	"fmt"
	"forum/models"
	"log"
	"time"
)

func CreatePostsTable(db *sql.DB) {
	postsTable := `CREATE TABLE Posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,		
		content TEXT,
		userId INTEGER references Users(id),
		userName TEXT,
		created_at TEXT,
		title TEXT		
	  );`

	log.Println("Creating Posts table...")
	statement, err := db.Prepare(postsTable)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("Table created")
}

func InsertPost(db *sql.DB, userId int, content, userName, title string) int64 {
	created_at := time.Now().Format("01-02-2006 15:04:05")
	postSQL := "INSERT INTO Posts(content, userId, userName, created_at, title) VALUES (?,?,?,?,?)"
	var statement *sql.Stmt
	if res, err := db.Prepare(postSQL); err != nil {
		log.Fatal(err.Error())
	} else {
		statement = res
	}

	var post sql.Result
	if res, err := statement.Exec(content, userId, userName, created_at, title); err != nil {
		log.Fatalln(err.Error())
	} else {
		post = res
	}

	postId, _ := post.LastInsertId()
	return postId
}

func GetPost(db *sql.DB, identifier string) (models.Post, error) {
	var post models.Post

	q := "SELECT * FROM Posts WHERE id = ?;"
	row, _ := db.Query(q, identifier)
	var id, userId int
	var content, userName, createdAt, title string
	for row.Next() {
		switch err := row.Scan(
			&id,
			&content,
			&userId,
			&userName,
			&createdAt,
			&title); err {
		case sql.ErrNoRows:
			fmt.Println("getSinglePost: no rows returned!")
		case nil:
			post = models.Post{
				Id:           id,
				Content:      content,
				AuthorId:     userId,
				UserName:     userName,
				CreationTime: createdAt,
				Title:        title,
			}
		default:

			fmt.Println("getSinglePost: some other error")
			fmt.Println(err)
			return post, err
		}
	}

	return post, nil
}

func GetPosts(db *sql.DB) ([]models.Post, error) {
	var posts []models.Post

	rows, err := db.Query("SELECT * FROM Posts ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var content, userName, created_at, title string
		var id, userId int
		rows.Scan(&id, &content, &userId, &userName, &created_at, &title)

		posts = append(posts, models.Post{
			Id:           id,
			Content:      content,
			AuthorId:     userId,
			UserName:     userName,
			CreationTime: created_at,
			Title:        title,
		})
	}
	return posts, nil
}
