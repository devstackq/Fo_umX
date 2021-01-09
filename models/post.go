package models

import (
	"ForumX/general"
	"ForumX/utils"
	"bytes"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"
)

//global variable for package models
var (
	err     error
	DB      *sql.DB
	rows    *sql.Rows
	post    Post
	msg     = general.API.Message
	pageNum = 1
)

//Posts struct
type Post struct {
	ID            int              `json:"id"`
	Title         string           `json:"title"`
	Content       string           `json:"content"`
	CreatorID     int              `json:"creatorId"`
	CreateTime    time.Time        `json:"createTime"`
	UpdateTime    time.Time        `json:"updateTime"`
	Endpoint      string           `json:"endpoint"`
	FullName      string           `json:"fullName"`
	Image         []byte           `json:"image"`
	ImageHTML     string           `json:"imageHtml"`
	PostIDEdit    int              `json:"postEditId"`
	AuthorForPost int              `json:"authorPost"`
	Like          int              `json:"like"`
	Dislike       int              `json:"dislike"`
	SVG           bool             `json:"svg"`
	PBGID         int              `json:"pbgId"`
	PBGPostID     int              `json:"pbgPostId"`
	PBGCategory   string           `json:"pbgCategory"`
	FileS         multipart.File   `json:"fileS"`
	FileI         multipart.File   `json:"fileB"`
	Session       *general.Session `json:"session"`
	Categories    []string         `json:"categories"`
	Temp          string           `json:"temp"`
	IsPhoto       bool             `json:"isPhoto"`
	Time          string           `json:"time"`
	CountPost     int              `json:"countPost"`
	Authenticated bool             `json:"isAuth"`
	Edited        bool             `json:"edited"`
}

//PostCategory struct
type PostCategory struct {
	PostID   int64  `json:"postId"`
	Category string `json:"category"`
}

//Filter struct
type Filter struct {
	Category string `json:"cateogry"`
	Like     string `json:"like"`
	Date     string `json:"date"`
}

//GetAllPost function
func (filter *Filter) GetAllPost(r *http.Request, next, prev string) ([]Post, string, string) {
	//pageNum = 1
	var post Post
	var leftJoin bool
	var arrPosts []Post

	//each call +1
	if next == "next" {
		pageNum++
	}
	if prev == "prev" {
		pageNum--
	}

	//count pageNum, fix
	limit := 4
	offset := limit * (pageNum - 1)
	switch r.URL.Path {
	case "/":
		leftJoin = false
		post.Endpoint = "/"
		if filter.Date == "asc" {
			rows, err = DB.Query("SELECT * FROM posts  ORDER BY create_time ASC LIMIT 8 ")
			if err != nil {
				log.Println(err)
			}
		} else if filter.Date == "desc" {
			rows, err = DB.Query("SELECT * FROM posts  ORDER BY create_time DESC LIMIT 8")
			if err != nil {
				log.Println(err)
			}
		} else if filter.Like == "like" {
			rows, err = DB.Query("SELECT * FROM posts  ORDER BY count_like DESC LIMIT 8")
			if err != nil {
				log.Println(err)
			}
		} else if filter.Like == "dislike" {
			rows, err = DB.Query("SELECT * FROM posts  ORDER BY count_dislike DESC LIMIT 8")
			if err != nil {
				log.Println(err)
			}
		} else if filter.Category != "" {
			leftJoin = true
			rows, err = DB.Query("SELECT  * FROM posts  LEFT JOIN post_cat_bridge  ON post_cat_bridge.post_id = posts.id   WHERE category_id=? ORDER  BY create_time  DESC LIMIT 8", filter.Category)
			if err != nil {
				log.Println(err)
			}
		} else {
			rows, err = DB.Query("SELECT * FROM posts ORDER BY create_time DESC LIMIT ? OFFSET ?", limit, offset)
			if err != nil {
				log.Println(err)
			}
		}
	case "/science":
		leftJoin = true
		post.Temp = "Science"
		post.Endpoint = "/science"
		rows, err = DB.Query("SELECT * FROM posts  LEFT JOIN post_cat_bridge  ON post_cat_bridge.post_id = posts.id   WHERE category_id=?  ORDER  BY create_time  DESC LIMIT 5", 1)
		if err != nil {
			log.Println(err)
		}
	case "/love":
		leftJoin = true
		post.Temp = "Love"
		post.Endpoint = "/love"
		rows, err = DB.Query("SELECT  * FROM posts  LEFT JOIN post_cat_bridge  ON post_cat_bridge.post_id = posts.id  WHERE category_id=?   ORDER  BY create_time  DESC LIMIT 5", 2)
		if err != nil {
			log.Println(err)
		}
	case "/sapid":
		leftJoin = true
		post.Temp = "Sapid"
		post.Endpoint = "/sapid"
		rows, err = DB.Query("SELECT  * FROM posts  LEFT JOIN post_cat_bridge  ON post_cat_bridge.post_id = posts.id  WHERE category_id=?  ORDER  BY create_time  DESC LIMIT 5", 3)
		if err != nil {
			log.Println(err)
		}
	}
	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	for rows.Next() {
		if leftJoin {
			if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatorID, &post.CreateTime, &post.UpdateTime, &post.Image, &post.Like, &post.Dislike, &post.PBGID, &post.PBGPostID, &post.PBGCategory); err != nil {
				fmt.Println(err)
			}
		} else {
			if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatorID, &post.CreateTime, &post.UpdateTime, &post.Image, &post.Like, &post.Dislike); err != nil {
				fmt.Println(err)
			}
		}

		//send countr +1
		err = DB.QueryRow("SELECT COUNT(id) FROM posts").Scan(&post.CountPost)
		if err != nil {
			log.Println(err)
		}
		post.Time = post.CreateTime.Format("2006 Jan _2 15:04:05")
		arrPosts = append(arrPosts, post)
	}
	//err = DB.QueryRow("SELECT COUNT(id) FROM posts").Scan(&post.CountPost)
	return arrPosts, post.Endpoint, post.Temp
}

