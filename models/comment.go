package models

import (
	"ForumX/utils"
	"log"
	"time"
)

//Comment ID -> foreign key -> postID
type Comment struct {
	ID             int       `json:"id"`
	Content        string    `json:"content"`
	PostID         string    `json:"postId"`
	UserID         int       `json:"userId"`
	Author         string    `json:"author"`
	Like           int       `json:"like"`
	Dislike        int       `json:"dislike"`
	TitlePost      string    `json:"titlePost"`
	Time           string    `json:"time"`
	CreatedTime    time.Time `json:"createdTime"`
	UpdatedTime    time.Time `json:"updatedTime"`
	ToWhom         int       `json:"toWhom"`
	FromWhom       int       `json:"fromWhom"`
	Replied        string    `json:"replied"`
	ReplyID        int       `json:"replyId"`
	Parent         int       `json:"parent"`
	Children       []Comment `json:"children"`
	Edited         bool      `json:"edited"`
	CommentID      int       `json:"cid"`
	ParentID       int       `json:"parentId"`
	RepliedContent string    `json:"repliedContent"`
}

//LeaveComment for post by id
func (c *Comment) LeaveComment() int64 {

	var lid int64
	err = DB.QueryRow("INSERT INTO comments(content, post_id, creator_id, create_time) VALUES($1,$2,$3,$4) RETURNING id", c.Content, c.PostID, c.UserID, time.Now()).Scan(&lid)
	if err != nil {
		log.Println(err)
	}

	//commet content
	err = DB.QueryRow("SELECT creator_id FROM posts WHERE id=$1", c.PostID).Scan(&c.ToWhom)
	if err != nil {
		log.Println(err)
	}

	if err != nil {
		log.Println(err)
	}
	utils.SetCommentNotify(c.PostID, c.UserID, c.ToWhom, lid)
	return lid
}

//UpdateComment func
func (c *Comment) UpdateComment() {
	_, err := DB.Exec("UPDATE comments SET content=$1, update_time=$2 WHERE id=$3",
		c.Content, c.UpdatedTime, c.ID)
	if err != nil {
		log.Println(err)
	}
}

// DeleteComment func
func DeleteComment(id string) {

	_, err = DB.Exec("DELETE FROM notify  WHERE comment_id=$1", id)
	if err != nil {
		log.Println(err)
	}
	_, err = DB.Exec("DELETE FROM voteState  WHERE comment_id =$1", id)
	if err != nil {
		log.Println(err)
	}
	_, err = DB.Exec("DELETE FROM  comments  WHERE id =$1", id)
	if err != nil {
		log.Println(err)
	}
}
