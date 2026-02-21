package database

import (
	"database/sql"
	"forum/models"
	"log"
	"strconv"
)

func CreateUsersTable(db *sql.DB) error {
	userTable := `CREATE TABLE Users (
		id INTEGER PRIMARY KEY,		
		email TEXT,
		username TEXT,
		password TEXT	
	  );`

	log.Println("Creating Users table...")
	statement, err := db.Prepare(userTable)
	if err != nil {
		return err
	}
	statement.Exec()
	log.Println("Table created")
	return nil
}

func InsertUser(db *sql.DB, usr models.User) error {
	userSQL := `INSERT INTO Users(email, username, password) VALUES (?,?,?);`
	statement, err := db.Prepare(userSQL)
	if err != nil {
		return err
	}
	_, err = statement.Exec(usr.Email, usr.UserName, usr.PassWord)
	if err != nil {
		return err
	}
	return nil
}

func GetUserPostLikes(db *sql.DB, postId, userId int) (models.PostLikes, error) {
	var likes models.PostLikes

	q := `SELECT * FROM PostLikes WHERE posts_id = ? AND user_id = ?;`
	rows, err := db.Query(q, postId, userId)
	if err != nil {
		return likes, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, postId, userId, likeDislike int
		err = rows.Scan(&id, &postId, &userId, &likeDislike)
		if err != nil {
			return likes, err
		}
		likes = models.PostLikes{
			Id:          id,
			PostId:      postId,
			AuthorId:    userId,
			LikeDislike: likeDislike,
		}
	}
	return likes, nil
}

func GetUserCommentLikes(db *sql.DB, commentId, userId int) (models.CommentLikes, error) {
	var likes models.CommentLikes

	q := `SELECT * FROM CommentLikes WHERE comment_id = ? AND user_id = ?;`
	rows, err := db.Query(q, commentId, userId)
	if err != nil {
		return likes, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, commentId, userId, likeDislike int
		err = rows.Scan(&id, &commentId, &userId, &likeDislike)
		if err != nil {
			return likes, err
		}
		likes = models.CommentLikes{
			Id:          id,
			CommentId:   commentId,
			AuthorId:    userId,
			LikeDislike: likeDislike,
		}
	}
	return likes, nil
}

func GetUserDataStruct(db *sql.DB, userName string) (models.User, error) {
	var user models.User

	q := `SELECT rowid, email, username, password FROM Users WHERE username = ?;`

	row, err := db.Query(q, userName)
	if err != nil {
		return user, err
	}
	var rowid int
	var email, username, password string

	for row.Next() {
		err = row.Scan(&rowid, &email, &username, &password)
		if err != nil {
			return user, err
		}
		user = models.User{
			Id:       rowid,
			Email:    email,
			UserName: username,
			PassWord: password,
		}
	}

	return user, nil
}

func GetUserPosts(db *sql.DB, userId int) ([]models.Post, error) {
	var posts []models.Post

	rows, err := db.Query("SELECT * FROM Posts WHERE userId=? ORDER BY id", userId)
	if err != nil {
		return posts, nil
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

func GetUserLikedPosts(db *sql.DB, userId int) ([]models.Post, error) {
	//query all posts which the user has liked, 1 = like
	rows, err := db.Query("SELECT * FROM postLikes WHERE user_id = ? AND like_dislike = 1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var likedPostIds []int
	for rows.Next() {
		var id, posts_id, user_id, like_dislike int
		rows.Scan(&id, &posts_id, &user_id, &like_dislike)

		likedPostIds = append(likedPostIds, posts_id)
	}

	// get all the posts by id
	var likedPosts []models.Post
	for _, id := range likedPostIds {
		idAsStr := strconv.Itoa(id)

		var post models.Post
		if res, err := GetPost(db, idAsStr); err != nil {
			return nil, err
		} else {
			post = res
		}

		likedPosts = append(likedPosts, post)
	}
	return likedPosts, nil
}

func GetUserDislikedPosts(db *sql.DB, userId int) ([]models.Post, error) {
	rows, err := db.Query("SELECT * FROM postLikes WHERE user_id = ? AND like_dislike = 2", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notlikedPostIds []int
	for rows.Next() {
		var id, posts_id, user_id, like_dislike int
		rows.Scan(&id, &posts_id, &user_id, &like_dislike)

		notlikedPostIds = append(notlikedPostIds, posts_id)
	}

	// get all the posts by id
	var notlikedPosts []models.Post
	for _, id := range notlikedPostIds {
		idAsStr := strconv.Itoa(id)

		var post models.Post
		if res, err := GetPost(db, idAsStr); err != nil {
			return nil, err
		} else {
			post = res
		}

		notlikedPosts = append(notlikedPosts, post)
	}
	return notlikedPosts, nil
}
