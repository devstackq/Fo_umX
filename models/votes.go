package models

import (
	"ForumX/general"
	"ForumX/utils"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

//Votes struct
type Votes struct {
	ID           int `json:"id"`
	LikeState    int `json:"likeState"`
	DislikeState int `json:"dislikeState"`
	PostID       int `json:"postId"`
	UserID       int `json:"userId"`
	CommentID    int `json:"commentId"`
	OldLike      int `json:"oldLike"`
	OldDislike   int `json:"oldDislike"`
	CreatorID    int `json:"creatorId"`
}

//VoteDislike func
func VoteDislike(w http.ResponseWriter, r *http.Request, id, any string, s *general.Session) {

	vote := Votes{}
	field := any + "_id"
	table := any + "s"

	err = DB.QueryRow("SELECT id FROM voteState WHERE "+field+"=$1 AND user_id=$2", id, s.UserID).Scan(&vote.ID)
	if err != nil {
		log.Println(err)
	}
	err = DB.QueryRow("SELECT creator_id FROM "+table+"  WHERE id=$1", id).Scan(&vote.CreatorID)
	if err != nil {
		log.Println(err)
	}
	objID, _ := strconv.Atoi(id)
	v := Votes{}

	if vote.ID == 0 {
		fmt.Println(vote.ID, s.UserID, "start", objID, table, "table", "init Dislike field", field)

		votePrepare, err := DB.Prepare(`INSERT INTO voteState(` + field + `, user_id, dislike_state, like_state) VALUES($1,$2,$3,$4)`)
		if err != nil {
			log.Println(err)
		}
		_, err = votePrepare.Exec(id, s.UserID, 1, 0)
		if err != nil {
			log.Println(err)
		}
		defer votePrepare.Close()

		err = DB.QueryRow("SELECT count_dislike FROM "+table+" WHERE id=$1", id).Scan(&vote.OldDislike)
		if err != nil {
			log.Println(err)
		}
		_, err = DB.Exec("UPDATE "+table+" SET count_dislike=$1 WHERE id=$2", vote.OldDislike+1, id)

		if err != nil {
			log.Println(err)
		}

		utils.SetVoteNotify(any, vote.CreatorID, s.UserID, objID, false)

	} else {
		err = DB.QueryRow("SELECT count_like FROM "+table+" WHERE id=$1", id).Scan(&vote.OldLike)
		if err != nil {
			log.Println(err)
		}
		err = DB.QueryRow("SELECT count_dislike FROM "+table+" WHERE id=$1", id).Scan(&vote.OldDislike)
		if err != nil {
			log.Println(err)
		}
		DB.QueryRow("SELECT like_state, dislike_state FROM voteState where "+field+"=$1 and user_id=$2", id, s.UserID).Scan(&vote.LikeState, &vote.DislikeState)

		//set dislike default
		if vote.LikeState == 0 && vote.DislikeState == 1 {

			vote.OldDislike--
			_, err = DB.Exec("UPDATE "+table+" SET count_dislike = $1 WHERE id=$2", vote.OldDislike, id)
			if err != nil {
				log.Println(err)
			}
			_, err = DB.Exec("UPDATE voteState SET  dislike_state=$1 WHERE "+field+"=$2 and user_id=$3", 0, id, s.UserID)
			if err != nil {
				log.Println(err)
			}
			//remove notify table
			//research Pointer -> Struct, and Method Struct

			v.UpdateNotify(any, vote.CreatorID, s.UserID, objID, 0)
			fmt.Println("case 2 like 0, dis 1")
		}

		//set dislike -> to like
		if vote.LikeState == 1 && vote.DislikeState == 0 {
			fmt.Println("case 3 like 1, dis 0")

			vote.OldDislike++
			vote.OldLike--
			_, err = DB.Exec("UPDATE "+table+" SET count_dislike = $1 WHERE id=$2", vote.OldDislike, id)
			if err != nil {
				log.Println(err)
			}
			_, err = DB.Exec("UPDATE "+table+" SET count_like = $1 WHERE id=$2", vote.OldLike, id)
			if err != nil {
				log.Println(err)
			}
			_, err = DB.Exec("UPDATE voteState SET like_state = $1, dislike_state=$2 WHERE "+field+"=$3 and user_id=$4", 0, 1, id, s.UserID)
			if err != nil {
				log.Println(err)
			}

			v.UpdateNotify(any, vote.CreatorID, s.UserID, objID, 2)
		}

		if vote.LikeState == 0 && vote.DislikeState == 0 {
			fmt.Println("case 4 like 0, dis 0 LS 0, DS 1")
			vote.OldDislike++
			_, err = DB.Exec("UPDATE "+table+" SET count_dislike=$1 WHERE id=$2", vote.OldDislike, id)
			if err != nil {
				log.Println(err)
			}
			_, err = DB.Exec("UPDATE voteState SET dislike_state = $1, like_state=$2 WHERE "+field+"=$3 and user_id=$4", 1, 0, id, s.UserID)
			if err != nil {
				log.Println(err)
			}
			v.UpdateNotify(any, vote.CreatorID, s.UserID, objID, 2)
		}
	}
}

//VoteLike func
func VoteLike(w http.ResponseWriter, r *http.Request, id, any string, s *general.Session) {

	vote := Votes{}
	field := any + "_id"
	table := any + "s"
	//get current UserId by uuid
	//get by post_id and user_id -> row -> in voteState, if not -> create new row, set chenge likeState = 1, add post by ID  - like_count + 1
	DB.QueryRow("SELECT id FROM voteState where "+field+"=$1 and user_id=$2", id, s.UserID).Scan(&vote.ID)

	err = DB.QueryRow("SELECT creator_id FROM "+table+"  WHERE id=$1", id).Scan(&vote.CreatorID)
	if err != nil {
		log.Println(err)
	}
	pid, _ := strconv.Atoi(id)
	v := Votes{}

	if vote.ID == 0 {
		fmt.Println(vote.ID, s.UserID, "start", id, table, "table", "init Like field", field)

		votePrepare, err := DB.Prepare("INSERT INTO voteState(" + field + ", user_id, like_state, dislike_state) VALUES($1,$2,$3,$4)")
		if err != nil {
			log.Println(err)
		}
		_, err = votePrepare.Exec(id, s.UserID, 1, 0)
		if err != nil {
			log.Println(err)
		}
		defer votePrepare.Close()

		err = DB.QueryRow("SELECT count_like FROM "+table+" WHERE id=$1", id).Scan(&vote.OldLike)
		if err != nil {
			log.Println(err)
		}
		_, err = DB.Exec("UPDATE "+table+" SET count_like=$1 WHERE id=$2", vote.OldLike+1, id)
		if err != nil {
			log.Println(err)
		}
		utils.SetVoteNotify(any, vote.CreatorID, s.UserID, pid, true)

	} else {
		//if post -> liked or Disliked -> get CountLike & Dislike current Post, and get LikeState & dislike State
		err = DB.QueryRow("SELECT count_like FROM "+table+"  WHERE id=$1", id).Scan(&vote.OldLike)
		if err != nil {
			log.Println(err)
		}
		err = DB.QueryRow("SELECT count_dislike FROM "+table+"  WHERE id=$1", id).Scan(&vote.OldDislike)
		if err != nil {
			log.Println(err)
		}

		DB.QueryRow("SELECT like_state, dislike_state FROM voteState where "+field+"=$1 and user_id=$2", id, s.UserID).Scan(&vote.LikeState, &vote.DislikeState)

		fmt.Println(" old Dislike & like", vote.OldDislike, vote.OldLike)

		//set like 0, default
		if vote.LikeState == 1 && vote.DislikeState == 0 {

			fmt.Println("l-1, d-0 -> l-0,d-0")
			vote.OldLike--
			_, err = DB.Exec("UPDATE "+table+"  SET count_like =$1 WHERE id=$2", vote.OldLike, id)
			if err != nil {
				log.Println(err)
			}
			_, err = DB.Exec("UPDATE voteState SET like_state =$1 WHERE "+field+"=$2  and user_id=$3", 0, id, s.UserID)
			if err != nil {
				log.Println(err)
			}

			v.UpdateNotify(any, vote.CreatorID, s.UserID, pid, 0)

		}
		//set dislike -> to like
		if vote.LikeState == 0 && vote.DislikeState == 1 {

			fmt.Println("l-0,d-1, -> l-1, d-0")
			vote.OldDislike--
			vote.OldLike++
			_, err = DB.Exec("UPDATE "+table+"  SET count_dislike =$1, count_like=$2 WHERE id=$3", vote.OldDislike, vote.OldLike, id)
			if err != nil {
				log.Println(err)
			}
			_, err = DB.Exec("UPDATE voteState SET like_state = $1, dislike_state=$2 WHERE "+field+"=$3  and user_id=$4", 1, 0, id, s.UserID)
			if err != nil {
				log.Println(err)
			}

			//add like notify &  remove DislikeNotify
			v.UpdateNotify(any, vote.CreatorID, s.UserID, pid, 1)
		}
		//set like,
		if vote.LikeState == 0 && vote.DislikeState == 0 {

			fmt.Println("l-0, d-0 -> l-1, d-0")
			vote.OldLike++
			_, err = DB.Exec("UPDATE "+table+" SET count_like=$1 WHERE id=$2", vote.OldLike, id)
			if err != nil {
				log.Println(err)
			}
			_, err = DB.Exec("UPDATE voteState SET like_state =$1, dislike_state=$2 WHERE "+field+"= $3  and user_id=$4", 1, 0, id, s.UserID)
			if err != nil {
				log.Println(err)
			}
			v.UpdateNotify(any, vote.CreatorID, s.UserID, pid, 1)
		}
	}
}

//difference ? default func  utils.UpdateVoteNotify
func (v *Votes) UpdateNotify(table string, toWhom, fromWhom, objID, voteType int) {

	if table == "post" && toWhom != 0 {
		_, err = DB.Exec("UPDATE notify SET voteState=$1 WHERE comment_id=$2 AND post_id=$3 AND current_user_id=$4  AND to_whom=$5", voteType, 0, objID, fromWhom, toWhom)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(objID, fromWhom, toWhom, "update  Like/Dislike Post")

	} else if table == "comment" && toWhom != 0 {

		fmt.Println(objID, fromWhom, toWhom, "notify Update Vote Comment")
		_, err = DB.Exec("UPDATE notify SET voteState=$1 WHERE post_id=$2 AND  comment_id=$3 AND current_user_id=$4  AND to_whom=$5", voteType, 0, objID, fromWhom, toWhom)
		if err != nil {
			fmt.Println(err)
		}
	}
}
