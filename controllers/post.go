package controllers

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strconv"

	structure "github.com/devstackq/ForumX/general"
	"github.com/devstackq/ForumX/models"
	util "github.com/devstackq/ForumX/utils"
)

var (
	err  error
	DB   *sql.DB
	msg  = structure.API.Message
	auth = structure.API.Authenticated
)

//GetAllPosts  by category || all posts
func GetAllPosts(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" && r.URL.Path != "/science" && r.URL.Path != "/love" && r.URL.Path != "/sapid" {
		util.DisplayTemplate(w, "404page", http.StatusNotFound)
		return
	}

	filterValue := models.Filter{
		Like:     r.FormValue("likes"),
		Date:     r.FormValue("date"),
		Category: r.FormValue("cats"),
	}

	posts, endpoint, category, err := filterValue.GetAllPost(r, r.FormValue("next"), r.FormValue("prev"))

	if err != nil {
		log.Fatal(err)
	}

	util.DisplayTemplate(w, "header", util.IsAuth(r))

	if endpoint == "/" {
		util.DisplayTemplate(w, "index", posts)
	} else {
		//send category value
		msg := []byte(fmt.Sprintf("<h3 id='category'> %s </h3>", category))
		w.Header().Set("Content-Type", "application/json")
		w.Write(msg)
		util.DisplayTemplate(w, "category_post_template", posts)
	}
}

//GetPostByID  1 post by id
func GetPostByID(w http.ResponseWriter, r *http.Request) {

	if util.URLChecker(w, r, "/post") {

		id, _ := strconv.Atoi(r.FormValue("id"))
		pid := models.Post{ID: id}
		comments, post, err := pid.GetPostByID(r)

		if err != nil {
			log.Println(err)
		}
		util.DisplayTemplate(w, "header", util.IsAuth(r))
		util.DisplayTemplate(w, "posts", post)
		util.DisplayTemplate(w, "comment_post", comments)
	}
}

//CreatePost  function
func CreatePost(w http.ResponseWriter, r *http.Request) {

	if util.URLChecker(w, r, "/create/post") {

		//switch r.Method {
		if r.Method == "GET" {
			util.DisplayTemplate(w, "header", util.IsAuth(r))
			util.DisplayTemplate(w, "create_post", &msg)
		}

		if r.Method == "POST" {
			access, session := util.IsCookie(w, r)
			log.Println(access, "access status")
			if !access {
				http.Redirect(w, r, "/signin", 200)
				return
			}
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

	if util.URLChecker(w, r, "/edit/post") {

		pid, _ := strconv.Atoi(r.FormValue("id"))

		if r.Method == "GET" {

			var p models.Post
			DB.QueryRow("SELECT * FROM posts WHERE id = ?", pid).Scan(&p.ID, &p.Title, &p.Content, &p.CreatorID, &p.CreatedTime, &p.Image, &p.Like, &p.Dislike)
			p.ImageHTML = base64.StdEncoding.EncodeToString(p.Image)

			util.DisplayTemplate(w, "header", util.IsAuth(r))
			util.DisplayTemplate(w, "update_post", p)

		}

		if r.Method == "POST" {

			access, _ := util.IsCookie(w, r)
			if !access {
				http.Redirect(w, r, "/signin", 200)
				return
			}

			p := models.Post{
				Title:   r.FormValue("title"),
				Content: r.FormValue("content"),
				Image:   util.IsImage(r),
				ID:      pid,
			}

			err = p.UpdatePost()

			if err != nil {
				//try hadnler all error
				defer log.Println(err, "upd post err")
			}
		}
		http.Redirect(w, r, "/post?id="+strconv.Itoa(int(pid)), 302)

	}
}

//DeletePost function
func DeletePost(w http.ResponseWriter, r *http.Request) {

	if util.URLChecker(w, r, "/delete/post") {

		access, _ := util.IsCookie(w, r)
		if !access {
			http.Redirect(w, r, "/signin", 200)
			return
		}
		pid, _ := strconv.Atoi(r.URL.Query().Get("id"))
		p := models.Post{ID: pid}

		err = p.DeletePost()

		if err != nil {
			log.Println(err.Error())
		}
		http.Redirect(w, r, "/", 302)
	}
}

//Search
func Search(w http.ResponseWriter, r *http.Request) {

	if util.URLChecker(w, r, "/search") {

		if r.Method == "GET" {
			util.DisplayTemplate(w, "search", http.StatusFound)
		}

		if r.Method == "POST" {

			foundPosts, err := models.Search(w, r)

			if err != nil {
				log.Println(err)
			}
			if foundPosts == nil {
				util.DisplayTemplate(w, "header", util.IsAuth(r))
				msg := []byte(fmt.Sprintf("<h2 id='notFound'> Nihuya ne naideno </h2>"))
				w.Header().Set("Content-Type", "application/json")
				w.Write(msg)
				util.DisplayTemplate(w, "index", nil)
			} else {
				util.DisplayTemplate(w, "header", util.IsAuth(r))
				util.DisplayTemplate(w, "index", foundPosts)
			}
		}
	}
}
