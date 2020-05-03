package routing

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"

	"net/http"
	"strconv"

	"github.com/devstackq/Forum-X/models"
	uuid "github.com/satori/go.uuid"

	"golang.org/x/crypto/bcrypt"
)

var (
	rows *sql.Rows
	err  error
	DB   *sql.DB
	temp = template.Must(template.ParseFiles("templates/header.html", "templates/likedpost.html", "templates/likes.html", "templates/404page.html", "templates/postupdate.html", "templates/postuser.html", "templates/commentuser.html", "templates/userupdate.html", "templates/search.html", "templates/user.html", "templates/commentuser.html", "templates/postuser.html", "templates/profile.html", "templates/signin.html", "templates/user.html", "templates/signup.html", "templates/filter.html", "templates/posts.html", "templates/comment.html", "templates/create.html", "templates/footer.html", "templates/index.html"))
)

//cahce html file
func displayTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
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

	if r.URL.Path != "/" && r.URL.Path != "/school" && r.URL.Path != "/people" && r.URL.Path != "/events" && r.URL.Path != "/qa" && r.URL.Path != "/sapid" {
		displayTemplate(w, "404page", http.StatusNotFound)
		return
	}

	auth := models.API
	for _, cookie := range r.Cookies() {
		if cookie.Name == "_cookie" {
			auth.Authenticated = true
			break
		}
	}

	var post models.Posts
	r.ParseForm()
	like := r.FormValue("likes")
	date := r.FormValue("date")
	category, _ := strconv.Atoi(r.FormValue("cats"))

	switch r.URL.Path {
	//check what come client, cats, and filter by like, date and cats
	case "/":
		post.EndpointPost = "/"
		if date == "asc" {
			rows, err = DB.Query("SELECT  * FROM posts WHERE category_id =?  ORDER BY  created_time ASC limit 8", category)
		} else if date == "desc" {
			rows, err = DB.Query("SELECT  * FROM posts WHERE category_id  =?  ORDER BY  created_time  DESC limit 8", category)
		} else if like == "big" {
			rows, err = DB.Query("SELECT  * FROM posts  WHERE   category_id  =? ORDER BY  count_like DESC limit 8", category)
		} else if like == "letter" {
			rows, err = DB.Query("SELECT  * FROM posts  WHERE category_id  =?  ORDER BY  count_dislike DESC limit 8", category)
		} else if category > 0 {
			rows, err = DB.Query("SELECT  * FROM posts  WHERE category_id  =?   ORDER BY  created_time  DESC limit 8", category)
		} else {
			rows, err = DB.Query("SELECT  * FROM posts  ORDER BY  created_time  DESC limit 8")
		}

	case "/school":
		post.EndpointPost = "/school"
		rows, err = DB.Query("SELECT  * FROM posts  WHERE category_id =?   ORDER  BY created_time  DESC   LIMIT 4", 1)
	case "/people":
		post.EndpointPost = "/people"
		rows, err = DB.Query("SELECT  * FROM posts  WHERE category_id =?   ORDER  BY created_time  DESC LIMIT 4", 2)
	case "/events":
		post.EndpointPost = "/events"
		rows, err = DB.Query("SELECT  * FROM posts  WHERE category_id =?   ORDER  BY created_time  DESC LIMIT 4", 3)
	case "/qa":
		post.EndpointPost = "/qa"
		rows, err = DB.Query("SELECT  * FROM posts  WHERE category_id =?   ORDER  BY created_time  DESC LIMIT 4", 4)
	case "/sapid":
		post.EndpointPost = "/sapid"
		rows, err = DB.Query("SELECT  * FROM posts  WHERE category_id =?   ORDER  BY created_time  DESC LIMIT 4", 5)
	}

	defer rows.Close()
	if err != nil {
		panic(err)
	}

	var PostsAll []models.Posts

	for rows.Next() {
		postik := models.Posts{}
		if err := rows.Scan(&postik.ID, &postik.Title, &postik.Content, &postik.CreatorID, &postik.CategoryID, &postik.CreationTime, &postik.Image, &postik.CountLike, &postik.CountDislike); err != nil {
			panic(err)
		}
		DB.QueryRow("SELECT title FROM categories WHERE id=?", postik.CategoryID).Scan(&postik.CategoryName)
		PostsAll = append(PostsAll, postik)
	}

	displayTemplate(w, "header", auth)
	displayTemplate(w, "index", PostsAll)
}

//view 1 post by id
func GetPostById(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/post" {
		displayTemplate(w, "404page", http.StatusNotFound)
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
		return
	}
	displayTemplate(w, "header", auth)
	displayTemplate(w, "posts", post)
	displayTemplate(w, "comment", comments)
}

