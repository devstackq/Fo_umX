package config

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/devstackq/ForumX/controllers"
	"github.com/devstackq/ForumX/models"
	util "github.com/devstackq/ForumX/utils"
)

type Config struct {
	BindAddr string
}

func New() *Config {
	return &Config{
		BindAddr: ":8080",
	}
}

func Init() {
	// create DB and columns
	db, err := sql.Open("sqlite3", "forumx.db")
	if err != nil {
		log.Fatalln(err)
	}
	db.Exec("PRAGMA foreign_keys=ON")

	postCategoryBridge, err := db.Prepare(`CREATE TABLE IF NOT EXISTS post_cat_bridge(id INTEGER PRIMARY KEY AUTOINCREMENT, post_id INTEGER, category TEXT, FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE )`)
	comment, err := db.Prepare(`CREATE TABLE IF NOT EXISTS comments(id	INTEGER PRIMARY KEY AUTOINCREMENT, content TEXT, post_id	INTEGER, user_idx	INTEGER, created_time	datetime,  count_like	INTEGER DEFAULT 0, count_dislike  INTEGER DEFAULT 0, CONSTRAINT fk_p FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE, CONSTRAINT fk_u FOREIGN KEY(user_idx) REFERENCES users(id) )`)
	post, err := db.Prepare(`CREATE TABLE  IF NOT EXISTS "posts" ("id"	INTEGER PRIMARY KEY AUTOINCREMENT, "title"	TEXT, "content"	TEXT, "creator_id"	INTEGER,  "created_time"	datetime, "image"	BLOB NOT NULL, "count_like"	INTEGER DEFAULT 0, "count_dislike"	INTEGER DEFAULT 0, FOREIGN KEY("creator_id") REFERENCES "users"("id") )`)
	session, err := db.Prepare(`CREATE TABLE IF NOT EXISTS "session" ("id"	INTEGER PRIMARY KEY AUTOINCREMENT, "uuid"	TEXT, "user_id"	INTEGER UNIQUE,	FOREIGN KEY("user_id") REFERENCES  "users"("id") ON DELETE CASCADE )`)
	user, err := db.Prepare(`CREATE TABLE IF NOT EXISTS "users"  ("id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, "full_name"	TEXT NOT NULL, "email"	TEXT NOT NULL UNIQUE, "password" TEXT NOT NULL, "isAdmin"	INTEGER DEFAULT 0, "age" INTEGER, "sex" TEXT, "created_time"	datetime, "city" TEXT, "image"	BLOB NOT NULL)`)
	voteState, err := db.Prepare(`CREATE  TABLE IF NOT EXISTS voteState (id INTEGER PRIMARY KEY AUTOINCREMENT,  user_id INTEGER , post_id INTEGER, comment_id INTEGER,   like_state INTEGER  DEFAULT 0, dislike_state INTEGER  DEFAULT 0, unique(post_id, user_id) FOREIGN KEY(post_id) REFERENCES  posts(id) ON DELETE CASCADE )`)

	if err != nil {
		log.Println(err)
	}

	postCategoryBridge.Exec()
	session.Exec()
	post.Exec()
	comment.Exec()
	user.Exec()
	voteState.Exec()

	//add connection - controllers/models & utils
	controllers.DB = db
	models.DB = db
	util.DB = db
	fmt.Println("Сукцесс коннект")
}
