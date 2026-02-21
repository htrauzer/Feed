package models

type Category struct {
	Id   int
	Name string
}

type Comment struct {
	Id             int
	Content        string
	Author         string
	PostId         int
	CreationTime   string
	Likes          int
	Dislikes       int
	UserLikeStatus int
}

type CommentLikes struct {
	Id          int
	CommentId   int
	AuthorId    int
	LikeDislike int
}

type Post struct {
	Id             int
	Content        string
	AuthorId       int // user.Id
	UserName       string
	Comments       []Comment
	Categories     []Category
	CreationTime   string
	Likes          int
	Dislikes       int
	UserLikeStatus int
	Title          string
	NumComments    int
}

type PostCategory struct {
	Id       int
	Post     *Post
	Category string //Category.Name
}

type PostLikes struct {
	Id          int
	PostId      int
	AuthorId    int
	LikeDislike int
}

type User struct {
	Id       int
	Email    string
	UserName string
	PassWord string
}
