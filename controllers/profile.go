package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/devstackq/ForumX/models"
	util "github.com/devstackq/ForumX/utils"
)

//GetUserProfile  current -> user page
func GetUserProfile(w http.ResponseWriter, r *http.Request) {

	if util.URLChecker(w, r, "/profile") {

		if r.Method == "GET" {
			access, _ := util.IsCookie(w, r)
			if !access {
				http.Redirect(w, r, "/signin", 200)
				return
			}
			cookie, _ := r.Cookie("_cookie")
			//if userId now, createdPost uid equal -> show
			likedPost, posts, comments, user, err := models.GetUserProfile(r, w, cookie)
			if err != nil {
				log.Println(err)
			}

			//check if current cookie equal - cookie
			util.DisplayTemplate(w, "header", util.IsAuth(r))
			util.DisplayTemplate(w, "profile", user)
			util.DisplayTemplate(w, "favorited_post", likedPost)
			util.DisplayTemplate(w, "created_post", posts)
			util.DisplayTemplate(w, "comment_user", comments)

			//delete coookie db
			go func() {
				for now := range time.Tick(299 * time.Minute) {
					util.IsCookieExpiration(now, cookie, w, r)
					//next logout each 300 min
					time.Sleep(299 * time.Minute)
				}
			}()
		}
	}
}

//GetAnotherProfile  other user page
func GetAnotherProfile(w http.ResponseWriter, r *http.Request) {

	if util.URLChecker(w, r, "/user/id/") {

		if r.Method == "POST" {

			uid := models.User{Temp: r.FormValue("uid")}
			posts, user, err := uid.GetAnotherProfile(r)
			if err != nil {
				log.Println(err)
			}
			util.DisplayTemplate(w, "header", util.IsAuth(r))
			util.DisplayTemplate(w, "another_user", user)
			util.DisplayTemplate(w, "created_post", posts)
		}
	}
}

//UpdateProfile function
func UpdateProfile(w http.ResponseWriter, r *http.Request) {

	if util.URLChecker(w, r, "/edit/user") {

		if r.Method == "GET" {
			util.DisplayTemplate(w, "header", util.IsAuth(r))
			util.DisplayTemplate(w, "profile_update", "")
		}

		if r.Method == "POST" {

			access, s := util.IsCookie(w, r)
			if !access {
				return
			}

			DB.QueryRow("SELECT user_id FROM session WHERE uuid = ?", s.UUID).
				Scan(&s.UserID)

			is, _ := strconv.Atoi(r.FormValue("age"))

			p := models.User{
				FullName: r.FormValue("fullname"),
				Age:      is,
				Sex:      r.FormValue("sex"),
				City:     r.FormValue("city"),
				Image:    util.FileByte(r, "user"),
				ID:       s.UserID,
			}

			err = p.UpdateProfile()

			if err != nil {
				log.Println(err.Error())
			}
		}
		http.Redirect(w, r, "/profile", http.StatusFound)
	}
}

func DeleteAccount(w http.ResponseWriter, r *http.Request) {

	if util.URLChecker(w, r, "/delete/account") {

		if r.Method == "POST" {

			access, _ := util.IsCookie(w, r)
			if !access {
				return
			}
			var p models.User

			err := json.NewDecoder(r.Body).Decode(&p.ID)
			if err != nil {
				log.Println(err)
			}

			p.DeleteAccount(w, r)
			fmt.Println("delete account by ID", p.ID)
		}
		http.Redirect(w, r, "/", 302)
	}
}