//UpdatePost fucntion
func (p *Post) UpdatePost() {
	_, err := DB.Exec("UPDATE  posts SET thread=?, content=?, image=?, update_time=? WHERE id=?",
		p.Title, p.Content, p.Image, p.UpdateTime, p.ID)
	if err != nil {
		log.Println(err)
	}
}

//DeletePost function, delete rows, notify, voteState, comment, by postId
func (p *Post) DeletePost() {

	_, err = DB.Exec("DELETE FROM voteState  WHERE post_id =?", p.ID)
	if err != nil {
		log.Println(err, "3")
	}
	_, err := DB.Exec("DELETE FROM comments  WHERE post_id =?", p.ID)
	if err != nil {
		log.Println(err, "1")
	}
	
	_, err = DB.Exec("DELETE FROM notify  WHERE post_id =?", p.ID)
	if err != nil {
		log.Println(err, "2")
	}

	_, err = DB.Exec("DELETE FROM post_cat_bridge  WHERE post_id =?", p.ID)
	if err != nil {
		log.Println(err, "4")
	}
	_, err = DB.Exec("DELETE FROM posts WHERE id =?", p.ID)
	if err != nil {
		log.Println(err, "5")
	}
}

//CreatePost function
func (p *Post) CreatePost(w http.ResponseWriter, r *http.Request) {

	var fileBytes []byte
	var buff bytes.Buffer

	if p.IsPhoto {
		fileSize, _ := buff.ReadFrom(p.FileS)
		defer p.FileS.Close()

		if fileSize < 20000000 {
			if err != nil {
				log.Fatal(err)
			}
			fileBytes, err = ioutil.ReadAll(p.FileI)
		} else {
			utils.RenderTemplate(w, "header", utils.IsAuth(r))
			utils.RenderTemplate(w, "create", "Large file, more than 20mb")
		}
	} else {
		//set empty photo post
		fileBytes = []byte{0, 0}
	}

	//check empty values
	if utils.IsValidLetter(p.Title, "post") && utils.IsValidLetter(p.Content, "post") {

		createPostPrepare, err := DB.Prepare(`INSERT INTO posts(thread, content, creator_id, create_time, image) VALUES(?,?,?,?,?)`)
		if err != nil {
			log.Println(err)
		}
		
		createPostExec, err := createPostPrepare.Exec(p.Title, p.Content, p.Session.UserID, time.Now(), fileBytes)
		if err != nil {
			log.Println(err)
		}
		defer createPostPrepare.Close()

		last, err := createPostExec.LastInsertId()

		if err != nil {
			log.Println(err)
		}
		pcb := PostCategory{}
		//set def category
		if len(p.Categories) == 0 {
			pcb = PostCategory{
				PostID:   last,
				Category: "3",
			}
			pcb.CreateBridge()

		} else if len(p.Categories) == 1 {
			pcb = PostCategory{
				PostID:   last,
				Category: p.Categories[0],
			}
			pcb.CreateBridge()

		} else if len(p.Categories) > 1 {
			//loop add > 1 category post
			for _, v := range p.Categories {
				pcb = PostCategory{
					PostID:   last,
					Category: v,
				}
				pcb.CreateBridge()
			}
		}
		http.Redirect(w, r, "/post?id="+strconv.Itoa(int(last)), 302)
	} else {
		msg = "Empty title or content"
		utils.RenderTemplate(w, "header", utils.IsAuth(r))
		utils.RenderTemplate(w, "create_post", &msg)
	}
}

