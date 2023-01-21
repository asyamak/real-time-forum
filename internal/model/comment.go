package model

type Comment struct {
	ID           int         `json:"id"`
	Author       User        `json:"author"`
	PostID       int         `json:"postID"`
	Content      string      `json:"content"`
	ImagePath    string      `json:"image_path"`
	CreationTime interface{} `json:"creation_time"`
	UserRate     int         `json:"userRate"`
	Rating       int         `json:"rating"`
}