//create post
func CreatePost(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/create/post" {
		displayTemplate(w, "404page", http.StatusNotFound)
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
		displayTemplate(w, "header", auth)
		displayTemplate(w, "create", &msg)
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
		if err != nil {
			panic(err)

		}

		var buff bytes.Buffer
		fileSize, _ := buff.ReadFrom(file)
		defer file.Close()

		// fmt.Println(fileSize)
		var max int64
		max = 20000000

		var fileBytes []byte

		if fileSize < max {
			file2, _, _ := r.FormFile("uploadfile")
			fileBytes, _ = ioutil.ReadAll(file2)
		} else {
			fmt.Print("file more 20mb")
			//messga clinet send
			msg.Msg = "Large file, more than 20mb"
			displayTemplate(w, "header", auth)
			displayTemplate(w, "create", &msg)
			return
		}
		// fmt.Println(fileBytes)

		if err != nil {
			panic(err)
		}

		DB.QueryRow("SELECT user_id FROM session WHERE uuid = ?", s.UUID).Scan(&s.UserID)

		category, _ := strconv.Atoi(r.FormValue("cats"))

		title := r.FormValue("title")
		content := r.FormValue("content")
		//check empty values
		norm := false
		nutaksebe := false

		for _, v := range title {
			if v >= 97 && v <= 122 || v >= 65 && v <= 90 && v >= 32 && v <= 64 || v > 128 {
				norm = true
			}
		}
		for _, v := range content {
			if v >= 97 && v <= 122 || v >= 65 && v <= 90 && v >= 32 && v <= 64 || v > 128 {
				nutaksebe = true
			}
		}

		if norm && nutaksebe {

			p := models.Posts{
				Title:      title,
				Content:    content,
				CreatorID:  s.UserID,
				CategoryID: category,
				Image:      fileBytes,
			}

			err = p.CreatePost()

			if err != nil {
				panic(err.Error())

			}

			http.Redirect(w, r, "/", http.StatusFound)
			w.WriteHeader(http.StatusCreated)
		} else {
			msg.Msg = "Empty title or content"
			displayTemplate(w, "header", auth)
			displayTemplate(w, "create", &msg)
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
		displayTemplate(w, "updatepost", p)

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

		cat, _ := strconv.Atoi(r.FormValue("cats"))
		pid := r.FormValue("pid")
		pidnum, _ := strconv.Atoi(pid)

		p := models.Posts{
			Title:      r.FormValue("title"),
			Content:    r.FormValue("content"),
			CategoryID: cat,
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
		displayTemplate(w, "404page", http.StatusNotFound)
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

		norm := false
		for _, v := range comment {
			if v >= 97 && v <= 122 || v >= 65 && v <= 90 && v >= 32 && v <= 64 || v > 128 {
				norm = true
			}
		}

		if norm {
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

//profile current user page
func GetProfileById(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/profile" {
		displayTemplate(w, "404page", http.StatusNotFound)
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

		displayTemplate(w, "header", auth)
		displayTemplate(w, "profile", user)
		displayTemplate(w, "likedpost", likedpost)
		displayTemplate(w, "postuser", posts)
		displayTemplate(w, "commentuser", comments)
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

		displayTemplate(w, "header", auth)
		displayTemplate(w, "user", user)
		displayTemplate(w, "postuser", posts)
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
		displayTemplate(w, "header", auth)
		displayTemplate(w, "updateuser", "")
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
		displayTemplate(w, "404page", http.StatusNotFound)
		return
	}

	if r.Method == "GET" {
		displayTemplate(w, "search", http.StatusFound)
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

		displayTemplate(w, "header", auth)
		displayTemplate(w, "index", findPosts)
	}
}

//signup system
func Signup(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/signup" {
		displayTemplate(w, "404page", http.StatusNotFound)
		return
	}
	msg := models.API

	if r.Method == "GET" {
		displayTemplate(w, "signup", &msg)
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
				displayTemplate(w, "signup", &msg)
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
		displayTemplate(w, "404page", http.StatusNotFound)
		return
	}

	r.Header.Add("Accept", "text/html")
	r.Header.Add("User-Agent", "MSIE/15.0")
	msg := models.API
	msg.Msg = ""
	if r.Method == "GET" {
		displayTemplate(w, "signin", &msg)
	}

	if r.Method == "POST" {
		r.ParseForm()

		email := r.FormValue("email")
		pwd := r.FormValue("password")

		u := DB.QueryRow("SELECT id, password FROM users WHERE email=?", email)

		var user models.Users
		//check pwd, if not correct, error
		u.Scan(&user.ID, &user.Password)

		if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pwd)); err != nil {
			// If the two passwords don't match, return a 401 status
			w.WriteHeader(http.StatusUnauthorized)
			msg.Msg = "Email or password incorrect lel"
			displayTemplate(w, "signin", &msg)
			return
		}

		//get user by Id, and write session struct
		s := models.Session{
			UserID: user.ID,
		}
		uuid := uuid.Must(uuid.NewV4(), err).String()
		if err != nil {
			panic(err)
		}
		//create uuid and set uid DB table session by userid,
		_, err = DB.Exec("INSERT INTO session(uuid, user_id) VALUES (?, ?)", uuid, s.UserID)
		if err != nil {
			panic(err)
		}
		// get user in info by session Id
		DB.QueryRow("SELECT id, uuid FROM session WHERE user_id = ?", s.UserID).Scan(&s.ID, &s.UUID)
		//set cookie
		//uuidUSoro@mail.com -> 9128ueq9widjaisdh238yrhdeiuwandijsan
		//CLient, DB
		// Crete post -> Cleint cookie == session, Userd
		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    s.UUID,
			Path:     "/",
			MaxAge:   84000,
			HttpOnly: false,
		}
		fmt.Println(cookie.MaxAge)
		if cookie.MaxAge == 0 {
			_, err = DB.Exec("DELETE FROM session WHERE id = ?", s.ID)
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/profile", http.StatusFound)
	}
}

// Logout
func Logout(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/logout" {
		displayTemplate(w, "404page", http.StatusNotFound)
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
