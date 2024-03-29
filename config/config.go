package config

import (
	"ForumX/controllers"
	"ForumX/models"
	"ForumX/utils"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var (
	db  *sql.DB
	DBQ *sql.DB
	err error
)

//Init Db, if not table -> create
func Init() {

	db, err := sql.Open("postgres", "postgres://myxaiifqsidmff:09843ef69b080501b4e1b5ab2c977f31372715e3d33113d9e07dff80109b325c@ec2-54-75-199-252.eu-west-1.compute.amazonaws.com:5432/d5hu4u82tl17ms")
	if err != nil {
		log.Println("can't connect inDb")
	}

	err = db.Ping()
	if err != nil {
		log.Println("can't Ping")
	}
	//db.Exec("PRAGMA foreign_keys=ON")

	user, err := db.Prepare("CREATE TABLE IF NOT EXISTS users(id SERIAL NOT NULL PRIMARY KEY, full_name varchar(255) NOT NULL, email	varchar(255) NOT NULL UNIQUE, username varchar(255) NOT NULL UNIQUE, password varchar(255), isAdmin int DEFAULT 0, age int, sex varchar(255), created_time	timestamp, last_seen timestamp, city varchar(255), image bytea);")
	if err != nil {
		log.Println(err)
	}

	_, err = user.Exec()

	if err != nil {
		log.Println(err)
	}

	post, err := db.Prepare("CREATE TABLE IF NOT EXISTS posts(id serial PRIMARY KEY, thread text, content text, creator_id int, create_time timestamp,   update_time timestamp DEFAULT current_timestamp, image	bytea NOT NULL, count_like int DEFAULT 0, count_dislike int DEFAULT 0);")

	if err != nil {
		log.Println(err)
	}
	_, err = post.Exec()

	if err != nil {
		log.Println(err)
	}

	postCategoryBridge, err := db.Prepare("CREATE TABLE IF NOT EXISTS post_cat_bridge(id SERIAL PRIMARY KEY, post_id int, category_id int);")
	if err != nil {
		log.Println(err)
	}
	_, err = postCategoryBridge.Exec()
	if err != nil {
		log.Println(err)
	}

	//postCategoryBridge, err := db.Prepare("CREATE TABLE IF NOT EXISTS post_cat_bridge(id SERIAL PRIMARY KEY, post_id int, category_id int, CONSTRAINT fk_key_pcb_cat FOREIGN KEY(category_id) REFERENCES category(id), CONSTRAINT fk_key_pcb_post FOREIGN KEY(post_id) REFERENCES posts(id) )")
	comment, err := db.Prepare("CREATE TABLE IF NOT EXISTS comments(id SERIAL PRIMARY KEY, parent_id int DEFAULT 0, content text, post_id int, creator_id int DEFAULT 0, toWho int DEFAULT 0, fromWho int DEFAULT 0, create_time timestamp,  update_time	timestamp DEFAULT CURRENT_TIMESTAMP,  count_like int DEFAULT 0, count_dislike  int DEFAULT 0, CONSTRAINT fk_key_comment_post  FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE );")
	if err != nil {
		log.Println(err)
	}

	_, err = comment.Exec()
	if err != nil {
		log.Println(err)
	}

	session, err := db.Prepare("CREATE TABLE IF NOT EXISTS session(id SERIAL PRIMARY KEY, uuid	varchar(255), user_id int UNIQUE, cookie_time timestamp);")
	if err != nil {
		log.Println(err)
	}
	_, err = session.Exec()
	if err != nil {
		log.Println(err)
	}

	voteState, err := db.Prepare("CREATE TABLE IF NOT EXISTS voteState(id SERIAL PRIMARY KEY,  user_id int, post_id int, comment_id int,   like_state int  DEFAULT 0, dislike_state int  DEFAULT 0, unique(post_id, user_id),CONSTRAINT fk_key_vote_comment  FOREIGN KEY(comment_id) REFERENCES comments(id), CONSTRAINT fk_key_vote_post FOREIGN KEY(post_id) REFERENCES posts(id));")
	if err != nil {
		log.Println(err)
	}
	_, err = voteState.Exec()
	if err != nil {
		log.Println(err)
	}

	notify, err := db.Prepare("CREATE TABLE IF NOT EXISTS notify(id SERIAL PRIMARY KEY, post_id int,  current_user_id int, voteState int DEFAULT 0, created_time timestamp, to_whom int, comment_id int );")

	if err != nil {
		log.Println(err)
	}
	_, err = notify.Exec()

	if err != nil {
		log.Println(err)
	}

	category, err := db.Prepare("CREATE TABLE IF NOT EXISTS  category(id SERIAL PRIMARY KEY, name varchar(255) UNIQUE );")
	if err != nil {
		log.Println(err)
	}
	_, err = category.Exec()
	if err != nil {
		log.Println(err)
	}

	putCategoriesInDb(db)

	//send packege - DB conn
	controllers.DB = db
	models.DB = db
	utils.DB = db

	fmt.Println("Сукцесс коннект")
}

//first call -> put categories values
func putCategoriesInDb(db *sql.DB) {

	count := utils.GetCountTable("category", db)

	if count != 3 {
		categories := []string{"science", "love", "sapid"}
		for i := 0; i < 3; i++ {
			categoryPrepare, err := db.Prepare(`INSERT INTO category(name) VALUES($1)`)
			if err != nil {
				log.Println(err)
			}
			_, err = categoryPrepare.Exec(categories[i])
			if err != nil {
				log.Println(err)
			}
			defer categoryPrepare.Close()
		}
	}
}
