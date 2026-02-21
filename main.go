package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
    _ "github.com/mattn/go-sqlite3"
)

var (
	tpl            *template.Template
	sqliteDatabase *sql.DB
)

func init() {
	ReCreateAndConnSqlDataBase()
	tpl = template.Must(template.ParseGlob("static/templates/*"))
}

func main() {
	tpl, _ = template.ParseGlob("static/templates/*.html")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/posts/", makeHandler(postsMiddleware))
	http.HandleFunc("/categories/", makeHandler(categoryHandler))
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/comments", commentsHandler)
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/reactpost", reactPostHandler)
	http.HandleFunc("/reactcomment", reactCommentHandler)
	http.HandleFunc("/showliked", showLikedHandler)

	fmt.Printf("Starting server at port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
