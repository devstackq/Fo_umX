package general

import "time"

//general structure -> for child packages use
//give in Browser coookie, save in server, then each call handler -> compare server uuid - Db uuid
type Session struct {
	ID          int    `json:"id"`
	UUID        string `json:"uuid"`
	UserID      int    `json:"userId"`
	DbCookie    string `json:"db_cookie"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	StartTimeCookie time.Time `json:"cookieTime"`
	Time string 
}

//general global variable
var API struct {
	Authenticated bool   `json:"authenticated"`
	Message       string `json:"message"`
}
