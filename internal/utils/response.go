package utils

import (
	"forum/forum/internal/models"
	json "github.com/mailru/easyjson"
	"io"
)

func GetForumFromRequest(r io.Reader) (*models.Forum, error) {
	forumInput := new(models.Forum)
	err := json.UnmarshalFromReader(r, forumInput)
	if err != nil {
		return nil, models.ErrJSONDecoding
	}
	return forumInput, nil
}

func GetProfileFromRequest(r io.Reader) (*models.User, error) {
	userInput := new(models.User)
	err := json.UnmarshalFromReader(r, userInput)
	if err != nil {
		return nil, models.ErrJSONDecoding
	}
	return userInput, nil
}

func GetThreadFromRequest(r io.Reader) (*models.Thread, error) {
	threadInput := new(models.Thread)
	err := json.UnmarshalFromReader(r, threadInput)
	if err != nil {
		return nil, models.ErrJSONDecoding
	}
	return threadInput, nil
}

func GetPostFromRequest(r io.Reader) (*models.Post, error) {
	postInput := new(models.Post)
	err := json.UnmarshalFromReader(r, postInput)
	if err != nil {
		return nil, models.ErrJSONDecoding
	}
	return postInput, nil
}

func GetPostsFromRequest(r io.Reader) (*models.Posts, error) {
	postsInput := new(models.Posts)
	err := json.UnmarshalFromReader(r, postsInput)
	if err != nil {
		return nil, models.ErrJSONDecoding
	}
	return postsInput, nil
}

func GetVoteFromRequest(r io.Reader) (*models.Vote, error) {
	voteInput := new(models.Vote)
	err := json.UnmarshalFromReader(r, voteInput)
	if err != nil {
		return nil, models.ErrJSONDecoding
	}
	return voteInput, nil
}
