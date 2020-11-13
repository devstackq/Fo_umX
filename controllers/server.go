package controllers

import (
	"log"
	"net/http"
)

// type Server struct {
// 	config *config.Config
// }

// func New(conf *config.Config) *Server {
// 	return &Server{
// 		config: conf,
// 	}
// }

// func (s *Server) Run() error {
// 	return nil
// }

//handlers
//mux own server,  route init  - google, config FileServer
//handlers
func Init() {
	http.Handle("/statics/", http.StripPrefix("/statics/", http.FileServer(http.Dir("./statics/"))))

	http.HandleFunc("/", GetAllPosts)
	http.HandleFunc("/sapid", GetAllPosts)
	http.HandleFunc("/love", GetAllPosts)
	http.HandleFunc("/science", GetAllPosts)

	http.HandleFunc("/post", GetPostByID)
	http.HandleFunc("/create/post", CreatePost)
	http.HandleFunc("/edit/post", UpdatePost)
	http.HandleFunc("/delete/post", DeletePost)

	http.HandleFunc("/comment", LeaveComment)

	http.HandleFunc("/votes", VotesPost)
	http.HandleFunc("/votes/comment", VotesComment)
	http.HandleFunc("/search", Search)

	http.HandleFunc("/profile", GetUserProfile)
	http.HandleFunc("/user/id/", GetAnotherProfile)
	http.HandleFunc("/edit/user", UpdateProfile)
	http.HandleFunc("/delete/account", DeleteAccount)

	http.HandleFunc("/signup", Signup)
	http.HandleFunc("/signin", Signin)
	http.HandleFunc("/logout", Logout)
	// http.HandleFunc("/chat", routing.StartChat)
	log.Fatal(http.ListenAndServe(":6969", nil))
}
