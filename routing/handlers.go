package routing

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"

	"net/http"
	"strconv"

	"github.com/devstackq/ForumX/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	err  error
	DB   *sql.DB
	temp = template.Must(template.ParseFiles("templates/header.html", "templates/category_temp.html", "templates/likedpost.html", "templates/likes.html", "templates/404page.html", "templates/postupdate.html", "templates/postuser.html", "templates/commentuser.html", "templates/userupdate.html", "templates/search.html", "templates/user.html", "templates/commentuser.html", "templates/postuser.html", "templates/profile.html", "templates/signin.html", "templates/user.html", "templates/signup.html", "templates/filter.html", "templates/posts.html", "templates/comment.html", "templates/create.html", "templates/footer.html", "templates/index.html"))
)

//cahce html file
func DisplayTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := temp.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		http.Error(w, err.Error(),
			http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
	}
}

//getAllPost and Posts by cateogry
//receive request, from client, query params, category ID, then query DB, depends catID, get Post this catID
func GetAllPosts(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" && r.URL.Path != "/science" && r.URL.Path != "/love" && r.URL.Path != "/sapid" {
		DisplayTemplate(w, "404page", http.StatusNotFound)
		return
	}

	auth := models.API
	for _, cookie := range r.Cookies() {
		if cookie.Name == "_cookie" {
			auth.Authenticated = true
			break
		}
	}
	posts, endpoint, err := models.GetAllPost(r)
	if err != nil {
		log.Fatal(err)
	}

	DisplayTemplate(w, "header", auth)

	// endpoint -> get post by category
	// profile/ fix, create, get post fix
	if endpoint == "/" {
		DisplayTemplate(w, "index", posts)
	} else {
		DisplayTemplate(w, "catTemp", posts)
	}
}

//view 1 post by id
func GetPostById(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/post" {
		DisplayTemplate(w, "404page", http.StatusNotFound)
		return
	}
	//check cookie for  navbar, if not cookie - signin
	auth := models.API
	for _, cookie := range r.Cookies() {
		if cookie.Name == "_cookie" {
			auth.Authenticated = true
			break
		}
	}

	comments, post, err := models.GetPostById(r)
	if err != nil {
		panic(err)
	}
	DisplayTemplate(w, "header", auth)
	DisplayTemplate(w, "posts", post)
	DisplayTemplate(w, "comment", comments)
}

//create post
func CreatePost(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/create/post" {
		DisplayTemplate(w, "404page", http.StatusNotFound)
		return
	}

	auth := models.API
	for _, cookie := range r.Cookies() {
		if cookie.Name == "_cookie" {
			auth.Authenticated = true
			break
		}
	}

	msg := models.API
	msg.Msg = ""

	switch r.Method {
	case "GET":
		DisplayTemplate(w, "header", auth)
		DisplayTemplate(w, "create", &msg)
	case "POST":
		access := CheckCookies(w, r)
		if !access {
			http.Redirect(w, r, "/signin", 302)
			return
		}
		r.ParseForm()

		c, _ := r.Cookie("_cookie")
		s := models.Session{UUID: c.Value}

		r.ParseMultipartForm(10 << 20)
		file, _, err := r.FormFile("uploadfile")

		// fImg, err := os.Open("./1553259670.jpg")

		// if err != nil {
		// 	fmt.Println(err)
		// 	os.Exit(1)
		// }
		// defer fImg.Close()

		// imgInfo, err := fImg.Stat()
		// if err != nil {
		// 	fmt.Println(err, "stats")
		// 	os.Exit(1)
		// }

		// var size int64 = imgInfo.Size()
		// fmt.Println(size, "size")
		// byteArr := make([]byte, size)

		// read file into bytes
		// buffer := bufio.NewReader(fImg)
		// _, err = buffer.Read(byteArr)
		//defer fImg.Close()
		var fileBytes []byte
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			//file = fImg
			//fileBytes = byteArr
		}

		var buff bytes.Buffer
		fileSize, _ := buff.ReadFrom(file)
		defer file.Close()

		// fmt.Println(fileSize)
		// var max int64
		// max = 20000000

		if fileSize < 20000000 {
			file2, _, err := r.FormFile("uploadfile")
			if err != nil {
				log.Fatal(err)
				//file2 = fImg
				//fileBytes = byteArr
			}
			defer file2.Close()
			fileBytes, _ = ioutil.ReadAll(file2)
		} else {
			fmt.Print("file more 20mb")
			//messga clinet send
			msg.Msg = "Large file, more than 20mb"
			DisplayTemplate(w, "header", auth)
			DisplayTemplate(w, "create", &msg)
			return
		}

		DB.QueryRow("SELECT user_id FROM session WHERE uuid = ?", s.UUID).Scan(&s.UserID)

		title := r.FormValue("title")
		content := r.FormValue("content")
		//check empty values
		checkInputTitle := false
		checkInputContent := false

		for _, v := range title {
			if v >= 97 && v <= 122 || v >= 65 && v <= 90 && v >= 32 && v <= 64 || v > 128 {
				checkInputTitle = true
			}
		}
		for _, v := range content {
			if v >= 97 && v <= 122 || v >= 65 && v <= 90 && v >= 32 && v <= 64 || v > 128 {
				checkInputContent = true
			}
		}

		if checkInputTitle && checkInputContent {

			p := models.Posts{
				Title:     title,
				Content:   content,
				CreatorID: s.UserID,
				Image:     fileBytes,
			}

			lastPost, err := p.CreatePost()
			//fmt.Println(lastPost)
			//query last post -> db
			if err != nil {
				fmt.Println(err)
				//	panic(err.Error())
			}

			//insert cat_post_bridge value
			categories, _ := r.Form["input"]
			if len(categories) == 1 {
				pcb := models.PostCategory{
					PostID:   lastPost,
					Category: categories[0],
				}
				err = pcb.CreateBridge()
			} else if len(categories) > 1 {
				//loop
				for _, v := range categories {
					pcb := models.PostCategory{
						PostID:   lastPost,
						Category: v,
					}
					err = pcb.CreateBridge()
				}
			}

			http.Redirect(w, r, "/", http.StatusFound)
			w.WriteHeader(http.StatusCreated)
		} else {
			msg.Msg = "Empty title or content"
			DisplayTemplate(w, "header", auth)
			DisplayTemplate(w, "create", &msg)
		}
	}
}

