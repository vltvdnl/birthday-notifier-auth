package models

type User struct {
	ID       int64
	Email    string
	Passhash string
}
type App struct {
	ID     int32
	Name   string
	Secret string
}
