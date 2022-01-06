package delivery

import (
	"forum/forum/internal/models"
	"forum/forum/internal/service"
	"forum/forum/internal/utils"
	"forum/forum/pkg/response"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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
	forum, err := utils.GetForumFromRequest(r.Body)
	if err != nil {
		response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
	}
	newForum, err := d.useCase.CreateForum(forum)
	if err != nil {
		if err == models.ErrForumExists {
			response.SendResponse(w, http.StatusConflict, newForum)
		} else if err == models.ErrUserNotFound {
			response.SendResponse(w, http.StatusNotFound, models.Error{Message: err.Error()})
		} else {
			response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
		}
	}
	response.SendResponse(w, http.StatusCreated, newForum)
}

func (d *ForumDelivery) GetForumDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]
	forum, err := d.useCase.GetForumDetails(slug)
	if err != nil {
		if err == models.ErrForumNotFound {
			response.SendResponse(w, http.StatusNotFound, models.Error{Message: err.Error()})
		} else {
			response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
		}
	}
	response.SendResponse(w, http.StatusOK, forum)
}

func (d *ForumDelivery) CreateForumThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]
	thread, err := utils.GetThreadFromRequest(r.Body)
	if err != nil {
		response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
	}
	newThread, err := d.useCase.CreateForumThread(slug, thread)
	if err != nil {
		if err == models.ErrThreadExists {
			response.SendResponse(w, http.StatusConflict, newThread)
		} else if err == models.ErrUserNotFound || err == models.ErrForumNotFound {
			response.SendResponse(w, http.StatusNotFound, models.Error{Message: err.Error()})
		} else {
			response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
		}
	}
	response.SendResponse(w, http.StatusCreated, newThread)
}

func (d *ForumDelivery) GetForumUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	q := r.URL.Query()
	var limit int32
	var since string
	var desc bool

	if len(q["limit"]) > 0 {
		limitInt, _ := strconv.Atoi(q["limit"][0])
		limit = int32(limitInt)
	}
	if len(q["since"]) > 0 {
		since = q["since"][0]
	}
	if len(q["desc"]) > 0 {
		descStr := q["desc"][0]
		if descStr == "true" {
			desc = true
		} else {
			desc = false
		}
	}

	users, err := d.useCase.GetForumUsers(slug, limit, since, desc)
	if err != nil {
		if err == models.ErrForumNotFound {
			response.SendResponse(w, http.StatusNotFound, models.Error{Message: err.Error()})
		} else {
			response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
		}
	}
	response.SendResponse(w, http.StatusOK, users)
}

func (d *ForumDelivery) GetForumThreads(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	q := r.URL.Query()
	var limit int32
	var since string
	var desc bool

	if len(q["limit"]) > 0 {
		limitInt, _ := strconv.Atoi(q["limit"][0])
		limit = int32(limitInt)
	}
	if len(q["since"]) > 0 {
		since = q["since"][0]
	}
	if len(q["desc"]) > 0 {
		descStr := q["desc"][0]
		if descStr == "true" {
			desc = true
		} else {
			desc = false
		}
	}

	threads, err := d.useCase.GetForumThreads(slug, limit, since, desc)
	if err != nil {
		if err == models.ErrForumNotFound {
			response.SendResponse(w, http.StatusNotFound, models.Error{Message: err.Error()})
		} else {
			response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
		}
	}
	response.SendResponse(w, http.StatusOK, threads)
}
