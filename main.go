package main

import (
	"ForumX/config"
	"ForumX/controllers"
)

func main() {
	config.Init()
	controllers.Init()
	//github signin - not empty data,  google signin - 1 - redirect profile
	//add -  all post, each post - show tags(category)s
	//try - comment under - replies comment show
	//try errors -> with gorutine
	//no row set db - fix handle
	//superflios writeheader

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
