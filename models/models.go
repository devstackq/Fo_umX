package models

import (
	"database/sql"
	"encoding/base64"
	"log"
	"net/http"
	"time"
)

var (
	err error
	DB  *sql.DB
)

type Users struct {
	ID          int
	FullName    string
	Email       string
	Password    string
	IsAdmin     bool
	Age         int
	Sex         string
	CreatedTime time.Time
	City        string
	Image       []byte
	ImageHtml   string
	Role        string
	SVG         bool
}

type Category struct {
	ID     int
	Name   string
	UserID int
}

type Posts struct {
	ID            int
	Title         string
	Content       string
	CreatorID     int
	CategoryID    int
	CreationTime  time.Time
	EndpointPost  string
	FullName      string
	CategoryName  string
	Image         []byte
	ImageHtml     string
	PostIDEdit    int
	AuthorForPost int
	CountLike     int
	CountDislike  int
	SVG           bool
}

//comment ID -> foreign key -> postID
type Comments struct {
	ID             int
	Commentik      string
	PostID         int
	UserID         int
	CreatedTime    time.Time
	AuthorComment  string
	CommentLike    int
	CommentDislike int
}

var API struct {
	Authenticated bool
	Msg           string
}

//save session, by client cookie
type Session struct {
	ID     int
	UUID   string
	UserID int
}

type Likes struct {
	ID      int
	Like    int
	Dislike int
	PostID  int
	UserID  int
	Voted   bool
}

//link to COmments struct, then call func(r), return arr comments, post, err
func GetPostById(r *http.Request) ([]Comments, Posts, error) {

	r.ParseForm()
	id := r.FormValue("id")
	p := Posts{}

	//take from all post, only post by id, then write data struct Post
	DB.QueryRow("SELECT * FROM posts WHERE  id = ?", id).Scan(&p.ID, &p.Title, &p.Content, &p.CreatorID, &p.CategoryID, &p.CreationTime, &p.Image, &p.CountLike, &p.CountDislike)

	//write values from tables Likes, and write data table Post fileds like, dislikes
	//[]byte -> encode string, client render img base64

	if len(p.Image) > 0 {
		if p.Image[0] == 60 {
			p.SVG = true
		}
	}

	encodedString := base64.StdEncoding.EncodeToString(p.Image)
	p.ImageHtml = encodedString

	//creator post
	DB.QueryRow("SELECT full_name FROM users WHERE id = ?", p.CreatorID).Scan(&p.FullName)
	//cateogry post
	DB.QueryRow("SELECT title FROM categories WHERE id=?", p.CategoryID).Scan(&p.CategoryName)
	//get all comments from post
	stmp, err := DB.Query("SELECT * FROM comments WHERE  post_id =?", p.ID)
	if err != nil {
		log.Fatal(err)
	}
	defer stmp.Close()
	//write each fileds inside Comment struct -> then  append Array Comments
	ComentsPost := []Comments{}

	for stmp.Next() {
		comment := Comments{}
		var id, postID, userID, comLike, comDislike int
		var content string
		var myTime time.Time
		err = stmp.Scan(&id, &content, &postID, &userID, &myTime, &comLike, &comDislike)
		if err != nil {
			panic(err.Error)
		}

		comment.ID = id
		comment.Commentik = content
		comment.PostID = postID
		comment.UserID = userID
		comment.CreatedTime = myTime
		comment.CommentLike = comLike
		comment.CommentDislike = comDislike

		DB.QueryRow("SELECT full_name FROM users WHERE id = ?", userID).Scan(&comment.AuthorComment)
		ComentsPost = append(ComentsPost, comment)
	}

	if err != nil {
		return nil, p, err
	}
	return ComentsPost, p, nil
}

//get data from client, put data in Handler, then models -> query db
func (c *Comments) LostComment() error {

	_, err := DB.Exec("INSERT INTO comments( content, post_id, user_idx) VALUES(?,?,?)",
		c.Commentik, c.PostID, c.UserID)
	if err != nil {
		panic(err.Error())
	}
	return nil
}

