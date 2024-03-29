package controllers

import (
	"ForumX/utils"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(1, 3)

func limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			//http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			utils.RenderTemplate(w, "404page", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// high order function func(func)(callback)
//case 1: signin -> set session, & cookie Browser, -> redirect Middleware(Profile)
//each handler - isCookie() - check Browser cookie value - and Db, if ok -> save session - global variable
func IsValidCookie(f http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		//check expires cookie
		c, err := r.Cookie("_cookie")
		if err != nil {
			log.Println(err, "expires timeout || cookie deleted")
			utils.Logout(w, r, *session)
			return
		}
		//cookie Browser -> send IsCookie(check if this user ->)
		// then call handler -> middleware
		if isValidCookie, sessionF := utils.IsCookie(w, r, c.Value); isValidCookie {
			err = DB.QueryRow("SELECT cookie_time FROM session WHERE user_id = $1", sessionF.UserID).Scan(&sessionF.Time)
			if err != nil {
				log.Println(err, "err session time")
			}
			strToTime, _ := time.Parse(time.RFC3339, sessionF.Time)
			diff := time.Now().Sub(strToTime)

			if int(diff.Minutes()) > 290 && int(diff.Seconds()) < 298 {
				uuid := utils.CreateUuid()
				utils.SetCookie(w, uuid)
				utils.ReSession(sessionF.UserID, session, "timeout", uuid)
				fmt.Println("change cookie Browser and update sessiontime and uuid in Db")
			}
			*session = sessionF
			f(w, r)
		}
	}
}

//Init func handlers
func Init() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "6969" // Default port if not specified
	}
	//create multiplexer
	mux := http.NewServeMux()
	//file server
	mux.Handle("/statics/", http.StripPrefix("/statics/", http.FileServer(http.Dir("./statics/"))))

	mux.HandleFunc("/", GetAllPosts)
	mux.HandleFunc("/sapid", GetAllPosts)
	mux.HandleFunc("/love", GetAllPosts)
	mux.HandleFunc("/science", GetAllPosts)

	mux.HandleFunc("/post", GetPostByID)
	mux.HandleFunc("/create/post", IsValidCookie(CreatePost))
	mux.HandleFunc("/edit/post", IsValidCookie(UpdatePost))
	mux.HandleFunc("/delete/post", IsValidCookie(DeletePost))

	mux.HandleFunc("/comment", IsValidCookie(LeaveComment))
	mux.HandleFunc("/edit/comment", IsValidCookie(UpdateComment))
	mux.HandleFunc("/delete/comment", IsValidCookie(DeleteComment))
	mux.HandleFunc("/reply/comment", IsValidCookie(ReplyComment))

	mux.HandleFunc("/votes/post", IsValidCookie(VotesPost))
	mux.HandleFunc("/votes/comment", IsValidCookie(VotesComment))

	mux.HandleFunc("/signin", Signin)
	mux.HandleFunc("/signup", Signup)
	mux.HandleFunc("/googleSignin", GoogleSignin)
	mux.HandleFunc("/googleUserInfo", GoogleUserData)

	mux.HandleFunc("/githubSignin", GithubSignin)
	mux.HandleFunc("/githubUserInfo", GithubUserData)
	mux.HandleFunc("/logout", IsValidCookie(Logout))

	mux.HandleFunc("/profile", IsValidCookie(GetUserProfile))
	mux.HandleFunc("/user/id", GetAnotherProfile)
	mux.HandleFunc("/edit/user", IsValidCookie(UpdateProfile))
	mux.HandleFunc("/delete/account", IsValidCookie(DeleteAccount))

	mux.HandleFunc("/activity", IsValidCookie(GetUserActivities))
	mux.HandleFunc("/search", Search)
	// http.HandleFunc("/chat", routing.StartChat)
	log.Println("Listening port:", port)
	log.Fatal(http.ListenAndServe(":"+port, limit(mux)))
}
