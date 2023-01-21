package model

type PostVotes struct {
	ID     int
	UserID int
	PostID int
	Vote   int
}

type CommentVotes struct {
	ID        int
	UserID    int
	CommentID int
	Vote      int
}
