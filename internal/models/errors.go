package models

import "errors"

var (
	ErrJSONDecoding          = errors.New("json decoding error")
	ErrUserExists            = errors.New("user already exists")
	ErrUserNotFound          = errors.New("user not found")
	ErrProfileUpdateConflict = errors.New("can't update profile")
	ErrForumExists           = errors.New("forum already exists")
	ErrForumNotFound         = errors.New("forum not found")
	ErrThreadExists          = errors.New("thread already exists")
	ErrThreadNotFound        = errors.New("thread not found")
	ErrPostNotFound          = errors.New("post not found")
)
