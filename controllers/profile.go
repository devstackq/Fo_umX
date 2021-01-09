package controllers

import (
	"ForumX/models"
	"ForumX/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

//GetUserProfile  current -> user page
func GetUserProfile(w http.ResponseWriter, r *http.Request) {

	if utils.URLChecker(w, r, "/profile") {

		if r.Method == "GET" {
			//if userId now, createdPost uid equal -> show
			u := models.User{
				Session : session,
			}
			dislikedPost, likedPost, posts, comments, user := u.GetUserProfile(r, w)
			//check if current cookie equal - cookie
			utils.RenderTemplate(w, "header", utils.IsAuth(r))
			utils.RenderTemplate(w, "profile", user)
			utils.RenderTemplate(w, "created_post", posts)
			utils.RenderTemplate(w, "favorited_post", likedPost)
			utils.RenderTemplate(w, "disliked_post", dislikedPost)
			utils.RenderTemplate(w, "comment_user", comments)
		}
	}
}

//GetAnotherProfile  other user page
func GetAnotherProfile(w http.ResponseWriter, r *http.Request) {

	if utils.URLChecker(w, r, "/user/id") {

		//if r.Method == "POST" {
			uid := models.User{Temp: r.FormValue("uid")}

			if uid.Temp == "0" {
				uid.Temp = r.URL.Query().Get("uid")
			}
			fmt.Println(uid, "UID")
			posts, user, err := uid.GetAnotherProfile(r)
			if err != nil {
				log.Println(err)
			}

			utils.RenderTemplate(w, "header", utils.IsAuth(r))
			utils.RenderTemplate(w, "another_user", user)
			utils.RenderTemplate(w, "created_post", posts)
	}
}

//GetUserActivities func
func GetUserActivities(w http.ResponseWriter, r *http.Request) {

	if utils.URLChecker(w, r, "/activity") {

		notifies := models.GetUserActivities(w, r, session)
		utils.RenderTemplate(w, "header", utils.IsAuth(r))
		if notifies != nil {
			
			utils.RenderTemplate(w, "activity", notifies)
		}else {
			msg := []byte(fmt.Sprintf("<h3> %s </h3>", "Empty notification"))
			w.Header().Set("Content-Type", "application/json")
			w.Write(msg)
		}
	}
}

//UpdateProfile function
func UpdateProfile(w http.ResponseWriter, r *http.Request) {

	if utils.URLChecker(w, r, "/edit/user") {

		if r.Method == "GET" {
			var u models.User

		err = DB.QueryRow("SELECT id, full_name, email, username, age, sex, city, image FROM users where id=?", session.UserID).Scan(&u.ID, &u.FullName, &u.Email, &u.Username, &u.Age, &u.Sex, &u.City, &u.Image)	
			if err != nil {
				log.Println(err)
			}
		u.ImageHTML = base64.StdEncoding.EncodeToString(u.Image)
			utils.RenderTemplate(w, "header", utils.IsAuth(r))
			utils.RenderTemplate(w, "profile_update", u)
		}

		if r.Method == "POST" {
			age, _ := strconv.Atoi(r.FormValue("age"))
		
			p := models.User{
				FullName: r.FormValue("fullname"),
				Age:      age,
				Username: r.FormValue("nickname"),
				Sex:      r.FormValue("sex"),
				City:     r.FormValue("city"),
				Image:    utils.FileByte(r, "user"),
				ID:       session.UserID,
			}

			p.UpdateProfile()
			http.Redirect(w, r, "/profile", http.StatusFound)
		}
	}
}

func DeleteAccount(w http.ResponseWriter, r *http.Request) {

	if utils.URLChecker(w, r, "/delete/account") {

		if r.Method == "POST" {
			var p models.User
			err := json.NewDecoder(r.Body).Decode(&p.ID)
			if err != nil {
				log.Println(err)
			}
			p.DeleteAccount(w, r)
			fmt.Println("deleted account by ID", p.ID)
			http.Redirect(w, r, "/", 302)
		}
	}
}
