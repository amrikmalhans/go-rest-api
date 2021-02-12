package models

//User ...
type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
	ID       uint64 `json:"ID"`
}
