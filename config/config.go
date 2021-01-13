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
	err error
)

//Init Db, if not table -> create
func Init() {

	// retrieve the url
	//	dbURL := os.Getenv("postgres://iczkybfluwphwj:12d122053793fe4ba376b339f5911d6a6cdfa16836b8e5068bfb904adfb0b2ad@ec2-52-30-161-203.eu-west-1.compu")7
	// connect to the db
	db, err := sql.Open("postgres", "postgres://rhwyoybdcpfqge:91b7f85d5fe2999acedec578e377dc63a941a0e8a320f6092b4071c4eec85b72@ec2-34-248-148-63.eu-west-1.compute.amazonaws.com:5432/dcmpipelt02b2h")
	if err != nil {
		log.Println("can't connect inDb")
	}
	fmt.Println(db, "db psq data")
	err = db.Ping()
	if err != nil {
		log.Println("can't Ping")
	}

	// create DB and table
	//db, err = sql.Open("sqlite3", "forumx.db")
	// if err != nil {
	// 	log.Println(err)
	// }
	//db.Exec("PRAGMA foreign_keys=ON")

	// CREATE TABLE IF NOT EXISTS app_user (
	// 	username varchar(45) NOT NULL,
	// 	password varchar(450) NOT NULL,
	// 	enabled integer NOT NULL DEFAULT '1',
	// 	PRIMARY KEY (username)
	// )

	// if _, err := db.Exec("CREATE TABLE IF NOT EXISTS ticks (tick timestamp)"); err != nil {
	// 	fmt.Println("Error creating database table: %q", err)
	// 	return
	// }

	post, err := db.Prepare("CREATE TABLE IF NOT EXISTS posts(id serial PRIMARY KEY, thread varchar(255), content varchar(255), creator_id integer, create_time datetime,   update_time datetime DEFAULT current_timestamp, image	bytea NOT NULL, count_like INTEGER DEFAULT 0, count_dislike INTEGER DEFAULT 0, FOREIGN KEY(creator_id) REFERENCES users(id) ON DELETE CASCADE ); ")
	if err != nil {
		log.Println(err, "1")
	}
	post.Exec()

	postCategoryBridge, err := db.Prepare("CREATE TABLE IF NOT EXISTS post_cat_bridge(id SERIAL PRIMARY KEY, post_id INTEGER, category_id INTEGER);")
	if err != nil {
		log.Println(err, "2")
	}
	pcb, err := postCategoryBridge.Exec()
	if err != nil {
		log.Println(err, "exec err 2")
	}
	fmt.Println(pcb, "pcb")
	//postCategoryBridge, err := db.Prepare("CREATE TABLE IF NOT EXISTS post_cat_bridge(id SERIAL PRIMARY KEY, post_id INTEGER, category_id INTEGER, CONSTRAINT fk_key_pcb_cat FOREIGN KEY(category_id) REFERENCES category(id), CONSTRAINT fk_key_pcb_post FOREIGN KEY(post_id) REFERENCES posts(id) )")
	comment, err := db.Prepare("CREATE TABLE IF NOT EXISTS comments(id SERIAL PRIMARY KEY, parent_id INTEGER DEFAULT 0, content varchar(255), post_id INTEGER, creator_id INTEGER DEFAULT 0, toWho INTEGER DEFAULT 0, fromWho INTEGER DEFAULT 0, create_time datetime,  update_time	datetime DEFAULT CURRENT_TIMESTAMP,  count_like INTEGER DEFAULT 0, count_dislike  INTEGER DEFAULT 0, CONSTRAINT fk_key_post_comment FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE );")
	if err != nil {
		log.Println(err, "3")
	}
	comment.Exec()

	session, err := db.Prepare("CREATE TABLE IF NOT EXISTS session(id SERIAL PRIMARY KEY, uuid	varchar(255), user_id INTEGER UNIQUE, cookie_time datetime);")
	if err != nil {
		log.Println(err, "4")
	}
	session.Exec()

	user, err := db.Prepare("CREATE TABLE IF NOT EXISTS users(id SERIAL NOT NULL PRIMARY KEY, full_name varchar(255) NOT NULL, email	varchar(255) NOT NULL UNIQUE, username varchar(255) NOT NULL UNIQUE, password varchar(255), isAdmin INTEGER DEFAULT 0, age INTEGER, sex varchar(255), created_time	datetime, last_seen datetime, city varchar(255), image	bytea NOT NULL);")
	if err != nil {
		log.Println(err, "5")
	}
	user.Exec()

	voteState, err := db.Prepare("CREATE TABLE IF NOT EXISTS voteState(id SERIAL PRIMARY KEY,  user_id INTEGER, post_id INTEGER, comment_id INTEGER,   like_state INTEGER  DEFAULT 0, dislike_state INTEGER  DEFAULT 0, unique(post_id, user_id), FOREIGN KEY(comment_id) REFERENCES comments(id), FOREIGN KEY(post_id) REFERENCES posts(id));")
	if err != nil {
		log.Println(err, "6")
	}
	voteState.Exec()

	notify, err := db.Prepare("CREATE TABLE IF NOT EXISTS notify(id SERIAL PRIMARY KEY, post_id INTEGER,  current_user_id INTEGER, voteState INTEGER DEFAULT 0, created_time datetime, to_whom INTEGER, comment_id INTEGER );")
	if err != nil {
		log.Println(err, "7")
	}
	notify.Exec()

	category, err := db.Prepare("CREATE TABLE IF NOT EXISTS  category(id SERIAL PRIMARY KEY, name varchar(255) UNIQUE) );")
	if err != nil {
		log.Println(err)
	}
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
	log.Println("test")
}
