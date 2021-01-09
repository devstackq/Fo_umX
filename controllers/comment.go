package controllers

import (
	"ForumX/models"
	"ForumX/utils"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

//LeaveComment function
func LeaveComment(w http.ResponseWriter, r *http.Request) {

	if utils.URLChecker(w, r, "/comment") {

		commentInput := r.FormValue("comment-text")

		if utils.IsValidLetter(commentInput, "post") {
			comment := models.Comment{
				Content: commentInput,
				PostID:  r.FormValue("curr"),
				UserID:  session.UserID,
			}
			comment.LeaveComment()
		}
		http.Redirect(w, r, "/post?id="+r.FormValue("curr"), 302)
	}
}

//UpdateComment func
func UpdateComment(w http.ResponseWriter, r *http.Request) {

	if utils.URLChecker(w, r, "/edit/comment") {

		cid, _ := strconv.Atoi(r.FormValue("id"))

		if r.Method == "GET" {

			var comment models.Comment
			err = DB.QueryRow("SELECT * FROM comments WHERE id = ?", cid).Scan(&comment.ID, &comment.ParentID, &comment.Content, &comment.PostID, &comment.UserID, &comment.ToWhom, &comment.FromWhom, &comment.CreatedTime, &comment.UpdatedTime, &comment.Like, &comment.Dislike)
			if err != nil {
				fmt.Println(err)
			}

			utils.RenderTemplate(w, "header", utils.IsAuth(r))
			utils.RenderTemplate(w, "update_comment", comment)
		}
		if r.Method == "POST" {

			comment := models.Comment{
				ID:          cid,
				Content:     r.FormValue("content"),
				UpdatedTime: time.Now(),
			}

			comment.UpdateComment()
		http.Redirect(w, r, "/profile", 302)

		}
	}
}

//DeleteComment dsa
func DeleteComment(w http.ResponseWriter, r *http.Request) {

	if utils.URLChecker(w, r, "/delete/comment") {
		models.DeleteComment(r.FormValue("id"))
	}
	http.Redirect(w, r, "/profile", 302)
}

// if have parentId-, when answer another comment, getByIdComment ->  get All data, then, show
//message cleint side, fromWho, toWhom answer
func ReplyComment(w http.ResponseWriter, r *http.Request) {

	if utils.URLChecker(w, r, "/reply/comment") {

		content := r.FormValue("answerComment")
		parent, _ := strconv.Atoi(r.FormValue("parentId"))
		postId := r.FormValue("postId")
		comment := models.Comment{}

		if utils.IsValidLetter(content, "post") {
			//get creator_id -> when answer by post, else get fromWho, reply another comment
		err = DB.QueryRow("SELECT creator_id FROM comments WHERE id = ?", parent).Scan(&comment.ToWhom)
		if err != nil {
			log.Println(err, "not find user")
		}

		fmt.Println(parent, "parentId", parent, "replyID", parent, "cid", postId, "pid")

			comment.CommentID= parent
			comment.Content=   content
			comment.PostID=   postId
			comment.FromWhom= session.UserID
			comment.ParentID = parent
			comment.UserID = session.UserID
			
			commentPrepare, err := DB.Prepare(`INSERT INTO comments(parent_id, content, post_id, creator_id, toWho, fromWho, create_time) VALUES(?,?,?,?,?,?,?)`)
			if err != nil {
				log.Println(err)
			}
			_, err = commentPrepare.Exec(comment.ParentID, comment.Content, comment.PostID, comment.UserID, comment.ToWhom,  comment.FromWhom, time.Now())
			if err != nil {
				log.Println(err)
			}
		}
		http.Redirect(w, r, "/post?id="+postId, 302)
	}
}
