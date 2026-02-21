package main

import (
	"fmt"
	"forum/database"
	"forum/models"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	validPath     = regexp.MustCompile("^/(posts|categories)/([a-zA-Z0-9-]+)$")
	justRegistred = false
	regUserName   = ""
)

func authMiddleware(w http.ResponseWriter, r *http.Request) bool {
	if isAuthed, _ := IsLoggedIn(r); !isAuthed {
		errorHandler(w, r, http.StatusForbidden)
		return false
	}
	return true
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		validURL := validPath.FindStringSubmatch(req.URL.Path)
		if validURL == nil {
			errorHandler(w, req, http.StatusNotFound)
			return
		}
		fn(w, req, validURL[2])
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	var user models.User

	isAuthed, username := IsLoggedIn(r)

	if isAuthed {
		if res, err := database.GetUserDataStruct(sqliteDatabase, username); err != nil {
			errorHandler(w, r, 500)
			return
		} else {
			user = res
		}
	}

	var posts []models.Post
	if res, err := indexPosts(user); err != nil {
		errorHandler(w, r, 500)
		return
	} else {
		posts = res
	}

	tpl.ExecuteTemplate(w, "index.html", map[string]interface{}{"posts": posts, "user": user})
}

func signupHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		// get form values
		eMail := req.FormValue("email")
		userName := req.FormValue("username")
		passWord := req.FormValue("password")

		data := map[string]interface{}{
			"Email":          eMail,
			"Username":       userName,
			"Password":       passWord,
			"InvalidEmail":   "",
			"NoSignUpReason": "",
		}

		if !isEmailValid(eMail) {
			data["InvalidEmail"] = "Could not verify email domain."
			tpl.ExecuteTemplate(w, "signup.html", data)
			return
		}

		answer, reason := CanRegister(sqliteDatabase, eMail, userName)

		if !answer {
			data["NoSignUpReason"] = reason
			tpl.ExecuteTemplate(w, "signup.html", data)
			return
		}

		err := signup(eMail, userName, passWord)
		if err != nil {
			errorHandler(w, req, http.StatusInternalServerError)
			fmt.Println(err)
		}
		//upon success goto login and auto log in
		justRegistred = true
		regUserName = userName
		loginHandler(w, req)
		return
	}
	tpl.ExecuteTemplate(w, "signup.html", nil)
}

func postsMiddleware(w http.ResponseWriter, r *http.Request, path string) {
	if _, err := strconv.Atoi(path); err == nil {
		postHandler(w, r, path)
		return
	}

	if path == "new" {
		if !authMiddleware(w, r) {
			return
		}
		newPostHandler(w, r)
		return
	}

	if path == "made" {
		if !authMiddleware(w, r) {
			return
		}
		madePostsHandler(w, r)
		return
	}

	if path == "liked" {
		if !authMiddleware(w, r) {
			return
		}
		likedPostsHandler(w, r)
		return
	}
	return
}

func postHandler(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method == http.MethodGet {
		isAuthed, username := IsLoggedIn(r)

		var user models.User
		if res, err := database.GetUserDataStruct(sqliteDatabase, username); err != nil {
			errorHandler(w, r, 500)
			return
		} else {
			user = res
		}

		var post models.Post
		if res, err := showPost(id, user); err != nil {
			errorHandler(w, r, 500)
			return
		} else {
			post = res
		}
		tpl.ExecuteTemplate(w, "post.html", map[string]interface{}{"IsLoggedIn": isAuthed, "post": post, "user": user})
	}
}