//update post
func UpdatePost(w http.ResponseWriter, r *http.Request) {
	var pid int
	if r.Method == "GET" {

		pid, _ = strconv.Atoi(r.URL.Query().Get("id"))
		p := models.Posts{}
		p.PostIDEdit = pid
		DisplayTemplate(w, "updatepost", p)

	}
	if r.Method == "POST" {
		access := CheckCookies(w, r)
		if !access {
			http.Redirect(w, r, "/signin", 302)
			return
		}

		r.ParseForm()
		r.ParseMultipartForm(10 << 20)
		file, _, err := r.FormFile("uploadfile")

		if err != nil {
			panic(err)

		}
		defer file.Close()

		fileBytes, err := ioutil.ReadAll(file)

		if err != nil {
			panic(err)
		}

		pid := r.FormValue("pid")
		pidnum, _ := strconv.Atoi(pid)

		p := models.Posts{
			Title:      r.FormValue("title"),
			Content:    r.FormValue("content"),
			Image:      fileBytes,
			PostIDEdit: pidnum,
		}

		err = p.UpdatePost()

		if err != nil {
			panic(err.Error())
		}
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

//delete post
func DeletePost(w http.ResponseWriter, r *http.Request) {

	var pid int

	pid, _ = strconv.Atoi(r.URL.Query().Get("id"))
	p := models.Posts{}
	p.PostIDEdit = pid

	access := CheckCookies(w, r)
	if !access {
		http.Redirect(w, r, "/signin", 302)
		return
	}

	err = p.DeletePost()

	if err != nil {
		panic(err.Error())
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

//create comment
func CreateComment(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/comment" {
		DisplayTemplate(w, "404page", http.StatusNotFound)
		return
	}

	if r.Method == "POST" {

		access := CheckCookies(w, r)
		if !access {
			http.Redirect(w, r, "/signin", 302)
			return
		}

		r.ParseForm()
		c, _ := r.Cookie("_cookie")
		s := models.Session{UUID: c.Value}
		DB.QueryRow("SELECT user_id FROM session WHERE uuid = ?", s.UUID).Scan(&s.UserID)

		pid, _ := strconv.Atoi(r.FormValue("curr"))
		comment := r.FormValue("comment-text")

		checkLetter := false
		for _, v := range comment {
			if v >= 97 && v <= 122 || v >= 65 && v <= 90 && v >= 32 && v <= 64 || v > 128 {
				checkLetter = true
			}
		}

		if checkLetter {
			com := models.Comments{
				Commentik: comment,
				PostID:    pid,
				UserID:    s.UserID,
			}

			err = com.LostComment()

			if err != nil {
				panic(err.Error())

			}
		}
		http.Redirect(w, r, "post?id="+r.FormValue("curr"), 301)
	}
}

//profile current -> user page
func GetProfileById(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/profile" {
		DisplayTemplate(w, "404page", http.StatusNotFound)
		return
	}

	if r.Method == "GET" {

		auth := models.API
		for _, cookie := range r.Cookies() {
			if cookie.Name == "_cookie" {
				auth.Authenticated = true
				break
			}
		}
		//if userId now, createdPost uid equal -> show
		likedpost, posts, comments, user, err := models.GetUserProfile(r)
		if err != nil {
			panic(err)
		}

		DisplayTemplate(w, "header", auth)
		DisplayTemplate(w, "profile", user)
		DisplayTemplate(w, "likedpost", likedpost)
		DisplayTemplate(w, "postuser", posts)
		DisplayTemplate(w, "commentuser", comments)
	}
}

//user page, other anyone
func GetUserById(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		auth := models.API
		for _, cookie := range r.Cookies() {
			if cookie.Name == "_cookie" {
				auth.Authenticated = true
				break
			}
		}

		posts, user, err := models.GetOtherUser(r)
		if err != nil {
			panic(err)
		}

		DisplayTemplate(w, "header", auth)
		DisplayTemplate(w, "user", user)
		DisplayTemplate(w, "postuser", posts)
	}
}

//update profile
func UpdateProfile(w http.ResponseWriter, r *http.Request) {

	//check cookie for  navbar
	auth := models.API
	for _, cookie := range r.Cookies() {
		if cookie.Name == "_cookie" {
			auth.Authenticated = true
			break
		}
	}

	if r.Method == "GET" {
		DisplayTemplate(w, "header", auth)
		DisplayTemplate(w, "updateuser", "")
	}

	if r.Method == "POST" {
		access := CheckCookies(w, r)
		if !access {
			http.Redirect(w, r, "/signin", 302)
			return
		}

		r.ParseForm()
		r.ParseMultipartForm(10 << 20)
		file, _, err := r.FormFile("uploadfile")

		if err != nil {
			panic(err)
		}

		defer file.Close()

		fileBytes, err := ioutil.ReadAll(file)

		if err != nil {
			panic(err)
		}

		c, _ := r.Cookie("_cookie")
		s := models.Session{UUID: c.Value}

		DB.QueryRow("SELECT user_id FROM session WHERE uuid = ?", s.UUID).
			Scan(&s.UserID)

		is, _ := strconv.Atoi(r.FormValue("age"))

		p := models.Users{
			FullName: r.FormValue("fullname"),
			Age:      is,
			Sex:      r.FormValue("sex"),
			City:     r.FormValue("city"),
			Image:    fileBytes,
			ID:       s.UserID,
		}

		err = p.UpdateProfile()

		if err != nil {
			panic(err.Error())
		}
	}

	http.Redirect(w, r, "/profile", http.StatusFound)
}

//search
func Search(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/search" {
		DisplayTemplate(w, "404page", http.StatusNotFound)
		return
	}

	if r.Method == "GET" {
		DisplayTemplate(w, "search", http.StatusFound)
	}

	if r.Method == "POST" {

		auth := models.API
		for _, cookie := range r.Cookies() {
			if cookie.Name == "_cookie" {
				auth.Authenticated = true
				break
			}
		}

		findPosts, err := models.Search(w, r)

		if err != nil {
			panic(err)
		}

		DisplayTemplate(w, "header", auth)
		DisplayTemplate(w, "index", findPosts)
	}
}

