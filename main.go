package main

import (
	"github.com/devstackq/ForumX/config"
	"github.com/devstackq/ForumX/controllers"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	config.Init()
	controllers.Init()
}

//statrt - Auth
//try - event -> add sound & confetti -Logiin

//domen check - org, kz ru, etc
// save photo, like - source DB refactor
//config, router refactor

//if cookie = 0, notify message  user, logout etc
//обработать ошикбки, log & http errors check http etc

//google acc signin -> -> back signin ? what??
//start Auth
//google token, client id, event signin Google, -> get data User,
//Name. email, photo, -> then save Db. -> authorized Forum
// Logout event, logout system, delete cookie, logout Google
//272819090705-qu6arlmkvs66hc5fuvalv6liuf2n9fj8.apps.googleusercontent.com   || W42c6sfYqhPc4O5wXMobY3av

// 1 request, 910 additional, 0904 - 101202 ->
// 2 request -7575
// 3 request 910 additional, 090410 - 101202 ->Otegen batyr etc
