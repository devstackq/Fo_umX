package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/devstackq/ForumX/models"
	util "github.com/devstackq/ForumX/utils"
)

//LeaveComment function
func LeaveComment(w http.ResponseWriter, r *http.Request) {

	if util.URLChecker(w, r, "/comment") {

		if r.Method == "POST" {

			pid, _ := strconv.Atoi(r.FormValue("curr"))
			commentInput := r.FormValue("comment-text")

			access, s := util.IsCookie(w, r)
			if !access {
				return
			}

			DB.QueryRow("SELECT user_id FROM session WHERE uuid = ?", s.UUID).Scan(&s.UserID)

			if util.CheckLetter(commentInput) {

				comment := models.Comment{
					Content: commentInput,
					PostID:  pid,
					UserID:  s.UserID,
				}

				err = comment.LeaveComment()

				if err != nil {
					log.Println(err.Error())
				}
			}
		}
		http.Redirect(w, r, "/post?id="+r.FormValue("curr"), 302)
	}
}
