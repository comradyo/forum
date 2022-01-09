package models

import "time"

type Post struct {
	Id       int64     `json:"id,omitempty"`
	Parent   int64     `json:"parent,omitempty"`
	Author   string    `json:"author"`
	Message  string    `json:"message"`
	IsEdited bool      `json:"isEdited,omitempty"`
	Forum    string    `json:"forum,omitempty"`
	Thread   int32     `json:"thread,omitempty"`
	Created  time.Time `json:"created,omitempty"`
}

type Posts struct {
	Posts []Post `json:"posts"`
}

type PostUpdate struct {
	Message string `json:"message,omitempty"`
}

type PostFull struct {
	Post   *Post   `json:"post"`
	Author *User   `json:"author"`
	Thread *Thread `json:"thread"`
	Forum  *Forum  `json:"forum"`
}
