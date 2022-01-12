package response

import (
	json1 "encoding/json"
	"forum/forum/internal/models"
	log "forum/forum/pkg/logger"
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

func GetPostsFromRequest(r io.Reader) ([]models.Post, error) {
	var posts []models.Post
	decoder := json1.NewDecoder(r)
	err := decoder.Decode(&posts)
	if err != nil {
		log.Error(err)
		return nil, models.ErrJSONDecoding
	}
	return posts, nil
}

func GetVoteFromRequest(r io.Reader) (*models.Vote, error) {
	voteInput := new(models.Vote)
	err := json.UnmarshalFromReader(r, voteInput)
	if err != nil {
		return nil, models.ErrJSONDecoding
	}
	return voteInput, nil
}