//signup system
func Signup(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/signup" {
		DisplayTemplate(w, "404page", http.StatusNotFound)
		return
	}
	msg := models.API

	if r.Method == "GET" {
		DisplayTemplate(w, "signup", &msg)
	}

	if r.Method == "POST" {

		fn := r.FormValue("fullname")
		e := r.FormValue("email")
		p := r.FormValue("password")
		a := r.FormValue("age")
		s := r.FormValue("sex")
		c := r.FormValue("city")

		r.ParseMultipartForm(10 << 20)
		file, _, err := r.FormFile("uploadfile")

		if err != nil {
			panic(err)
		}

		defer file.Close()

		fB, err := ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(p), 8)
		if err != nil {
			panic(err)
		}

		//check email by unique, if have same email
		checkEmail, err := DB.Query("SELECT email FROM users")
		if err != nil {
			panic(err)
		}

		all := []models.Users{}

		for checkEmail.Next() {
			user := models.Users{}
			var email string
			err = checkEmail.Scan(&email)
			if err != nil {
				panic(err.Error)
			}

			user.Email = email
			all = append(all, user)
		}

		for _, v := range all {
			if v.Email == e {
				msg.Msg = "Not unique email lel"
				DisplayTemplate(w, "signup", &msg)
				return
			}
		}

		_, err = DB.Exec("INSERT INTO users( full_name, email, password, age, sex, city, image) VALUES (?, ?, ?, ?, ?, ?, ?)",
			fn, e, hash, a, s, c, fB)

		if err != nil {
			panic(err.Error())
		}

		http.Redirect(w, r, "/signin", 301)
	}
}