//GetPostByID function take from all post, only post by id, then write p struct Post
func (post *Post) GetPostByID(r *http.Request) ([]Comment, Post) {
	//write new structure data
	p := Post{}
	err = DB.QueryRow("SELECT * FROM posts WHERE id = ?", post.ID).Scan(&p.ID, &p.Title, &p.Content, &p.CreatorID, &p.CreateTime, &p.UpdateTime, &p.Image, &p.Like, &p.Dislike)
	if err != nil {
		log.Println(err)
	}
	//[]byte -> encode string, client render img base64
	//check svg || jpg,png
	if len(p.Image) > 0 {
		if p.Image[0] == 60 {
			p.SVG = true
		}
	}
	//difference time -> ,mean edited
	diff := p.UpdateTime.Sub(p.CreateTime)
	if diff > 0 {
		p.Time = p.UpdateTime.Format("2006 Jan _2 15:04:05")
		p.Edited = true
	} else {
		p.Time = p.CreateTime.Format("2006 Jan _2 15:04:05")
	}

	p.ImageHTML = base64.StdEncoding.EncodeToString(p.Image)
	DB.QueryRow("SELECT full_name FROM users WHERE id = ?", p.CreatorID).Scan(&p.FullName)

	cmtq, err := DB.Query("SELECT * FROM comments WHERE  post_id=?", p.ID)
	if err != nil {
		log.Fatal(err)
	}
	defer cmtq.Close()
	//write each fields inside Comment struct -> then  append Array Comments
	var comments []Comment
	var was []int

	for cmtq.Next() {
		//get each comment Post -> by Id, -> get each replyComment by comment_id -> get replyAnswer by reply_com_id
		comment := Comment{}
		err = cmtq.Scan(&comment.ID, &comment.ParentID, &comment.Content, &comment.PostID, &comment.UserID, &comment.ToWhom, &comment.FromWhom, &comment.CreatedTime, &comment.UpdatedTime, &comment.Like, &comment.Dislike)
		if err != nil {
			log.Println(err.Error())
		}
		diff := comment.UpdatedTime.Sub(comment.CreatedTime)
		if diff > 0 {
			comment.Time = comment.UpdatedTime.Format("2006 Jan _2 15:04:05")
			comment.Edited = true
		} else {
			comment.Time = comment.CreatedTime.Format("2006 Jan _2 15:04:05")
		}
		
		err = 	DB.QueryRow("SELECT full_name FROM users WHERE id = ?", comment.UserID).Scan(&comment.Author)
		if err != nil {
			log.Println(err)
		}

		cmtReplq, err := DB.Query("SELECT * FROM comments WHERE  post_id=?", p.ID)
		if err != nil {
			log.Fatal(err)
		}
		defer cmtReplq.Close()
		//write each fields inside Comment struct -> then  append Array Comments
		//var comments []Comment
	
		for cmtReplq.Next() {

			comRep := Comment{}
			err = cmtReplq.Scan(&comRep.ID, &comRep.ParentID, &comRep.Content, &comRep.PostID, &comRep.UserID, &comRep.ToWhom, &comRep.FromWhom, &comRep.CreatedTime, &comRep.UpdatedTime, &comRep.Like, &comRep.Dislike)
			if err != nil {
				log.Fatal(err)
			}
		//9, 10 -> 8, iskluchit 9,10
			if comRep.ParentID == comment.ID {
				comment.Children = append(comment.Children,  comRep)
				was = append(was, comRep.ID)
			}
	}
	//write - in comments [], uniqum comments, add toggle
	if comment.ParentID  > 0 {

			err = 	DB.QueryRow("SELECT creator_id, toWho FROM comments WHERE id = ?", comment.ID).Scan(&comment.FromWhom, &comment.ToWhom)
			if err != nil {
				log.Println(err)
			}
			err = 	DB.QueryRow("SELECT full_name FROM users WHERE id = ?", comment.FromWhom).Scan(&comment.Author)
			if err != nil {
				log.Println(err)
			}
			err = 	DB.QueryRow("SELECT full_name FROM users WHERE id = ?", comment.ToWhom).Scan(&comment.Replied)
			if err != nil {
				log.Println(err)
			}
			err = 	DB.QueryRow("SELECT content FROM comments WHERE id = ?", comment.ParentID).Scan(&comment.RepliedContent)
			if err != nil {
				log.Println(err)
			}
		}
		
		//write - uniq by comment.ID
		fmt.Println(was)

		comments = append(comments, comment)
	}
	//8, 11 if current Comment.ID - have, inner -  comment.ParentID, Exclude this comment

	//if comment.ID == parentID, recursive append
	// for _, com := range comments {
	// 	for _, child := range com.Children {
	// 		// if RecursiveAppend(&com, child.ID) {
	// 		// 	res = append(res, com)
	// 		// }
	// 	}
	// }

	//check, comapre comment.ID == Children.ID, || add 1 table

	return comments, p
}


