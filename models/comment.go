package models

import (
	"ForumX/utils"
	"log"
	"time"
)

//Comment ID -> foreign key -> postID
type Comment struct {
	ID          int       `json:"id"`
	Content     string    `json:"content"`
	PostID      string    `json:"postId"`
	UserID      int       `json:"userId"`
	Author      string    `json:"author"`
	Like        int       `json:"like"`
	Dislike     int       `json:"dislike"`
	TitlePost   string    `json:"titlePost"`
	Time         string `json:"time"`
	CreatedTime  time.Time    `json:"createdTime"`
	UpdatedTime time.Time    `json:"updatedTime"`
	ToWhom      int       `json:"toWhom"`
	FromWhom    int       `json:"fromWhom"`
	Replied    string       `json:"replied"`
	ReplyID     int       `json:"replyId"`
	Parent int `json:"parent"`
	Children []Comment `json:"children"`
	Edited bool `json:"edited"`
	CommentID int `json:"cid"`
	ParentID int `json:"parentId"`
	RepliedContent string `json:"repliedContent"`
}


//LeaveComment for post by id
func (c *Comment) LeaveComment() (int64) {

	commentPrepare, err := DB.Prepare(`INSERT INTO comments(content, post_id, creator_id, create_time) VALUES(?,?,?,?)`)
	if err != nil {
		log.Println(err)
	}
	commentExec, err := commentPrepare.Exec(c.Content, c.PostID, c.UserID, time.Now())
	if err != nil {
		log.Println(err)
	}
	defer commentPrepare.Close()

	//commet content
	err = DB.QueryRow("SELECT creator_id FROM posts WHERE id=?", c.PostID).Scan(&c.ToWhom)
	if err != nil {
		log.Println(err)
	}
	lid, err := commentExec.LastInsertId()
	if err != nil {
		log.Println(err)
	}
	utils.SetCommentNotify(c.PostID, c.UserID, c.ToWhom, lid)
	return lid
}

//UpdateComment func
func (c *Comment) UpdateComment() {
	_, err := DB.Exec("UPDATE comments SET content=?, update_time=? WHERE id =?",
		c.Content, c.UpdatedTime, c.ID )
	if err != nil {
		log.Println(err)
	}
}

// DeleteComment func
func DeleteComment(id string) {

	_, err = DB.Exec("DELETE FROM notify  WHERE comment_id =?", id)
	if err != nil {
		log.Println(err)
	}
	_, err = DB.Exec("DELETE FROM voteState  WHERE comment_id =?", id)
	if err != nil {
		log.Println(err)
	}
	_, err = DB.Exec("DELETE FROM  comments  WHERE id =?", id)
	if err != nil {
		log.Println(err)
	}
}