//signin system
func Signin(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/signin" {
		DisplayTemplate(w, "404page", http.StatusNotFound)
		return
	}
	r.Header.Add("Accept", "text/html")
	r.Header.Add("User-Agent", "MSIE/15.0")
	msg := models.API
	msg.Msg = ""

	if r.Method == "GET" {
		DisplayTemplate(w, "signin", &msg)
	}

	if r.Method == "POST" {
		var person models.Users
		//b, _ := ioutil.ReadAll(r.Body)
		err := json.NewDecoder(r.Body).Decode(&person)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println(person)

		if person.Type == "default" {
			fmt.Println(" default auth")

			email := person.Email
			pwd := person.Password

			models.Signin(w, r, email, pwd)

			//http.Redirect(w, r, "/profile", 200)
			//citiesArtist := FindCityArtist(w, r, strings.ToLower(string(body)))
			//w.Header().Set("Content-Type", "application/json")
			//json.NewEncoder(w).Encode(citiesArtist)

			// json.NewEncoder(w).Encode(msg)
			// ok := "okay"
			// b := []byte(ok)

			// msg.Msg = "okay"
			// w.Header().Set("Content-Type", "application/json")
			// m, _ := json.Marshal(msg)
			// w.Write(m)
		} else if person.Type == "google" {
			fmt.Println("todo google auth")
			http.Redirect(w, r, "/profile", http.StatusFound)
		} else if person.Type == "github" {
			fmt.Println("todo github")
			http.Redirect(w, r, "/profile", http.StatusFound)
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		//http.Redirect(w, r, "/profile", 200)
	}
}

