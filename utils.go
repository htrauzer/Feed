package main

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/database"
	"forum/models"
	"log"
	"net/http"
	"net/mail"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func signup(email, username, password string) error {
	encryptedPassWord, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		error := errors.New("password encryption failed")
		fmt.Println(error)
		return error
	}

	user := models.User{0, email, username, string(encryptedPassWord)}

	if err := database.InsertUser(sqliteDatabase, user); err != nil {
		return err
	}
	return nil
}

func indexPosts(user models.User) ([]models.Post, error) {
	var posts []models.Post
	if res, err := database.GetPosts(sqliteDatabase); err != nil {
		return nil, err
	} else {
		posts = res
	}

	for i, post := range posts {
		if res, err := populatePostData(post, user); err != nil {
			return posts, err
		} else {
			posts[i] = res
		}
	}

	return posts, nil
}

func showPost(id string, user models.User) (models.Post, error) {
	var post models.Post

	if res, err := database.GetPost(sqliteDatabase, id); err != nil {
		return post, err
	} else {
		post = res
	}

	if res, err := populatePostData(post, user); err != nil {
		return post, err
	} else {
		post = res
	}

	return post, nil
}

func isEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil

}

func populatePostData(post models.Post, user models.User) (models.Post, error) {
	var comments []models.Comment
	if res, err := database.GetComments(sqliteDatabase, post); err != nil {
		return post, err
	} else {
		comments = res
	}
	for i := range comments {
		var commentLikes []models.CommentLikes
		if res, err := database.GetCommentLikes(sqliteDatabase, comments[i].Id); err != nil {
			return post, err
		} else {
			commentLikes = res
		}

		var likes int
		var dislikes int
		var userHasReacted int

		for _, cL := range commentLikes {
			if cL.LikeDislike == 1 {
				likes++
			} else if cL.LikeDislike == 2 {
				dislikes++
			}
			if cL.AuthorId == user.Id { // only if logged in
				userHasReacted = cL.LikeDislike
			}
		}
		comments[i].Likes = likes
		comments[i].Dislikes = dislikes
		comments[i].UserLikeStatus = userHasReacted
	}
	post.Comments = comments
	post.NumComments = len(comments)

	var postLikes []models.PostLikes
	if res, err := database.GetPostLikes(sqliteDatabase, post.Id); err != nil {
		return post, err
	} else {
		postLikes = res
	}

	var likes int
	var dislikes int
	var userHasReacted int

	for _, pL := range postLikes {
		if pL.LikeDislike == 1 {
			likes++
		} else if pL.LikeDislike == 2 {
			dislikes++
		}
		if pL.AuthorId == user.Id { // only if logged in
			userHasReacted = pL.LikeDislike
		}
	}
	post.Likes = likes
	post.Dislikes = dislikes
	post.UserLikeStatus = userHasReacted

	return post, nil
}

func IsLoggedIn(r *http.Request) (bool, string) {
	var answer bool
	var who string
	// get cookie "session" from browser request if available
	cookie, err := r.Cookie("session")
	if err != nil {
		return false, "NA"
	}

	answer, who, _ = database.HasSession(sqliteDatabase, cookie.Value)
	return answer, who
}

func ReCreateAndConnSqlDataBase() {

	new := false
	_, err := os.Stat("sqlite.db")

	if os.IsNotExist(err) {
		file, err2 := os.Create("sqlite.db") // Create SQLite file
		if err2 != nil {
			log.Fatal(err.Error())

		}
		file.Close()
		log.Println("sqlite.db created")
		new = true
	}

	sqliteDatabase, err = sql.Open("sqlite3", "./sqlite.db")
	if err != nil {
		log.Fatalln(err.Error())
	}
	if new {
		database.CreateUsersTable(sqliteDatabase)
		database.CreateCommentsTable(sqliteDatabase)
		database.CreatePostCategoriesTable(sqliteDatabase)
		database.CreateSessionsTable(sqliteDatabase)
		database.CreatePostsTable(sqliteDatabase)
		database.CreatePostLikesTable(sqliteDatabase)
		database.CreateCommentLikesTable(sqliteDatabase)
		database.CreateCategoriesTable(sqliteDatabase)
		database.InsertCategory(sqliteDatabase)
	}
}

func CanRegister(db *sql.DB, inputEmail, inputUserName string) (answer bool, reason string) {
	answer = true

	row, err := db.Query("SELECT email, username FROM Users")
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer row.Close()

	for row.Next() {
		var email, username string
		row.Scan(&email, &username)
		if email == inputEmail {
			answer = false
			reason = "User with that email already registred"
			return answer, reason
		} else if username == inputUserName {
			answer = false
			reason = "Username already taken"
			return answer, reason
		}
	}
	return
}
