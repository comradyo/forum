package delivery

import (
	"forum/forum/internal/models"
	"forum/forum/internal/service"

	"forum/forum/pkg/response"
	"net/http"

	"github.com/gorilla/mux"
)

const threadLogMessage = "delivery:thread:"

type ThreadDelivery struct {
	useCase service.ThreadUseCaseInterface
}

func NewThreadDelivery(useCase service.ThreadUseCaseInterface) *ThreadDelivery {
	return &ThreadDelivery{
		useCase: useCase,
	}
}

func (d *ThreadDelivery) CreateThreadPosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slugOrId := vars["slug_or_id"]
	posts, err := response.GetPostsFromRequest(r.Body)
	if err != nil {
		response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
		return
	}
	newPosts, err := d.useCase.CreateThreadPosts(slugOrId, &models.Posts{Posts: posts})
	if err != nil {
		if err == models.ErrThreadNotFound || err == models.ErrUserNotFound {
			response.SendResponse(w, http.StatusNotFound, models.Error{Message: err.Error()})
			return
		} else if err == models.ErrPostNotFound {
			response.SendResponse(w, http.StatusConflict, models.Error{Message: err.Error()})
			return
		} else {
			response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
			return
		}
	}
	if len(newPosts.Posts) == 0 {
		response.SendResponse(w, http.StatusCreated, []models.Post{})
		return
	}
	response.SendResponse(w, http.StatusCreated, newPosts.Posts)
	return
}

func (d *ThreadDelivery) GetThreadDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slugOrId := vars["slug_or_id"]
	thread, err := d.useCase.GetThreadDetails(slugOrId)
	if err != nil {
		if err == models.ErrThreadNotFound {
			response.SendResponse(w, http.StatusNotFound, models.Error{Message: err.Error()})
			return
		} else {
			response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
			return
		}
	}
	response.SendResponse(w, http.StatusOK, thread)
	return
}

func (d *ThreadDelivery) UpdateThreadDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slugOrId := vars["slug_or_id"]
	thread, err := response.GetThreadFromRequest(r.Body)
	if err != nil {
		response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
		return
	}
	updatedThread, err := d.useCase.UpdateThreadDetails(slugOrId, thread)
	if err != nil {
		if err == models.ErrThreadNotFound {
			response.SendResponse(w, http.StatusNotFound, models.Error{Message: err.Error()})
			return
		} else {
			response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
			return
		}
	}
	response.SendResponse(w, http.StatusOK, updatedThread)
	return
}

func (d *ThreadDelivery) GetThreadPosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slugOrId := vars["slug_or_id"]

	q := r.URL.Query()
	var limit string
	var since string
	var sort string
	var desc string
	if len(q["limit"]) > 0 {
		limit = q["limit"][0]
	}
	if len(q["since"]) > 0 {
		since = q["since"][0]
	}
	if len(q["sort"]) > 0 {
		sort = q["sort"][0]
	}
	if len(q["desc"]) > 0 {
		desc = q["desc"][0]
	}

	posts, err := d.useCase.GetThreadPosts(slugOrId, limit, since, sort, desc)
	if err != nil {
		if err == models.ErrThreadNotFound {
			response.SendResponse(w, http.StatusNotFound, models.Error{Message: err.Error()})
			return
		} else {
			response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
			return
		}
	}

	if len(posts.Posts) == 0 {
		response.SendResponse(w, http.StatusOK, []models.Post{})
		return
	}
	response.SendResponse(w, http.StatusOK, posts.Posts)
	return
}

func (d *ThreadDelivery) VoteForThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slugOrId := vars["slug_or_id"]
	vote, err := response.GetVoteFromRequest(r.Body)
	if err != nil {
		response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
		return
	}
	thread, err := d.useCase.VoteForThread(slugOrId, vote)
	if err != nil {
		if err == models.ErrThreadNotFound || err == models.ErrUserNotFound {
			response.SendResponse(w, http.StatusNotFound, models.Error{Message: err.Error()})
			return
		} else {
			response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
			return
		}
	}
	response.SendResponse(w, http.StatusOK, thread)
	return
}