func newPostHandler(w http.ResponseWriter, r *http.Request) {
	isAuthed, username := IsLoggedIn(r)

	var user models.User
	if res, err := database.GetUserDataStruct(sqliteDatabase, username); err != nil {
		errorHandler(w, r, 500)
		return
	} else {
		user = res
	}

	if r.Method == http.MethodGet {
		tpl.ExecuteTemplate(w, "new-post.html", map[string]interface{}{"user": user, "IsLoggedIn": isAuthed})
		return
	}

	if r.Method == http.MethodPost {
		content := r.FormValue("content")

		var user models.User
		if res, err := database.GetUserDataStruct(sqliteDatabase, username); err != nil {
			errorHandler(w, r, 500)
			return
		} else {
			user = res
		}

		userId := user.Id
		title := r.FormValue("title")
		postId := database.InsertPost(sqliteDatabase, userId, content, username, title)

		categories := r.Form["categories[]"]

		for _, category := range categories {
			categoryInt, _ := strconv.Atoi(category)
			database.InsertPostCategory(sqliteDatabase, int(postId), categoryInt)
		}

		urlSuffix := "/posts/" + strconv.Itoa(int(postId))
		http.Redirect(w, r, urlSuffix, http.StatusSeeOther)
	}

}

func commentsHandler(w http.ResponseWriter, r *http.Request) {
	if !authMiddleware(w, r) {
		return
	}

	if r.Method == http.MethodPost {
		postrow, _ := strconv.Atoi(r.FormValue("postId"))
		single := r.FormValue("single")

		isAuthed, username := IsLoggedIn(r)

		if !isAuthed {
			if single == "true" {
				http.Redirect(w, r, "/posts/"+strconv.Itoa(postrow), http.StatusSeeOther) 
						}
			http.Redirect(w, r, "/posts", http.StatusSeeOther)
			return
		}

		err := r.ParseForm()
		if err != nil {
			errorHandler(w, r, 500)
		}
		id := 0 // discarded/placeholder value
		content := r.FormValue("commentContent")
		creationTime := time.Now().String()

		comment := models.Comment{
			Id:           id,
			Content:      content,
			Author:       username,
			PostId:       postrow,
			CreationTime: creationTime,
		}

		database.InsertComment(sqliteDatabase, comment)

		if single == "true" {
			http.Redirect(w, r, "/posts/"+strconv.Itoa(postrow), http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/posts", http.StatusSeeOther)
	}
}

func madePostsHandler(w http.ResponseWriter, r *http.Request) {
	_, username := IsLoggedIn(r)

	var user models.User
	if res, err := database.GetUserDataStruct(sqliteDatabase, username); err != nil {
		errorHandler(w, r, 500)
		return
	} else {
		user = res
	}

	var madePosts []models.Post
	if res, err := database.GetUserPosts(sqliteDatabase, user.Id); err != nil {
		errorHandler(w, r, 500)
		return
	} else {
		madePosts = res
	}

	tpl.ExecuteTemplate(w, "posts.html", madePosts)
}

func likedPostsHandler(w http.ResponseWriter, r *http.Request) {
	_, username := IsLoggedIn(r)

	var user models.User
	if res, err := database.GetUserDataStruct(sqliteDatabase, username); err != nil {
		errorHandler(w, r, 500)
		return
	} else {
		user = res
	}

	var madePosts []models.Post
	if res, err := database.GetUserLikedPosts(sqliteDatabase, user.Id); err != nil {
		errorHandler(w, r, 500)
		return
	} else {
		madePosts = res
	}

	tpl.ExecuteTemplate(w, "posts.html", madePosts)
}

func categoryHandler(w http.ResponseWriter, r *http.Request, category string) {
	categories := map[string]string{
		"all":      "All posts",
		"general":  "General",
		"topic-1":  "Topic-1",
		"topic-2":  "Topic-2",
		"topic-3":  "Topic-3",
		"topic-4":  "Topic-4",
		"my-likes": "My likes",
		"my-posts": "My posts",
	}

	_, username := IsLoggedIn(r)

	var user models.User
	if res, err := database.GetUserDataStruct(sqliteDatabase, username); err != nil {
		errorHandler(w, r, 500)
		return
	} else {
		user = res
	}

	var categoryPosts []models.Post
	if category == "all" {
		if res, err := database.GetPosts(sqliteDatabase); err != nil {
			errorHandler(w, r, 500)
		} else {
			categoryPosts = res
		}
	} else if category == "my-likes" {
		if res, err := database.GetUserLikedPosts(sqliteDatabase, user.Id); err != nil {
			errorHandler(w, r, 500)
		} else {
			categoryPosts = res
		}
	} else if category == "my-posts" {
		if res, err := database.GetUserPosts(sqliteDatabase, user.Id); err != nil {
			errorHandler(w, r, 500)
		} else {
			categoryPosts = res
		}
	} else if category == "my-dislikes" {
		if res, err := database.GetUserDislikedPosts(sqliteDatabase, user.Id); err != nil {
			errorHandler(w, r, 500)
		} else {
			categoryPosts = res
		}

	} else {
		if res, err := database.GetCategoryPosts(sqliteDatabase, categories[category]); err != nil {
			errorHandler(w, r, 500)
			return
		} else {
			categoryPosts = res
		}
	}

	for i, post := range categoryPosts {
		if res, err := populatePostData(post, user); err != nil {
			errorHandler(w, r, 500)
			return
		} else {
			categoryPosts[i] = res
		}

	}

	tpl.ExecuteTemplate(w, "index.html", map[string]interface{}{"posts": categoryPosts, "user": user})
}

func reactPostHandler(w http.ResponseWriter, r *http.Request) {
	if !authMiddleware(w, r) {
		return
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			log.Fatal(err.Error())
		}
		single := r.FormValue("single")
		postIdStr := r.FormValue("postId")

		isAuthed, username := IsLoggedIn(r)

		if !isAuthed {
			if single == "true" {
				http.Redirect(w, r, "/posts/"+postIdStr, http.StatusSeeOther)
				return
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		var user models.User
		if res, err := database.GetUserDataStruct(sqliteDatabase, username); err != nil {
			errorHandler(w, r, 500)
			return
		} else {
			user = res
		}

		likeStatus, _ := strconv.Atoi(r.FormValue("likestatus"))
		postId, _ := strconv.Atoi(postIdStr)

		reaction := models.PostLikes{
			PostId:      postId,
			AuthorId:    user.Id,
			LikeDislike: likeStatus, // 1 = like 2 = dislike
		}
		//check if user has already reacted
		var userReactions models.PostLikes
		if res, err := database.GetUserPostLikes(sqliteDatabase, postId, user.Id); err != nil {
			errorHandler(w, r, 500)
			return
		} else {
			userReactions = res
		}
		switch userReactions.LikeDislike {
		case 0:
			if err := database.InsertPostLike(sqliteDatabase, reaction); err != nil {
				errorHandler(w, r, 500)
				return
			}
		case 1:
			if likeStatus == 1 {
				if err := database.DeletePostLike(sqliteDatabase, userReactions); err != nil {
					errorHandler(w, r, 500)
					return
				}
			} else if likeStatus == 2 {
				userReactions.LikeDislike = likeStatus
				if err := database.UpdatePostLike(sqliteDatabase, userReactions); err != nil {
					errorHandler(w, r, 500)
					return
				}
			}
		case 2:
			if likeStatus == 2 {
				if err := database.DeletePostLike(sqliteDatabase, userReactions); err != nil {
					errorHandler(w, r, 500)
					return
				}
			} else if likeStatus == 1 {
				userReactions.LikeDislike = likeStatus

				if err := database.UpdatePostLike(sqliteDatabase, userReactions); err != nil {
					errorHandler(w, r, 500)
					return
				}
			}
		}

		if single == "true" {
			http.Redirect(w, r, "/posts/"+postIdStr, http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

}

func reactCommentHandler(w http.ResponseWriter, r *http.Request) {
	if !authMiddleware(w, r) {
		return
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			log.Fatal(err.Error())
		}
		single := r.FormValue("single")
		postIdStr := r.FormValue("postId")

		isAuthed, username := IsLoggedIn(r)

		if !isAuthed {
			if single == "true" {
				http.Redirect(w, r, "/posts/"+postIdStr, http.StatusSeeOther)
				return
			}
			http.Redirect(w, r, "/posts", http.StatusSeeOther)
			return
		}

		var user models.User
		if res, err := database.GetUserDataStruct(sqliteDatabase, username); err != nil {
			errorHandler(w, r, 500)
			return
		} else {
			user = res
		}

		likeStatus, _ := strconv.Atoi(r.FormValue("likestatus"))
		commentIdStr := r.FormValue("commentId")
		commentId, _ := strconv.Atoi(commentIdStr)
		reaction := models.CommentLikes{
			CommentId:   commentId,
			AuthorId:    user.Id,
			LikeDislike: likeStatus,
		}

		var userReactions models.CommentLikes
		if res, err := database.GetUserCommentLikes(sqliteDatabase, commentId, user.Id); err != nil {
			errorHandler(w, r, 500)
		} else {
			userReactions = res
		}

		switch userReactions.LikeDislike {
		case 0:
			database.InsertCommentLike(sqliteDatabase, reaction)
		case 1:
			if likeStatus == 1 {
				database.DeleteCommentLike(sqliteDatabase, userReactions)
			} else if likeStatus == 2 {
				userReactions.LikeDislike = likeStatus
				database.UpdateCommentLike(sqliteDatabase, userReactions)
			}
		case 2:
			if likeStatus == 2 {
				database.DeleteCommentLike(sqliteDatabase, userReactions)
			} else if likeStatus == 1 {
				userReactions.LikeDislike = likeStatus
				database.UpdateCommentLike(sqliteDatabase, userReactions)
			}
		}

		if single == "true" {
			http.Redirect(w, r, "/posts/"+postIdStr, http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/posts", http.StatusSeeOther)
	}
}

func showLikedHandler(w http.ResponseWriter, r *http.Request) {
	_, username := IsLoggedIn(r)

	var user models.User
	if res, err := database.GetUserDataStruct(sqliteDatabase, username); err != nil {
		errorHandler(w, r, 500)
		return
	} else {
		user = res
	}

	database.GetUserLikedPosts(sqliteDatabase, user.Id)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	isAuthed, _ := IsLoggedIn(r)

	if isAuthed {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet && !justRegistred {
		tpl.ExecuteTemplate(w, "login.html", nil)
		return
	} else if r.Method == http.MethodGet && justRegistred {

		sID := uuid.NewV4()
		cookie := &http.Cookie{
			Name:   "session",
			Value:  sID.String(),
			MaxAge: 600,
		}
		http.SetCookie(w, cookie)
		if err := database.NewSession(sqliteDatabase, cookie.Value, regUserName); err != nil {
			errorHandler(w, r, 500)
			return
		}

		justRegistred = false
		regUserName = ""
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var user models.User
		if res, err := database.GetUserDataStruct(sqliteDatabase, username); err != nil {
			errorHandler(w, r, http.StatusBadRequest)
			return
		} else {
			user = res
		}

		err := bcrypt.CompareHashAndPassword([]byte(user.PassWord), []byte(password))
		if err != nil {
			tpl.ExecuteTemplate(w, "login.html", "Incorrect username and/or password")
			return
		}

		if err := database.DeleteSessionByUserName(sqliteDatabase, username); err != nil {
			errorHandler(w, r, 500)
			return
		}

		sID := uuid.NewV4()
		cookie := &http.Cookie{
			Name:   "session",
			Value:  sID.String(),
			MaxAge: 600,
		}
		http.SetCookie(w, cookie)

		if err := database.NewSession(sqliteDatabase, cookie.Value, username); err != nil {
			errorHandler(w, r, 500)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// remove cookie from db
	cookie, _ := r.Cookie("session")
	if err := database.DeleteSession(sqliteDatabase, cookie); err != nil {
		errorHandler(w, r, 500)
	}

	// remove cookie from browser
	cookie = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {

	if status == http.StatusForbidden {
		http.Error(w, "403 Forbidden", http.StatusForbidden)
	}
	if status == http.StatusNotFound {
		http.Error(w, "404 Page not found", 404)
	}
	if status == http.StatusInternalServerError {
		http.Error(w, "500 Internal server error", 500)
	}
	if status == http.StatusBadRequest {
		http.Error(w, "400 Bad request", 400)
	}
}