// Logout
func Logout(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/logout" {
		DisplayTemplate(w, "404page", http.StatusNotFound)
		return
	}

	if r.Method == "GET" {
		cookie, _ := r.Cookie("_cookie")
		//add cookie -> fields uuid
		s := models.Session{UUID: cookie.Value}
		//get ssesion id, by local struct uuid
		DB.QueryRow("SELECT id FROM session WHERE uuid = ?", s.UUID).
			Scan(&s.ID)
		//delete session by id session
		_, err = DB.Exec("DELETE FROM session WHERE id = ?", s.ID)

		if err != nil {
			panic(err)
		}

		// then delete cookie from client
		cookieDelete := http.Cookie{
			Name:     "_cookie",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookieDelete)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

//check cookie client side and Db
func CheckCookies(w http.ResponseWriter, r *http.Request) bool {

	flag := false
	cookieHave := false
	for _, cookie := range r.Cookies() {
		if cookie.Name == "_cookie" {
			cookieHave = true
			break
		}
	}
	if !cookieHave {
		http.Redirect(w, r, "/signin", 302)
	} else {
		//get client cookie
		cookie, _ := r.Cookie("_cookie")
		//set local struct -> cookie value
		s := models.Session{UUID: cookie.Value}
		var tmp string
		// get userid by Client sessionId
		DB.QueryRow("SELECT user_id FROM session WHERE uuid = ?", s.UUID).
			Scan(&s.UserID)
		//get uuid by userid, and write UUID data
		DB.QueryRow("SELECT uuid FROM session WHERE user_id = ?", s.UserID).Scan(&tmp)
		//check local and DB session
		if cookie.Value == tmp {
			flag = true
		}
	}
	return flag

}

//like dislike post
func LostVotes(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/votes" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	access := CheckCookies(w, r)
	if !access {
		http.Redirect(w, r, "/signin", 302)
		return
	}
	c, _ := r.Cookie("_cookie")
	s := models.Session{UUID: c.Value}

	DB.QueryRow("SELECT user_id FROM session WHERE uuid = ?", s.UUID).
		Scan(&s.UserID)

	pid := r.URL.Query().Get("id")
	lukas := r.FormValue("lukas")
	diskus := r.FormValue("diskus")

	if r.Method == "POST" {

		if lukas == "1" {
			//check if not have post and user lost vote this post
			//1 like or 1 dislike 1 user lost 1 post, get previus value and +1
			var p, u int
			err = DB.QueryRow("SELECT post_id, user_id FROM likes WHERE post_id=? AND user_id=?", pid, s.UserID).Scan(&p, &u)

			if p == 0 && u == 0 {

				oldlike := 0
				err = DB.QueryRow("SELECT count_like FROM posts WHERE id=?", pid).Scan(&oldlike)
				nv := oldlike + 1
				_, err = DB.Exec("UPDATE  posts SET count_like = ? WHERE id= ?", nv, pid)
				if err != nil {
					panic(err)
				}

				_, err = DB.Exec("INSERT INTO likes(post_id, user_id) VALUES( ?, ?)", pid, s.UserID)
				if err != nil {
					panic(err)
				}
			}
		}

		if diskus == "1" {

			var p, u int
			err = DB.QueryRow("SELECT post_id, user_id FROM likes WHERE post_id=? AND user_id=?", pid, s.UserID).Scan(&p, &u)

			if p == 0 && u == 0 {

				oldlike := 0
				err = DB.QueryRow("select count_dislike from posts where id=?", pid).Scan(&oldlike)
				nv := oldlike + 1
				_, err = DB.Exec("UPDATE  posts SET count_dislike = ? WHERE id= ?", nv, pid)
				if err != nil {
					panic(err)
				}
				_, err = DB.Exec("INSERT INTO likes(post_id, user_id) VALUES( ?, ?)", pid, s.UserID)

				if err != nil {
					panic(err)
				}
			}
		}
	}
	http.Redirect(w, r, "post?id="+pid, 301)
}

func LostVotesComment(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/votes/comment" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	access := CheckCookies(w, r)
	if !access {
		http.Redirect(w, r, "/signin", 302)
		return
	}
	c, _ := r.Cookie("_cookie")
	s := models.Session{UUID: c.Value}
	DB.QueryRow("SELECT user_id FROM session WHERE uuid = ?", s.UUID).
		Scan(&s.UserID)

	cid := r.URL.Query().Get("cid")
	comdis := r.FormValue("comdis")
	comlike := r.FormValue("comlike")

	pidc := r.FormValue("pidc")

	if r.Method == "POST" {

		if comlike == "1" {

			var c, u int
			err = DB.QueryRow("SELECT comment_id, user_id FROM likes WHERE comment_id=? AND user_id=?", cid, s.UserID).Scan(&c, &u)

			if c == 0 && u == 0 {

				oldlike := 0
				err = DB.QueryRow("SELECT com_like FROM comments WHERE id=?", cid).Scan(&oldlike)
				nv := oldlike + 1

				_, err = DB.Exec("UPDATE  comments SET com_like = ? WHERE id= ?", nv, cid)

				if err != nil {
					panic(err)
				}

				_, err = DB.Exec("INSERT INTO likes(comment_id, user_id) VALUES( ?, ?)", cid, s.UserID)
				if err != nil {
					panic(err)
				}
			}
		}

		if comdis == "1" {

			var c, u int
			err = DB.QueryRow("SELECT comment_id, user_id FROM likes WHERE comment_id=? AND user_id=?", cid, s.UserID).Scan(&c, &u)

			if c == 0 && u == 0 {

				oldlike := 0
				err = DB.QueryRow("SELECT com_dislike FROM comments WHERE id=?", cid).Scan(&oldlike)
				nv := oldlike + 1

				_, err = DB.Exec("UPDATE  comments SET com_dislike = ? WHERE id= ?", nv, cid)

				if err != nil {
					panic(err)
				}

				_, err = DB.Exec("INSERT INTO likes(comment_id, user_id) VALUES( ?, ?)", cid, s.UserID)
				if err != nil {
					panic(err)
				}
			}
		}
		http.Redirect(w, r, "/post?id="+pidc, 301)
	}
}

//Likes table, filed posrid, userid, state_id
// 0,1,2 if state ==0, 1 || 2,
// next btn, if 1 == 1, state =0
