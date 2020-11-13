package models

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	structure "github.com/devstackq/ForumX/general"
	util "github.com/devstackq/ForumX/utils"
)

//global variable for package models
var (
	err                          error
	DB                           *sql.DB
	rows                         *sql.Rows
	id, creatorID, like, dislike int
	content, title               string
	createdTime                  time.Time
	image                        []byte
	postID                       int
	userID                       int
	post                         Post
	comment                      Comment
	msg                          = structure.API.Message
	pageNum                      = 1
)

//Posts struct
type Post struct {
	ID            int
	Title         string
	Content       string
	CreatorID     int
	CreatedTime   time.Time
	Endpoint      string
	FullName      string
	Image         []byte
	ImageHTML     string
	PostIDEdit    int
	AuthorForPost int
	Like          int
	Dislike       int
	SVG           bool
	PBGID         int
	PBGPostID     int
	PBGCategory   string
	FileS         multipart.File
	FileI         multipart.File
	Session       structure.Session
	Categories    []string
	Temp          string
	IsPhoto       bool
	Time          string
	CountPost     int
}

//PostCategory struct
type PostCategory struct {
	PostID   int64
	Category string
}

//Filter struct
type Filter struct {
	Category string
	Like     string
	Date     string
}

//GetAllPost function
func (f *Filter) GetAllPost(r *http.Request, next, prev string) ([]Post, string, string, error) {
	//pageNum = 1
	var post Post
	var leftJoin bool
	var arrPosts []Post

	//each call +1
	if next == "next" {
		pageNum++
	}
	if prev == "prev" {
		pageNum--
	}
	//count pageNum, fix
	//fmt.Print(pageNum)

	limit := 4
	offset := limit * (pageNum - 1)

	switch r.URL.Path {
	case "/":
		leftJoin = false
		post.Endpoint = "/"
		if f.Date == "asc" {
			rows, err = DB.Query("SELECT * FROM posts  ORDER BY created_time ASC LIMIT 8 ")
		} else if f.Date == "desc" {
			rows, err = DB.Query("SELECT * FROM posts  ORDER BY created_time DESC LIMIT 8")
		} else if f.Like == "like" {
			rows, err = DB.Query("SELECT * FROM posts  ORDER BY count_like DESC LIMIT 8")
		} else if f.Like == "dislike" {
			rows, err = DB.Query("SELECT * FROM posts  ORDER BY count_dislike DESC LIMIT 8")
		} else if f.Category != "" {
			leftJoin = true
			rows, err = DB.Query("SELECT  * FROM posts  LEFT JOIN post_cat_bridge  ON post_cat_bridge.post_id = posts.id   WHERE category=? ORDER  BY created_time  DESC LIMIT 8", f.Category)
		} else {
			rows, err = DB.Query("SELECT * FROM posts ORDER BY created_time DESC LIMIT ? OFFSET ?", limit, offset)
		}

	case "/science":
		leftJoin = true
		post.Temp = "Science"
		post.Endpoint = "/science"
		rows, err = DB.Query("SELECT * FROM posts  LEFT JOIN post_cat_bridge  ON post_cat_bridge.post_id = posts.id   WHERE category=?  ORDER  BY created_time  DESC LIMIT 5", "science")
	case "/love":
		leftJoin = true
		post.Temp = "Love"
		post.Endpoint = "/love"
		rows, err = DB.Query("SELECT  * FROM posts  LEFT JOIN post_cat_bridge  ON post_cat_bridge.post_id = posts.id  WHERE category=?   ORDER  BY created_time  DESC LIMIT 5", "love")
	case "/sapid":
		leftJoin = true
		post.Temp = "Sapid"
		post.Endpoint = "/sapid"
		rows, err = DB.Query("SELECT  * FROM posts  LEFT JOIN post_cat_bridge  ON post_cat_bridge.post_id = posts.id  WHERE category=?  ORDER  BY created_time  DESC LIMIT 5", "sapid")
	}

	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	for rows.Next() {
		if leftJoin {
			if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatorID, &post.CreatedTime, &post.Image, &post.Like, &post.Dislike, &post.PBGID, &post.PBGPostID, &post.PBGCategory); err != nil {
				fmt.Println(err)
			}
		} else {
			if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatorID, &post.CreatedTime, &post.Image, &post.Like, &post.Dislike); err != nil {
				fmt.Println(err)
			}
			//fmt.Print(post.ID)
		}

		if err != nil {
			log.Println(err)
		}
		//send countr +1
		err = DB.QueryRow("SELECT COUNT(id) FROM posts").Scan(&post.CountPost)
		post.Time = post.CreatedTime.Format("2006 Jan _2 15:04:05")
		arrPosts = append(arrPosts, post)
	}
	//err = DB.QueryRow("SELECT COUNT(id) FROM posts").Scan(&post.CountPost)
	return arrPosts, post.Endpoint, post.Temp, nil
}

//UpdatePost fucntion
func (p *Post) UpdatePost() error {

	_, err := DB.Exec("UPDATE  posts SET title=?, content=?, image=? WHERE id =?",
		p.Title, p.Content, p.Image, p.ID)

	if err != nil {
		return err
	}
	return nil
}

