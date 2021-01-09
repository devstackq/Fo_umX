package controllers

import (
	"ForumX/models"
	"ForumX/utils"
	"net/http"
)

//VotesPost func Post
func VotesPost(w http.ResponseWriter, r *http.Request) {

	if utils.URLChecker(w, r, "/votes/post") {

		pid := r.URL.Query().Get("id")

		if r.Method == "POST" {

			if r.FormValue("like") == "1" {
				models.VoteLike(w, r, pid, "post", session)
			}
			if r.FormValue("dislike") == "1" {
				models.VoteDislike(w, r, pid, "post", session)
			}
		}
		http.Redirect(w, r, "/post?id="+pid, 302)
	}
}

//VotesComment function
func VotesComment(w http.ResponseWriter, r *http.Request) {

	if utils.URLChecker(w, r, "/votes/comment") {

		commentID := r.URL.Query().Get("commentID")

		if r.Method == "POST" {
			if r.FormValue("commentLike") == "1" {
				models.VoteLike(w, r, commentID, "comment", session)
			}
			if r.FormValue("commentDislike") == "1" {
				models.VoteDislike(w, r, commentID, "comment", session)
			}
			http.Redirect(w, r, "/post?id="+r.FormValue("pidc"), 302)
		}
	}
}
