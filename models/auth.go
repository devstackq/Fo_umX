package models

import (
	"ForumX/general"
	"ForumX/utils"
	"fmt"
	"log"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

//Signup func
func (u User) Signup(w http.ResponseWriter, r *http.Request) {

	var hashPwd []byte
	if utils.AuthType == "default" {
		hashPwd, err = bcrypt.GenerateFromPassword([]byte(u.Password), 8)
		if err != nil {
			log.Println(err)
		}
	}

	emailCheck := utils.IsRegistered(w, r, u.Email)
	userCheck := utils.IsRegistered(w, r, u.Username)
	fmt.Println(userCheck, u.Username, "Sihnup check")

	if !emailCheck && !userCheck {
		userPrepare, err := DB.Prepare(`INSERT INTO users(full_name, email, username, password, age, sex, created_time, city, image) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)`)
		if err != nil {
			log.Println(err)
		}

		_, err = userPrepare.Exec(u.FullName, u.Email, u.Username, hashPwd, u.Age, u.Sex, time.Now(), u.City, u.Image)
		if err != nil {
			log.Println(err)
		}
		defer userPrepare.Close()

	} else {
		if utils.AuthType == "default" {

			if emailCheck {
				utils.AuthError(w, r, err, "Not unique email", utils.AuthType)
				return
			}
			if userCheck {
				utils.AuthError(w, r, err, "Not unique username", utils.AuthType)
				return
			}
			if userCheck && emailCheck {
				utils.AuthError(w, r, err, "Not unique username && email", utils.AuthType)
				return
			}
		}
	}
	if utils.AuthType == "default" {
		utils.AuthError(w, r, nil, "success", utils.AuthType)
		http.Redirect(w, r, "/profile", 302)
	}
}

//Signin function dsds
//if user no system -> create uuid, save save in Db & browser, -> redirect middleware(profile)
// if second time -> check by Email || username, if
func (uStr *User) Signin(w http.ResponseWriter, r *http.Request) {

	var user User

	err = DB.QueryRow("SELECT id FROM users WHERE email=$1", uStr.Email).Scan(&user.ID)
	if err != nil {
		log.Println(err)
	}

	if utils.AuthType == "default" {

		if uStr.Email != "" {
			err = DB.QueryRow("SELECT id, password FROM users WHERE email=$1", uStr.Email).Scan(&user.ID, &user.Password)
			if err != nil {
				log.Println("err email")
				utils.AuthError(w, r, err, "user by Email not found", utils.AuthType)
				return
			}
			utils.ReSession(user.ID, uStr.Session, "", "")
		} else if uStr.Username != "" {
			err = DB.QueryRow("SELECT id, password FROM users WHERE username=$1", uStr.Username).Scan(&user.ID, &user.Password)
			if err != nil {
				log.Println("errr username")
				utils.AuthError(w, r, err, "user by Username not found", utils.AuthType)
				return
			}
			utils.ReSession(user.ID, uStr.Session, "", "")
		}
		//check pwd, if not correct, error
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(uStr.Password))
		if err != nil {
			utils.AuthError(w, r, err, "password incorrect", utils.AuthType)
			return
		}
	} else {
		fmt.Println("google github")
		utils.ReSession(user.ID, uStr.Session, "", "")
	}

	//1 time set uuid user, set cookie in Browser
	newSession := general.Session{
		UserID: user.ID,
	}

	uuid := uuid.Must(uuid.NewV4(), err).String()
	if err != nil {
		utils.AuthError(w, r, err, "uuid problem", utils.AuthType)
		return
	}
	//create uuid and set uid DB table session by userid,
	userPrepare, err := DB.Prepare(`INSERT INTO session(uuid, user_id, cookie_time) VALUES ($1,$2,$3)`)
	if err != nil {
		log.Println(err)
	}

	_, err = userPrepare.Exec(uuid, newSession.UserID, time.Now())
	if err != nil {
		log.Println(err)
	}
	defer userPrepare.Close()

	// get user in info by session Id
	err = DB.QueryRow("SELECT id, uuid FROM session WHERE user_id = $1", newSession.UserID).Scan(&newSession.ID, &newSession.UUID)
	if err != nil {
		utils.AuthError(w, r, err, "not find user by session", utils.AuthType)
		log.Println(err)
		return
	}
	utils.SetCookie(w, newSession.UUID)
	utils.AuthError(w, r, nil, "success", utils.AuthType)
	fmt.Println(utils.AuthType, "auth type")
	http.Redirect(w, r, "/profile", 302)
}

//Logout function
func Logout(w http.ResponseWriter, r *http.Request, s general.Session) {

	utils.Logout(w, r, s)
	if utils.AuthType == "google" {
		_, err = http.Get("https://accounts.google.com/o/oauth2/revoke?token=" + utils.Token)
		if err != nil {
			log.Println(err)
		}
	}
}