//create post
func (p *Posts) CreatePost() error {
	_, err := DB.Exec("INSERT INTO posts (title, content, creator_id, category_id, image) VALUES ( ?,?, ?, ?, ?)",
		p.Title, p.Content, p.CreatorID, p.CategoryID, p.Image)
	if err != nil {
		return err
	}
	return nil
}

//update post
func (p *Posts) UpdatePost() error {

	_, err := DB.Exec("UPDATE  posts SET title=?, content=?, category_id=?, image=? WHERE id =?",
		p.Title, p.Content, p.CategoryID, p.Image, p.PostIDEdit)

	if err != nil {
		return err
	}
	return nil
}

//delete post
func (p *Posts) DeletePost() error {
	_, err := DB.Exec("DELETE FROM  posts  WHERE id =?", p.PostIDEdit)
	if err != nil {
		return err
	}
	return nil
}

//update profile
func (u *Users) UpdateProfile() error {

	_, err := DB.Exec("UPDATE  users SET full_name=?, age=?, sex=?, city=?, image=? WHERE id =?",
		u.FullName, u.Age, u.Sex, u.City, u.Image, u.ID)
	if err != nil {
		return err
	}
	return nil
}

//get profile by id
func GetUserProfile(r *http.Request) ([]Posts, []Posts, []Comments, Users, error) {

	cookie, _ := r.Cookie("_cookie")

	s := Session{UUID: cookie.Value}
	u := Users{}
	DB.QueryRow("SELECT user_id FROM session WHERE uuid = ?", s.UUID).Scan(&s.UserID)
	lps := []Likes{}
	lp, err := DB.Query("select post_id from likes where user_id =?", s.UserID)

	for lp.Next() {
		l := Likes{}

		var lpid int

		err = lp.Scan(&lpid)
		l.PostID = lpid
		lps = append(lps, l)
	}

	DB.QueryRow("SELECT * FROM users WHERE id = ?", s.UserID).Scan(&u.ID, &u.FullName, &u.Email, &u.Password, &u.IsAdmin, &u.Age, &u.Sex, &u.CreatedTime, &u.City, &u.Image)

	if u.Image[0] == 60 {
		u.SVG = true
	}

	encStr := base64.StdEncoding.EncodeToString(u.Image)
	u.ImageHtml = encStr

	var likedpost *sql.Rows
	LikedPosts := []Posts{}
	var can []int

	for _, v := range lps {
		can = append(can, v.PostID)
	}

	//unique liked post by user
	fin := isUnique(can)

	//accum liked post
	for _, v := range fin {
		//get each liked post by ID, then likedpost, puth array post
		likedpost, err = DB.Query("SELECT * FROM posts WHERE id=?", v)

		for likedpost.Next() {

			post := Posts{}

			var id, creatorid, categoryid, countlike, countdislike int
			var content, title string
			var creationtime time.Time
			var image []byte

			err = likedpost.Scan(&id, &title, &content, &creatorid, &categoryid, &creationtime, &image, &countlike, &countdislike)
			if err != nil {
				panic(err.Error)
			}

			post.ID = id
			post.Title = title
			post.Content = content
			post.CreatorID = creatorid
			post.CategoryID = categoryid
			post.CreationTime = creationtime
			post.Image = image
			post.CountLike = countlike
			post.CountDislike = countdislike

			LikedPosts = append(LikedPosts, post)
		}
	}

	psu, err := DB.Query("SELECT * FROM posts WHERE creator_id=?", s.UserID)

	PostsCreatedUser := []Posts{}

	for psu.Next() {

		post := Posts{}

		var id, creatorid, categoryid, countlike, countdislike int
		var content, title string
		var creationtime time.Time
		var image []byte

		err = psu.Scan(&id, &title, &content, &creatorid, &categoryid, &creationtime, &image, &countlike, &countdislike)
		if err != nil {
			panic(err.Error)
		}

		post.AuthorForPost = s.UserID

		post.ID = id
		post.Title = title
		post.Content = content
		post.CreatorID = creatorid
		post.CategoryID = categoryid
		post.CreationTime = creationtime
		post.Image = image
		post.CountLike = countlike
		post.CountDislike = countdislike

		PostsCreatedUser = append(PostsCreatedUser, post)
	}

	csu, err := DB.Query("SELECT * FROM comments WHERE user_idx=?", s.UserID)

	CommentsLostUser := []Comments{}

	defer csu.Close()

	for csu.Next() {

		comment := Comments{}

		var id, postid, useridx, comLike, comDislike int
		var content string
		var createdtime time.Time

		err = csu.Scan(&id, &content, &postid, &useridx, &createdtime, &comLike, &comDislike)
		if err != nil {
			panic(err.Error)
		}

		comment.ID = id
		comment.PostID = postid
		comment.UserID = useridx
		comment.Commentik = content
		comment.CreatedTime = createdtime
		comment.CommentLike = comLike
		comment.CommentDislike = comDislike

		CommentsLostUser = append(CommentsLostUser, comment)
	}

	if err != nil {
		return nil, nil, nil, u, err
	}

	return LikedPosts, PostsCreatedUser, CommentsLostUser, u, nil
}

