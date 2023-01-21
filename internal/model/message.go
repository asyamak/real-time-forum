package model

type Message struct {
	ID           int         `json:"id"`
	SenderID     int         `json:"sender_id"`
	RecipientID  int         `json:"recipient_id"`
	Message      string      `json:"message"`
	CreationTime interface{} `json:"creation_time"`
	Readed       bool        `json:"readed"`
}
