package model

type Post struct {
	ID           int         `json:"id"`
	Author       User        `json:"author"`
	Title        string      `json:"title"`
	Content      string      `json:"content"`
	CreationTime interface{} `json:"creation_time"`
	ImagePath    string      `json:"image_path"`
	Categories   []Category  `json:"categories"`
	Comments     []Comment   `json:"comments"`
	Rating       int         `json:"rating"`
	UserRate     int         `json:"user_rate"`
}
