package structure

import (
	"database/sql"
)

//general model
type Session struct {
	ID     int
	UUID   string
	UserID int
}

//general global variable
var API struct {
	Authenticated bool
	Message       string
	DB            *sql.DB
}
