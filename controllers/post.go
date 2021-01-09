package controllers

import (
	"ForumX/general"
	"ForumX/models"
	"ForumX/utils"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	err  error
	DB   *sql.DB
	msg  = general.API.Message
	auth = general.API.Authenticated
)

//GetAllPosts  by category || all posts
func GetAllPosts(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" && r.URL.Path != "/science" && r.URL.Path != "/love" && r.URL.Path != "/sapid" {
		utils.RenderTemplate(w, "404page", http.StatusNotFound)
		return
	}
	//methods Filter struct -> if no case filter : default get All post
	filterValue := models.Filter{
		Like:     r.FormValue("likes"),
		Date:     r.FormValue("date"),
		Category: r.FormValue("cats"),
	}

	posts, endpoint, category := filterValue.GetAllPost(r, r.FormValue("next"), r.FormValue("prev"))
	utils.RenderTemplate(w, "header", utils.IsAuth(r))
	if posts == nil {
		msg := []byte(fmt.Sprintf("<span id='notify-post'> Post nil </span>", ))
		w.Header().Set("Content-Type", "application/json")
		w.Write(msg)
		//utils.RenderTemplate(w, "category_post_template", posts)
		return
	}
	if endpoint == "/" {
		utils.RenderTemplate(w, "index", posts)
	} else {
		//send category value
		msg := []byte(fmt.Sprintf("<span id='category'> %s </span>", category))
		w.Header().Set("Content-Type", "application/json")
		w.Write(msg)
		utils.RenderTemplate(w, "category_post_template", posts)
	}
}

//GetPostByID  1 post by id
func GetPostByID(w http.ResponseWriter, r *http.Request) {

	if utils.URLChecker(w, r, "/post") {

		id, _ := strconv.Atoi(r.FormValue("id"))
		var temp string
		err = DB.QueryRow("select content from posts where id=?", id).Scan(&temp)
		if err != nil {
			log.Println(err)
			utils.RenderTemplate(w, "404page", 404)
			return
		}
		//fmt.Println(temp)
		pid := models.Post{ID: id}
		comments, post := pid.GetPostByID(r)

		utils.RenderTemplate(w, "header", utils.IsAuth(r))
		utils.RenderTemplate(w, "posts", post)
		utils.RenderTemplate(w, "comment_post", comments)
		//utils.RenderTemplate(w, "reply_comment", repliesComment)
	}
}

//CreatePost  function
func CreatePost(w http.ResponseWriter, r *http.Request) {

	if utils.URLChecker(w, r, "/create/post") {

		if r.Method == "GET" {
			utils.RenderTemplate(w, "header", utils.IsAuth(r))
			utils.RenderTemplate(w, "create_post", &msg)
		}

		if r.Method == "POST" {
			//r.ParseMultipartForm(10 << 20)
			f, _, _ := r.FormFile("uploadfile")
			f2, _, _ := r.FormFile("uploadfile")

			categories, _ := r.Form["input"]

			photoFlag := false
			if f != nil && f2 != nil {
				photoFlag = true
			}
			post := models.Post{
				Title:      r.FormValue("title"),
				Content:    r.FormValue("content"),
				Categories: categories,
				FileS:      f,
				FileI:      f2,
				Session:    session,
				IsPhoto:    photoFlag,
			}
			post.CreatePost(w, r)
		}
	}
}

//UpdatePost function
func UpdatePost(w http.ResponseWriter, r *http.Request) {

	if utils.URLChecker(w, r, "/edit/post") {

		pid, _ := strconv.Atoi(r.FormValue("id"))

		if r.Method == "GET" {
			//send data - client
			var p models.Post
			DB.QueryRow("SELECT * FROM posts WHERE id = ?", pid).Scan(&p.ID, &p.Title, &p.Content, &p.CreatorID, &p.CreateTime, &p.UpdateTime, &p.Image, &p.Like, &p.Dislike)
			p.ImageHTML = base64.StdEncoding.EncodeToString(p.Image)

			utils.RenderTemplate(w, "header", utils.IsAuth(r))
			utils.RenderTemplate(w, "update_post", p)
		}

		if r.Method == "POST" {

			p := models.Post {
				Title:      r.FormValue("title"),
				Content:    r.FormValue("content"),
				Image:      utils.IsImage(r),
				ID:         pid,
				UpdateTime: time.Now(),
			}
			p.UpdatePost()
			http.Redirect(w, r, "/post?id="+strconv.Itoa(int(pid)), 302)
		}
	}
}

//DeletePost function
func DeletePost(w http.ResponseWriter, r *http.Request) {

	if utils.URLChecker(w, r, "/delete/post") {
		pid, _ := strconv.Atoi(r.URL.Query().Get("id"))
		p := models.Post{ID: pid}
		p.DeletePost()
		http.Redirect(w, r, "/profile", 302)
	}
}

//Search
func Search(w http.ResponseWriter, r *http.Request) {

	if utils.URLChecker(w, r, "/search") {
		foundPosts := models.Search(w, r)
		utils.RenderTemplate(w, "header", utils.IsAuth(r))
		if foundPosts == nil {
			msg := []byte(fmt.Sprintf("<h2 id='notFound'> Not found</h2>"))
			w.Header().Set("Content-Type", "application/json")
			w.Write(msg)
			utils.RenderTemplate(w, "index", nil)
		} else {
			utils.RenderTemplate(w, "index", foundPosts)
		}
	}
}