//CreateBridge create post  -> post_id relation category
func (pcb *PostCategory) CreateBridge() {

	createBridgePrepare, err := DB.Prepare(`INSERT INTO post_cat_bridge(post_id, category_id) VALUES (?,?)`)
	if err != nil {
		log.Println(err)
	}
	_, err = createBridgePrepare.Exec(pcb.PostID, pcb.Category)
	if err != nil {
		log.Println(err)
	}
	defer createBridgePrepare.Close()
}

//Search post by contain title
func Search(w http.ResponseWriter, r *http.Request) []Post {

	var posts []Post
	psbt, err := DB.Query("SELECT * FROM posts WHERE thread LIKE ?", "%"+r.FormValue("search")+"%")
	if err != nil {
		log.Println(err, "not find by thread")
		return nil
	}
	psbc, err := DB.Query("SELECT * FROM posts WHERE content LIKE ?", "%"+r.FormValue("search")+"%")
	if err != nil {
		log.Println(err)
		return nil
	}
	
	defer psbc.Close()
	defer psbt.Close()
	var  pTID int
	
	for psbt.Next() {
		err = psbt.Scan(&post.ID, &post.Title, &post.Content, &post.CreatorID, &post.CreateTime, &post.UpdateTime, &post.Image, &post.Like, &post.Dislike)
		if err != nil {
			log.Println(err)
		}
		pTID = post.ID
		post.Time = post.CreateTime.Format("2006 Jan _2 15:04:05")
		posts = append(posts, post)
	}

	for psbc.Next() {
		err = psbc.Scan(&post.ID, &post.Title, &post.Content, &post.CreatorID, &post.CreateTime, &post.UpdateTime, &post.Image, &post.Like, &post.Dislike)
		if err != nil {
			log.Println(err)
		}
		//check duplicate id - appen 1 item
		post.Time = post.CreateTime.Format("2006 Jan _2 15:04:05")
		if pTID != post.ID {
			posts = append(posts, post)
		}
	}
	return posts
}
