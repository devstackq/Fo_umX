package util

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
	"unicode"

	structure "github.com/devstackq/ForumX/general"
)

var (
	DB   *sql.DB
	err  error
	temp = template.Must(template.ParseFiles("view/header.html", "view/category_post.html", "view/favorites.html", "view/404page.html", "view/update_post.html", "view/created_post.html", "view/comment_user.html", "view/profile_update.html", "view/search.html", "view/another_user.html", "view/profile.html", "view/signin.html", "view/signup.html", "view/filter.html", "view/post.html", "view/comment_post.html", "view/create_post.html", "view/footer.html", "view/index.html"))
)

type API struct {
	Authenticated bool
}

//IsAuth check user now authorized system ?
func IsAuth(r *http.Request) API {
	var auth API
	for _, cookie := range r.Cookies() {
		if cookie.Name == "_cookie" {
			auth.Authenticated = true
		}
	}
	return auth
}

//IsCookie check user cookie client and DB session value, if true -> give access
func IsCookie(w http.ResponseWriter, r *http.Request) (bool, structure.Session) {

	var flag, cookieHave bool
	cookie, _ := r.Cookie("_cookie")
	s := structure.Session{}

	if IsAuth(r).Authenticated {
		cookieHave = true
	}
	if !cookieHave {
		http.Redirect(w, r, "/signin", 302)

	} else {
		//get client cookie
		//set local struct -> cookie value
		s := structure.Session{UUID: cookie.Value}
		var tmp string
		// get userid by Client sessionId
		err = DB.QueryRow("SELECT user_id FROM session WHERE uuid = ?", s.UUID).
			Scan(&s.UserID)
		//get uuid by userid, and write UUID data
		if err != nil {
			log.Println(err)
		}
		err = DB.QueryRow("SELECT uuid FROM session WHERE user_id = ?", s.UserID).Scan(&tmp)
		if err != nil {
			log.Println(err)
		}
		//check local and DB session
		if cookie.Value == tmp {
			flag = true
		}
	}
	if flag {
		s.UUID = cookie.Value
		return flag, s
	}

	return flag, s
}

//CheckLetter correct letter
func CheckLetter(value string) bool {

	for _, v := range value {
		if v >= 97 && v <= 122 || v >= 65 && v <= 90 || v >= 32 && v <= 64 || v > 128 {
			return true
		}
	}
	return false
}

//DisplayTemplate function
func DisplayTemplate(w http.ResponseWriter, tmpl string, data interface{}) {

	err = temp.ExecuteTemplate(w, tmpl, data)

	if err != nil {
		fmt.Println(err, "exec ERR")
		http.Error(w, err.Error(),
			http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
		return
	}

}

//IsCookieExpiration if cookie time = 0, delete session and cookie client
func IsCookieExpiration(t time.Time, cookie *http.Cookie, w http.ResponseWriter, r *http.Request) {

	for _, cookie := range r.Cookies() {
		if cookie.Name == "_cookie" {
			s := structure.Session{UUID: cookie.Value}
			//get ssesion id, by local struct uuid
			DB.QueryRow("SELECT id FROM session WHERE uuid = ?", s.UUID).
				Scan(&s.ID)

			_, err = DB.Exec("DELETE FROM session WHERE id = ?", s.ID)

			// then delete cookie from client
			DeleteCookie(w)
			http.Redirect(w, r, "/", 200)
			return
		}
	}
}

//IsUnique find unique liked post
func IsUnique(intSlice []int) []int {
	keys := make(map[int]bool)
	list := []int{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

//FileByte func for convert receive file - to fileByte
func FileByte(r *http.Request, typePhoto string) []byte {
	//check user photo || post photo
	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("uploadfile")

	var defImg *os.File
	if err != nil {
		log.Println(err)
		//set default photo user
		if typePhoto == "user" {
			defImg, _ = os.Open("./utils/default-user.jpg")
		}
		file = defImg
	}
	defer file.Close()

	imgBytes, err := ioutil.ReadAll(file)

	if err != nil {
		log.Println(err)
	}

	return imgBytes
}

//AuthError show auth error
func AuthError(w http.ResponseWriter, err error, text string) {

	fmt.Println(text, "errka auth")
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		m, _ := json.Marshal(text)
		w.Write(m)
		return
	} else {
		w.WriteHeader(200)
		m, _ := json.Marshal(text)
		w.Write(m)
	}
}

//URLChecker function
func URLChecker(w http.ResponseWriter, r *http.Request, url string) bool {

	if r.URL.Path != url {
		DisplayTemplate(w, "404page", http.StatusNotFound)
		return false
	}
	return true
}

//IsEmailValid function
func IsEmailValid(email string) bool {
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return Re.MatchString(email)
}

//IsPasswordValid function
func IsPasswordValid(s string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(s) >= 7 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

//DeleteCookie func
func DeleteCookie(w http.ResponseWriter) {

	cookieDelete := http.Cookie{
		Name:     "_cookie",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: false,
	}
	http.SetCookie(w, &cookieDelete)
}

//IsImage func
func IsImage(r *http.Request) []byte {

	f, _, _ := r.FormFile("uploadfile")
	photoFlag := false

	if f != nil {
		photoFlag = true
	}
	var imgBytes []byte

	if !photoFlag {
		imgBytes = []byte{0, 0}
	} else {
		imgBytes = FileByte(r, "post")
	}
	return imgBytes
}
