package model

type User struct {
	ID           int         `json:"id"`
	Email        string      `json:"email"`
	Username     string      `json:"username"`
	Password     string      `json:"-"`
	FirstName    string      `json:"firstName"`
	LastName     string      `json:"lastName"`
	Age          int         `json:"age"`
	Gender       string      `json:"gender"`
	CreationTime interface{} `json:"registered"`
	Avatar       string      `json:"avatar"`
}
