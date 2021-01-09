package main

import (
	"ForumX/config"
	"ForumX/controllers"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	config.Init()
	controllers.Init()

	//wsl installed

	//signup page - todo js, like signin
	//try - comment under - replies comment show
	//create post -  out main page
	//design -> another site, copy colors, etc
	//try errors -> with gorutine
	//no row set db - fix handle
	//superflios writeheader
	// Heroku deploy

	//not require, optional:
	//not delete rows in table- add field - visible, if Client delete post/comment-> filed visible false
	//save image -> local folder, no Db
	//try - create div - content editable
	//create uniq Func -> queryDb(table, ...fields string, db)
	//todo another Func add CheckMethod
	//add valid Input data, and logger -> Middleware
	//mod Name -> change github/devstackq/...
	//try - event -> add sound & confetti -Login
	//config, router refactor
	// перегрузку методов - exp.go
	// use constructor
	// use anonim func
	// use gorutine
	// func use with Interface
	//10 principe write coding
}
