package data

// User is an application user and all their data
type User struct {
	Username string
	Points   uint
	Password []byte
}

type UserDAL interface {
	GetUser(*User) error
}