//find unique liked post
func isUnique(intSlice []int) []int {
	keys := make(map[int]bool)
	list := []int{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

//get other user
func GetOtherUser(r *http.Request) ([]Posts, Users, error) {

	r.ParseForm()
	uid := r.FormValue("uid")

	user := DB.QueryRow("SELECT * FROM users WHERE id = ?", uid)
	u := Users{}
	err = user.Scan(&u.ID, &u.FullName, &u.Email, &u.Password, &u.IsAdmin, &u.Age, &u.Sex, &u.CreatedTime, &u.City, &u.Image)
	if u.Image[0] == 60 {
		u.SVG = true
	}
	encStr := base64.StdEncoding.EncodeToString(u.Image)
	u.ImageHtml = encStr
	psu, err := DB.Query("SELECT * FROM posts WHERE creator_id=?", u.ID)

	PostsOtherUser := []Posts{}

	defer psu.Close()

	for psu.Next() {
		post := Posts{}

		var id, creatorid, categoryid, countlike, countdislike int
		var content, title string
		var creationtime time.Time
		var image []byte

		err = psu.Scan(&id, &title, &content, &creatorid, &categoryid, &creationtime, &image, &countlike, &countdislike)
		if err != nil {
			panic(err.Error)
		}

		post.ID = id
		post.Title = title
		post.Content = content
		post.CreatorID = creatorid
		post.CategoryID = categoryid
		post.CreationTime = creationtime
		post.CountLike = countlike
		post.CountDislike = countdislike
		PostsOtherUser = append(PostsOtherUser, post)
	}
	if err != nil {
		return nil, u, err
	}
	return PostsOtherUser, u, nil
}

//search
func Search(w http.ResponseWriter, r *http.Request) ([]Posts, error) {

	keyword := r.FormValue("search")
	psu, err := DB.Query("SELECT * FROM posts WHERE title LIKE ?", "%"+keyword+"%")
	PostsUser := []Posts{}

	defer psu.Close()

	for psu.Next() {

		post := Posts{}

		var id, creatorid, categoryid, countlike, countdislike int
		var content, title string
		var creationtime time.Time
		var image []byte

		err = psu.Scan(&id, &title, &content, &creatorid, &categoryid, &creationtime, &image, &countlike, &countdislike)

		if err != nil {
			panic(err.Error)
		}
		post.ID = id
		post.Title = title
		post.Content = content
		post.CreatorID = creatorid
		post.CategoryID = categoryid
		post.CreationTime = creationtime
		post.Image = image
		post.CountLike = countlike
		post.CountDislike = countdislike

		PostsUser = append(PostsUser, post)
	}

	if err != nil {
		return nil, err
	}
	return PostsUser, nil
}
