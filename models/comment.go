package models

import (
	"time"
)

//Comment ID -> foreign key -> postID
type Comment struct {
	ID      int
	Content string
	PostID  int
	UserID  int
	//CreatedTime time.Time
	Author      string
	Like        int
	Dislike     int
	TitlePost   string
	CreatedTime string
}

//LeaveComment for post by id
func (c *Comment) LeaveComment() error {
	_, err := DB.Exec("INSERT INTO comments(content, post_id, user_idx, created_time) VALUES(?,?,?,?)",
		c.Content, c.PostID, c.UserID, time.Now())
	if err != nil {
		return err
	}
	return nil
}

//AppendComment helper func
func AppendComment(id int, content string, postID, userID int, createdTime time.Time, like, dislike int, titlePost string) Comment {

	comment = Comment{
		ID:          id,
		Content:     content,
		PostID:      postID,
		UserID:      userID,
		CreatedTime: createdTime.Format("2006 Jan _2 15:04:05"),
		Like:        like,
		Dislike:     dislike,
		TitlePost:   titlePost,
	}
	return comment
}
