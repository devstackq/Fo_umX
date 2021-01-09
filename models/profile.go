package models

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"ForumX/general"
	"ForumX/utils"
)

//User struct
type User struct {
	ID          int       `json:"id"`
	FullName    string    `json:"fullName"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	IsAdmin     bool      `json:"isAdmin"`
	Age         int       `json:"age"`
	Sex         string    `json:"sex"`
	CreatedTime time.Time `json:"createdTime"`
	City        string    `json:"city"`
	Image       []byte    `json:"image"`
	ImageHTML   string    `json:"imageHtml"`
	Role        string    `json:"role"`
	SVG         bool      `json:"svg"`
	Type        string    `json:"type"`
	Temp        string    `json:"temp"`
	Name        string    `json:"name"`
	Location    string    `json:"location"`
	Username    string    `json:"username"`
	Session *general.Session
	LastTime  time.Time `json:"lastTime"`
	LastSeen  string `json:"lastSeen"`
}

//Notify struct
type Notify struct {
	UID          int    `json:"uid"`
	PID          int    `json:"pid"`
	CID          int    `json:"cid"`
	CLID         int    `json:"clid"`
	ID           int    `json:"id"`
	CIDPID       int    `json:"cidpid"`
	PostID       int    `json:"postId"`
	CommentID    int    `json:"commentId"`
	UserLostID   int    `json:"userLostId"`
	VoteState    int    `json:"voteState"`
	Time  string `json:"time"`
	CreatedTime  time.Time `json:"createdTime"`
	UpdatedTime time.Time `json:"updatedTime"`
	ToWhom       int    `json:"toWhom"`
	Content    string `json:"postTitle"`
	UserLost     string `json:"userLost"`
	Editted bool `json:"editted"`
}

//GetUserProfile function
func (user *User) GetUserProfile(r *http.Request, w http.ResponseWriter) ([]Post, []Post, []Post, []Comment, User) {

	//time.AfterFunc(10, checkCookieLife(cookie, w, r)) try check every 30 min cookie
	u := User{}
	liked := VotedPosts("like_state", user.Session.UserID)
	disliked := VotedPosts("dislike_state", user.Session.UserID)
	err = DB.QueryRow("SELECT id, full_name, email, username, isAdmin, age, sex, created_time, city, image, last_seen  FROM users WHERE id = ?", user.Session.UserID).Scan(&u.ID, &u.FullName, &u.Email, &u.Username, &u.IsAdmin, &u.Age, &u.Sex, &u.CreatedTime, &u.City, &u.Image, &u.LastTime)
	if err != nil {
		log.Println(err)
	}
	if u.Image[0] == 60 {
		u.SVG = true
	}
	u.LastSeen = u.LastTime.Format("2006 Jan _2 15:04:05")
	u.Temp = u.CreatedTime.Format("2006 Jan _2 15:04:05")
	u.ImageHTML = base64.StdEncoding.EncodeToString(u.Image)

	//get posts current user
	pStmp, err := DB.Query("SELECT * FROM posts WHERE creator_id=?", user.Session.UserID)
	if err != nil {
		log.Println(err.Error())
	}
	postsCreated := []Post{}

	for pStmp.Next() {
		err = pStmp.Scan(&post.ID, &post.Title, &post.Content, &post.CreatorID, &post.CreateTime, &post.UpdateTime, &post.Image, &post.Like, &post.Dislike)

		if err != nil {
			log.Println(err.Error())
		}
		post.AuthorForPost = user.Session.UserID
		post.Time = post.CreateTime.Format("2006 Jan _2 15:04:05")
		postsCreated = append(postsCreated, post)
	}

	commentQuery, err := DB.Query("SELECT * FROM comments WHERE creator_id=?", user.Session.UserID)
	if err != nil {
		log.Println(err.Error())
	}
	var comments []Comment
	var cmt Comment
	defer commentQuery.Close()

	for commentQuery.Next() {
		err = commentQuery.Scan(&cmt.ID, &cmt.ParentID, &cmt.Content, &cmt.PostID, &cmt.UserID, &cmt.ToWhom, &cmt.FromWhom, &cmt.CreatedTime, &cmt.UpdatedTime, &cmt.Like, &cmt.Dislike)
		if err != nil {
			log.Println(err.Error())
		}
		err = DB.QueryRow("SELECT post_id FROM comments WHERE id = ?", cmt.ID).Scan(&post.ID)
		if err != nil {
			log.Println(err.Error())
		}
		err = DB.QueryRow("SELECT thread FROM posts WHERE id = ?", post.ID).Scan(&cmt.TitlePost)
		if err != nil {
			log.Println(err.Error())
		}

		cmt.Time = cmt.CreatedTime.Format("2006 Jan _2 15:04:05")
		comments = append(comments, cmt)
	}

	return disliked, liked, postsCreated, comments, u
}

//GetUserActivities func
func GetUserActivities(w http.ResponseWriter, r *http.Request, s *general.Session) (result []Notify) {

	var notifies []Notify
	nQuery, err := DB.Query("SELECT * FROM notify WHERE to_whom=?", s.UserID)
	if err != nil {
		log.Println(err)
	}
	for nQuery.Next() {
		n := Notify{}
		err = nQuery.Scan(&n.ID, &n.PostID, &n.UserLostID, &n.VoteState, &n.CreatedTime, &n.ToWhom, &n.CommentID)
		if err != nil {
			log.Println(err)
		}
		notifies = append(notifies, n)
	}
	for _, v := range notifies {
		//get data/comment data by id
		n := Notify{}
		// commnentId - delete, but notify - Have row
		err = DB.QueryRow("SELECT thread, create_time, update_time FROM posts WHERE id = ?", v.PostID).Scan(&n.Content, &n.CreatedTime, &n.UpdatedTime)
		if err != nil {
			log.Println(err)
		}		
		err = DB.QueryRow("SELECT post_id, content, create_time, update_time FROM comments WHERE id = ?", v.CommentID).Scan(&n.CIDPID, &n.Content, &n.CreatedTime, &n.UpdatedTime)
		if err != nil {
			log.Println(err)
		}
		
		//compare create == update, time post/comment
		//like - to dislike, Edited vote ??
		diff := n.UpdatedTime.Sub(n.CreatedTime)
		
		if diff > 0 {
			n.Time = n.UpdatedTime.Format("2006 Jan _2 15:04:05")
			n.Editted = true
			//fmt.Println( "edited post/comment", diff)		
		}else {
			n.Time = n.CreatedTime.Format("2006 Jan _2 15:04:05")
		}

		err = DB.QueryRow("SELECT full_name FROM users WHERE id = ?", v.UserLostID).Scan(&n.UserLost)
		if err != nil {
			log.Println(err)
		}

		n.VoteState = v.VoteState
		n.UID = v.UserLostID

		fmt.Println("lost vote author:", n.UserLost)

		if v.VoteState == 1 && v.PostID != 0 {
			n.PID = v.PostID
			fmt.Println("user: ", n.UserLost, " liked your post : ", n.Content, " in ", v.CreatedTime, "")
		}
		if v.VoteState == 2 && v.PostID != 0 {
			n.PID = v.PostID
			fmt.Println("user: ", n.UserLost, " Dislike your post : ", n.Content, " in ", v.CreatedTime, "")
		}
		if v.VoteState == 1 && v.CommentID != 0 {
			n.CID = v.CommentID
			fmt.Println("user: ", n.UserLost, " liked u Comment : ", n.Content, " in ", v.CreatedTime, "")
		}

		if v.VoteState == 2 && v.CommentID != 0 {
			n.CID = v.CommentID
			fmt.Println("user: ", n.UserLost, " Dislike u Comment!: ", n.Content, " in ", v.CreatedTime, "", n.CID, n.CIDPID)
		}
		//comment lost case
		if v.VoteState == 0 && v.CommentID != 0 {
			fmt.Println("user: ", n.UserLost, " Comment u Post: ", n.Content, " in ", v.CreatedTime)
			n.CLID = v.PostID
		}
		if n.PID > 0 || n.CID > 0 || n.CLID > 0 || n.VoteState > 0 {
			result = append(result, n)
		}
	}
	return result
}

//GetAnotherProfile other user data
func (user *User) GetAnotherProfile(r *http.Request) ([]Post, User, error) {

	//userQR := DB.QueryRow("SELECT * FROM users WHERE id = ?", user.Temp)
	u := User{}
	postsU := []Post{}

	//err = userQR.Scan(&u.ID, &u.FullName, &u.Email, &u.Password, &u.IsAdmin, &u.Age, &u.Sex, &u.CreatedTime, &u.City, &u.Image)
	err = DB.QueryRow("SELECT id, full_name, email, isAdmin, age, sex, created_time, city, image  FROM users WHERE id = ?", user.Temp).Scan(&u.ID, &u.FullName, &u.Email, &u.IsAdmin, &u.Age, &u.Sex, &u.CreatedTime, &u.City, &u.Image)
	if u.Image[0] == 60 {
		u.SVG = true
	}
	u.ImageHTML = base64.StdEncoding.EncodeToString(u.Image)
	psu, err := DB.Query("SELECT * FROM posts WHERE creator_id=?", u.ID)

	defer psu.Close()

	for psu.Next() {

		err = psu.Scan(&post.ID, &post.Title, &post.Content, &post.CreatorID, &post.CreateTime, &post.UpdateTime, &post.Image, &post.Like, &post.Dislike)
		if err != nil {
			log.Println(err.Error())
		}
		//AuthorForPost
		post.Time = post.CreateTime.Format("2006 Jan _2 15:04:05")
		postsU = append(postsU, post)
	}
	if err != nil {
		return nil, u, err
	}
	return postsU, u, nil
}

//UpdateProfile function
func (u *User) UpdateProfile() {

	_, err := DB.Exec("UPDATE  users SET full_name=?, username=?, age=?, sex=?, city=?, image=? WHERE id =?",
		u.FullName, u.Username, u.Age, u.Sex, u.City, u.Image, u.ID)
	if err != nil {
	log.Println(err)	
}
}

//DeleteAccount then dlogut - delete cookie, delete lsot comment, session Db, voteState
func (u *User) DeleteAccount(w http.ResponseWriter, r *http.Request) {

	_, err = DB.Exec("DELETE FROM  session  WHERE user_id=?", u.ID)
	if err != nil {
		log.Println(err)
	}
	_, err = DB.Exec("DELETE FROM  voteState  WHERE user_id=?", u.ID)
	if err != nil {
		log.Println(err)
	}
	_, err = DB.Exec("DELETE FROM  comments  WHERE creator_id=?", u.ID)
	if err != nil {
		log.Println(err)
	}
	_, err = DB.Exec("DELETE FROM  users  WHERE id=?", u.ID)

	if err != nil {
		log.Println(err)
	}
	utils.DeleteCookie(w)
}

func VotedPosts(voteType string, uid int) (result []Post) {

	postArr := []Votes{}
	arrIDVote := []int{}

	votedPost, err := DB.Query("select post_id from voteState where user_id=? and  "+voteType+" and comment_id is null", uid, 1)
	if err != nil {
		log.Println(err)
	}
	for votedPost.Next() {
		voteLiked := Votes{}
		err = votedPost.Scan(&voteLiked.PostID)
		postArr = append(postArr, voteLiked)
	}
	defer votedPost.Close()

	for _, v := range postArr {
		arrIDVote = append(arrIDVote, v.PostID)
	}

	for _, v := range arrIDVote {
		smtp, err := DB.Query("SELECT * FROM posts WHERE id=?", v)
		if err != nil {
			log.Println(err)
		}
		p := Post{}
		for smtp.Next() {
			err = smtp.Scan(&p.ID, &p.Title, &p.Content, &p.CreatorID, &p.CreateTime, &p.UpdateTime, &p.Image, &p.Like, &p.Dislike)
			if err != nil {
				log.Println(err.Error())
			}
			p.Time = p.CreateTime.Format("2006 Jan _2 15:04:05")
			result = append(result, p)
		}
	}
	return result
}
