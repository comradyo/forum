package delivery

import (
	"forum/forum/internal/models"
	"forum/forum/internal/service"

	"forum/forum/pkg/response"
	"net/http"

	"github.com/gorilla/mux"
)

const forumLogMessage = "delivery:forum:"

type ForumDelivery struct {
	useCase service.ForumUseCaseInterface
}

func NewForumDelivery(useCase service.ForumUseCaseInterface) *ForumDelivery {
	return &ForumDelivery{
		useCase: useCase,
	}
}

func (d *ForumDelivery) CreateForum(w http.ResponseWriter, r *http.Request) {
	forum, err := response.GetForumFromRequest(r.Body)
	if err != nil {
		response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
		return
	}
	newForum, err := d.useCase.CreateForum(forum)
	if err != nil {
		if err == models.ErrForumExists {
			response.SendResponse(w, http.StatusConflict, newForum)
			return
		} else if err == models.ErrUserNotFound {
			response.SendResponse(w, http.StatusNotFound, models.Error{Message: err.Error()})
			return
		} else {
			response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
			return
		}
	}
	response.SendResponse(w, http.StatusCreated, newForum)
	return
}

func (d *ForumDelivery) GetForumDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]
	forum, err := d.useCase.GetForumDetails(slug)
	if err != nil {
		if err == models.ErrForumNotFound {
			response.SendResponse(w, http.StatusNotFound, models.Error{Message: err.Error()})
			return
		} else {
			response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
			return
		}
	}
	response.SendResponse(w, http.StatusOK, forum)
	return
}

func (d *ForumDelivery) CreateForumThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]
	thread, err := response.GetThreadFromRequest(r.Body)
	if err != nil {
		response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
		return
	}
	newThread, err := d.useCase.CreateForumThread(slug, thread)
	if err != nil {
		if err == models.ErrThreadExists {
			response.SendResponse(w, http.StatusConflict, newThread)
			return
		} else if err == models.ErrUserNotFound || err == models.ErrForumNotFound {
			response.SendResponse(w, http.StatusNotFound, models.Error{Message: err.Error()})
			return
		} else {
			response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
			return
		}
	}
	response.SendResponse(w, http.StatusCreated, newThread)
	return
}

func (d *ForumDelivery) GetForumUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	q := r.URL.Query()
	var limit string
	var since string
	var desc string
	if len(q["limit"]) > 0 {
		limit = q["limit"][0]
	}
	if len(q["since"]) > 0 {
		since = q["since"][0]
	}
	if len(q["desc"]) > 0 {
		desc = q["desc"][0]
	}

	users, err := d.useCase.GetForumUsers(slug, limit, since, desc)
	if err != nil {
		if err == models.ErrForumNotFound {
			response.SendResponse(w, http.StatusNotFound, models.Error{Message: err.Error()})
			return
		} else {
			response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
			return
		}
	}
	if len(users.Users) == 0 {
		response.SendResponse(w, http.StatusOK, []models.User{})
		return
	}
	response.SendResponse(w, http.StatusOK, users.Users)
	return
}

func (d *ForumDelivery) GetForumThreads(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	q := r.URL.Query()
	var limit string
	var since string
	var desc string
	if len(q["limit"]) > 0 {
		limit = q["limit"][0]
	}
	if len(q["since"]) > 0 {
		since = q["since"][0]
	}
	if len(q["desc"]) > 0 {
		desc = q["desc"][0]
	}

	threads, err := d.useCase.GetForumThreads(slug, limit, since, desc)
	if err != nil {
		if err == models.ErrForumNotFound {
			response.SendResponse(w, http.StatusNotFound, models.Error{Message: err.Error()})
			return
		} else {
			response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
			return
		}
	}
	if len(threads.Threads) == 0 {
		response.SendResponse(w, http.StatusOK, []models.Thread{})
		return
	}
	response.SendResponse(w, http.StatusOK, threads.Threads)
	return
}
