package domain

type User struct {
	ID       int    `json:"id" db:"id"`
	Login    string `json:"login" db:"login"`
	Password string `json:"password" db:"password"`
}
