package models

import "time"

type Post struct {
	Id       int64     `json:"id"`
	Parent   int64     `json:"parent"`
	Author   string    `json:"author"`
	Message  string    `json:"message"`
	IsEdited bool      `json:"isEdited"`
	Forum    string    `json:"forum"`
	Thread   int32     `json:"thread"`
	Created  time.Time `json:"created"`
}

type Posts struct {
	Posts []Post `json:"posts"`
}

type PostFull struct {
	Post   *Post   `json:"post"`
	Author *User   `json:"author"`
	Thread *Thread `json:"thread"`
	Forum  *Forum  `json:"forum"`
}
