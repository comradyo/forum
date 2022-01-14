package response

import (
	json1 "encoding/json"
	"forum/internal/models"

	json "github.com/mailru/easyjson"
)

func GetForumFromRequest(r []byte) (*models.Forum, error) {
	forumInput := new(models.Forum)
	err := json.Unmarshal(r, forumInput)
	if err != nil {
		return nil, models.ErrJSONDecoding
	}
	return forumInput, nil
}

func GetProfileFromRequest(r []byte) (*models.User, error) {
	userInput := new(models.User)
	err := json.Unmarshal(r, userInput)
	if err != nil {
		return nil, models.ErrJSONDecoding
	}
	return userInput, nil
}

func GetThreadFromRequest(r []byte) (*models.Thread, error) {
	threadInput := new(models.Thread)
	err := json.Unmarshal(r, threadInput)
	if err != nil {
		return nil, models.ErrJSONDecoding
	}
	return threadInput, nil
}

func GetPostFromRequest(r []byte) (*models.Post, error) {
	postInput := new(models.Post)
	err := json.Unmarshal(r, postInput)
	if err != nil {
		return nil, models.ErrJSONDecoding
	}
	return postInput, nil
}

func GetPostsFromRequest(r []byte) ([]models.Post, error) {
	var posts []models.Post
	err := json1.Unmarshal(r, &posts)
	if err != nil {
		return nil, models.ErrJSONDecoding
	}
	return posts, nil
}

func GetVoteFromRequest(r []byte) (*models.Vote, error) {
	voteInput := new(models.Vote)
	err := json.Unmarshal(r, voteInput)
	if err != nil {
		return nil, models.ErrJSONDecoding
	}
	return voteInput, nil
}