//DeletePost function
func (p *Post) DeletePost() error {
	_, err := DB.Exec("DELETE FROM  posts  WHERE id =?", p.ID)
	if err != nil {
		return err
	}
	return nil
}

//CreatePost function
func (p *Post) CreatePost(w http.ResponseWriter, r *http.Request) {

	var fileBytes []byte
	var buff bytes.Buffer

	if p.IsPhoto {

		fileSize, _ := buff.ReadFrom(p.FileS)
		defer p.FileS.Close()

		if fileSize < 20000000 {
			if err != nil {
				log.Fatal(err)
			}
			fileBytes, err = ioutil.ReadAll(p.FileI)
		} else {
			util.DisplayTemplate(w, "header", util.IsAuth(r))
			util.DisplayTemplate(w, "create", "Large file, more than 20mb")
		}
	} else {
		//set empty photo post
		fileBytes = []byte{0, 0}
	}

	DB.QueryRow("SELECT user_id FROM session WHERE uuid = ?", p.Session.UUID).Scan(&p.Session.UserID)

	//check empty values
	if util.CheckLetter(p.Title) && util.CheckLetter(p.Content) {

		db, err := DB.Exec("INSERT INTO posts (title, content, creator_id, created_time, image) VALUES ( ?,?, ?, ?, ?)",
			p.Title, p.Content, p.Session.UserID, time.Now(), fileBytes)
		if err != nil {
			log.Println(err)
		}
		last, err := db.LastInsertId()
		if err != nil {
			log.Println(err)
		}

		if len(p.Categories) == 1 {
			pcb := PostCategory{
				PostID:   last,
				Category: p.Categories[0],
			}
			pcb.CreateBridge()

		} else if len(p.Categories) > 1 {
			//loop add > 1 category post
			for _, v := range p.Categories {
				pcb := PostCategory{
					PostID:   last,
					Category: v,
				}
				pcb.CreateBridge()
			}
		}
		s := strconv.Itoa(int(last))
		http.Redirect(w, r, "/post?id="+s, 302)

	} else {
		msg = "Empty title or content"
		util.DisplayTemplate(w, "header", util.IsAuth(r))
		util.DisplayTemplate(w, "create_post", &msg)
	}
}

//GetPostById function take from all post, only post by id, then write p struct Post
func (post *Post) GetPostByID(r *http.Request) ([]Comment, Post, error) {

	p := Post{}
	DB.QueryRow("SELECT * FROM posts WHERE id = ?", post.ID).Scan(&p.ID, &p.Title, &p.Content, &p.CreatorID, &p.CreatedTime, &p.Image, &p.Like, &p.Dislike)

	//[]byte -> encode string, client render img base64
	//check svg || jpg,png
	if len(p.Image) > 0 {
		if p.Image[0] == 60 {
			p.SVG = true
		}
	}
	p.Time = p.CreatedTime.Format("2006 Jan _2 15:04:05")

	p.ImageHTML = base64.StdEncoding.EncodeToString(p.Image)

	DB.QueryRow("SELECT full_name FROM users WHERE id = ?", p.CreatorID).Scan(&p.FullName)

	stmp, err := DB.Query("SELECT * FROM comments WHERE  post_id =?", p.ID)
	if err != nil {
		log.Fatal(err)
	}
	defer stmp.Close()
	//write each fields inside Comment struct -> then  append Array Comments
	var comments []Comment

	for stmp.Next() {

		comment := Comment{}
		err = stmp.Scan(&id, &content, &postID, &userID, &createdTime, &like, &dislike)
		if err != nil {
			log.Println(err.Error())
		}
		comment = AppendComment(id, content, postID, userID, createdTime, like, dislike, "")
		DB.QueryRow("SELECT full_name FROM users WHERE id = ?", userID).Scan(&comment.Author)
		comments = append(comments, comment)
	}

	if err != nil {
		return nil, p, err
	}
	return comments, p, nil
}

//CreateBridge create post  -> post_id + category
func (pcb *PostCategory) CreateBridge() {

	_, err := DB.Exec("INSERT INTO post_cat_bridge (post_id, category) VALUES (?, ?)",
		pcb.PostID, pcb.Category)

	if err != nil {
		log.Println(err)
		return
	}
}

//Search post by contain title
func Search(w http.ResponseWriter, r *http.Request) ([]Post, error) {

	var posts []Post
	psu, err := DB.Query("SELECT * FROM posts WHERE title LIKE ?", "%"+r.FormValue("search")+"%")
	defer psu.Close()

	for psu.Next() {

		err = psu.Scan(&id, &title, &content, &creatorID, &createdTime, &image, &like, &dislike)
		if err != nil {
			log.Println(err.Error())
		}
		post = AppendPost(id, title, content, creatorID, image, like, dislike, 0, createdTime)
		posts = append(posts, post)
	}

	if err != nil {
		return nil, err
	}
	return posts, nil
}

//appendPost each post put value from Db
func AppendPost(id int, title, content string, creatorID int, image []byte, like, dislike, authorID int, createdTime time.Time) Post {

	post = Post{
		ID:            id,
		Title:         title,
		Content:       content,
		CreatorID:     creatorID,
		Image:         image,
		Like:          like,
		Dislike:       dislike,
		AuthorForPost: authorID,
		Time:          createdTime.Format("2006 Jan _2 15:04:05"),
	}
	return post
}
