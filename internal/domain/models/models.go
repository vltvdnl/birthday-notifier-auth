package models

type User struct {
	ID       int64
	Email    string
	Passhash []byte
}
type App struct {
	ID     int32
	Name   string
	Secret string
}
