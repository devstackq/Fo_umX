package models

import (
	"fmt"
	"log"
	"net/http"
	"time"

	structure "github.com/devstackq/ForumX/general"
	util "github.com/devstackq/ForumX/utils"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func (u *User) Signup(w http.ResponseWriter, r *http.Request) {

	users := []User{}

	hashPwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), 8)
	if err != nil {
		log.Println(err)
	}
	//check email by unique, if have same email
	checkEmail, err := DB.Query("SELECT email FROM users")
	if err != nil {
		log.Println(err)
	}

	for checkEmail.Next() {
		user := User{}
		var email string
		err = checkEmail.Scan(&email)
		if err != nil {
			log.Println(err.Error)
		}

		user.Email = email
		users = append(users, user)
	}

	for _, v := range users {
		if v.Email == u.Email {
			msg = "Not unique email lel"
			util.DisplayTemplate(w, "signup", &msg)
			log.Println(err)
		}
	}

	_, err = DB.Exec("INSERT INTO users( full_name, email, password, age, sex, created_time, city, image) VALUES (?,?, ?, ?, ?, ?, ?, ?)",
		u.FullName, u.Email, hashPwd, u.Age, u.Sex, time.Now(), u.City, u.Image)

	if err != nil {
		log.Println(err)
	}

}

//Signin function
func (uStr *User) Signin(w http.ResponseWriter, r *http.Request) {

	u := DB.QueryRow("SELECT id, password FROM users WHERE email=?", uStr.Email)

	var user User
	var err error
	//check pwd, if not correct, error
	err = u.Scan(&user.ID, &user.Password)
	if err != nil {
		util.AuthError(w, err, "user not found")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(uStr.Password))
	if err != nil {
		util.AuthError(w, err, "password incorrect")
		return
	}
	//get user by Id, and write session struct
	s := structure.Session{
		UserID: user.ID,
	}
	uuid := uuid.Must(uuid.NewV4(), err).String()
	if err != nil {
		util.AuthError(w, err, "uuid trouble")
		return
	}
	//create uuid and set uid DB table session by userid,
	_, err = DB.Exec("INSERT INTO session(uuid, user_id) VALUES (?, ?)", uuid, s.UserID)

	if err != nil {
		util.AuthError(w, err, "the user is already in the system")
		//get ssesion id, by local struct uuid
		return
	}
	// get user in info by session Id
	err = DB.QueryRow("SELECT id, uuid FROM session WHERE user_id = ?", s.UserID).Scan(&s.ID, &s.UUID)
	if err != nil {
		util.AuthError(w, err, "not find user from session")
		return
	}

	//set cookie 9128ueq9widjaisdh238yrhdeiuwandijsan
	cookie := http.Cookie{
		Name:     "_cookie",
		Value:    s.UUID,
		Path:     "/",
		Expires:  time.Now().Add(300 * time.Minute),
		HttpOnly: false,
	}
	http.SetCookie(w, &cookie)
	util.AuthError(w, nil, "success")
}

//Logout function
func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("_cookie")
	if err != nil {
		fmt.Println(err, "cookie err")
	}
	//add cookie -> fields uuid
	s := structure.Session{UUID: cookie.Value}
	//get ssesion id, by local struct uuid
	DB.QueryRow("SELECT id FROM session WHERE uuid = ?", s.UUID).
		Scan(&s.ID)
	fmt.Println(s.ID, "user id deleted session")
	//delete session by id session
	_, err = DB.Exec("DELETE FROM session WHERE id = ?", s.ID)

	if err != nil {
		log.Println(err)
	}
	// then delete cookie from client
	util.DeleteCookie(w)
}
