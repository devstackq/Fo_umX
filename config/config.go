package config

import (
	"ForumX/controllers"
	"ForumX/models"
	"ForumX/utils"
	"database/sql"
	"fmt"
	"log"
)

var (
	db  *sql.DB
	err error
)

//Init Db, if not table -> create
func Init() {
	// create DB and table
	db, err = sql.Open("sqlite3", "forumx.db")
	if err != nil {
		log.Println(err)
	}
	db.Exec("PRAGMA foreign_keys=ON")

	postCategoryBridge, err := db.Prepare(`CREATE TABLE IF NOT EXISTS post_cat_bridge(id INTEGER PRIMARY KEY AUTOINCREMENT, post_id INTEGER, category_id INTEGER, FOREIGN KEY(category_id) REFERENCES category(id), FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE )`)
	if err != nil {
		log.Println(err)
	}
	comment, err := db.Prepare(`CREATE TABLE IF NOT EXISTS comments(id INTEGER PRIMARY KEY AUTOINCREMENT, parent_id INTEGER DEFAULT 0, content TEXT, post_id INTEGER, creator_id INTEGER DEFAULT 0, toWho INTEGER DEFAULT 0, fromWho INTEGER DEFAULT 0, create_time datetime,  update_time	datetime DEFAULT CURRENT_TIMESTAMP,  count_like INTEGER DEFAULT 0, count_dislike  INTEGER DEFAULT 0, CONSTRAINT fk_key_post_comment FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE )`)
	if err != nil {
		log.Println(err)
	}
	post, err := db.Prepare(`CREATE TABLE IF NOT EXISTS posts(id INTEGER PRIMARY KEY AUTOINCREMENT, thread TEXT, content TEXT, creator_id INTEGER, create_time datetime,   update_time datetime DEFAULT CURRENT_TIMESTAMP, image	BLOB NOT NULL, count_like INTEGER DEFAULT 0, count_dislike INTEGER DEFAULT 0, FOREIGN KEY(creator_id) REFERENCES users(id) ON DELETE CASCADE ) `)
	if err != nil {
		log.Println(err)
	}
	session, err := db.Prepare(`CREATE TABLE IF NOT EXISTS "session"("id" INTEGER PRIMARY KEY AUTOINCREMENT, "uuid"	TEXT, "user_id"	INTEGER UNIQUE, cookie_time datetime)`)
	if err != nil {
		log.Println(err)
	}
	user, err := db.Prepare(`CREATE TABLE IF NOT EXISTS "users"("id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, "full_name" TEXT NOT NULL, "email"	TEXT NOT NULL UNIQUE, "username" TEXT NOT NULL UNIQUE, "password" TEXT, "isAdmin" INTEGER DEFAULT 0, "age" INTEGER, "sex" TEXT, "created_time"	datetime, "last_seen" datetime, "city" TEXT, "image"	BLOB NOT NULL)`)
	if err != nil {
		log.Println(err)
	}
	voteState, err := db.Prepare(`CREATE TABLE IF NOT EXISTS voteState(id INTEGER PRIMARY KEY AUTOINCREMENT,  user_id INTEGER, post_id INTEGER, comment_id INTEGER,   like_state INTEGER  DEFAULT 0, dislike_state INTEGER  DEFAULT 0, unique(post_id, user_id), FOREIGN KEY(comment_id) REFERENCES comments(id), FOREIGN KEY(post_id) REFERENCES posts(id))`)
	if err != nil {
		log.Println(err)
	}
	notify, err := db.Prepare(`CREATE TABLE IF NOT EXISTS notify(id INTEGER PRIMARY KEY AUTOINCREMENT, post_id INTEGER,  current_user_id INTEGER, voteState INTEGER DEFAULT 0, created_time datetime, to_whom INTEGER, comment_id INTEGER )`)
	if err != nil {
		log.Println(err)
	}

	category, err := db.Prepare(`CREATE TABLE IF NOT EXISTS  category(id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT UNIQUE)`)
	if err != nil {
		log.Println(err)
	}
	postCategoryBridge.Exec()
	session.Exec()
	post.Exec()
	comment.Exec()
	user.Exec()
	voteState.Exec()
	notify.Exec()
	category.Exec()
	putCategoriesInDb()

	//send packege - DB conn
	controllers.DB = db
	models.DB = db
	utils.DB = db
	fmt.Println("Сукцесс коннект")
}

//first call -> put categories values
func putCategoriesInDb() {

	count := utils.GetCountTable("category", db)

	if count != 3 {
		categories := []string{"science", "love", "sapid"}
		for i := 0; i < 3; i++ {
			categoryPrepare, err := db.Prepare(`INSERT INTO category(name) VALUES(?)`)
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
