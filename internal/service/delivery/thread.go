package delivery

import (
	"forum/forum/internal/models"
	"forum/forum/internal/service"
	log "forum/forum/pkg/logger"
	"forum/forum/pkg/response"
	"github.com/gorilla/mux"
	"net/http"
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
	message := threadLogMessage + "CreateThreadPosts:"
	log.Info(message + "started")
	vars := mux.Vars(r)
	slugOrId := vars["slug_or_id"]
	posts, err := response.GetPostsFromRequest(r.Body)
	if err != nil {
		log.Error(message+"error = ", err)
		response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
	}
	newPosts, err := d.useCase.CreateThreadPosts(slugOrId, posts)
	if err != nil {
		log.Error(message+"error = ", err)
		if err == models.ErrThreadNotFound {
			response.SendResponse(w, http.StatusNotFound, models.Error{Message: err.Error()})
		} else if err == models.ErrPostNotFound {
			response.SendResponse(w, http.StatusConflict, models.Error{Message: err.Error()})
		} else {
			response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
		}
	}
	response.SendResponse(w, http.StatusCreated, newPosts)
	log.Info(message + "ended")
}

func (d *ThreadDelivery) GetThreadDetails(w http.ResponseWriter, r *http.Request) {
	message := threadLogMessage + "GetThreadDetails:"
	log.Info(message + "started")
	vars := mux.Vars(r)
	slugOrId := vars["slug_or_id"]
	thread, err := d.useCase.GetThreadDetails(slugOrId)
	if err != nil {
		log.Error(message+"error = ", err)
		if err == models.ErrThreadNotFound {
			response.SendResponse(w, http.StatusNotFound, models.Error{Message: err.Error()})
		} else {
			response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
		}
	}
	response.SendResponse(w, http.StatusOK, thread)
	log.Info(message + "ended")
}

func (d *ThreadDelivery) UpdateThreadDetails(w http.ResponseWriter, r *http.Request) {
	message := threadLogMessage + "UpdateThreadDetails:"
	log.Info(message + "started")
	vars := mux.Vars(r)
	slugOrId := vars["slug_or_id"]
	thread, err := response.GetThreadFromRequest(r.Body)
	if err != nil {
		log.Error(message+"error = ", err)
		response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
	}
	updatedThread, err := d.useCase.UpdateThreadDetails(slugOrId, thread)
	if err != nil {
		log.Error(message+"error = ", err)
		if err == models.ErrThreadNotFound {
			response.SendResponse(w, http.StatusNotFound, models.Error{Message: err.Error()})
		} else {
			response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
		}
	}
	response.SendResponse(w, http.StatusOK, updatedThread)
	log.Info(message + "ended")
}

func (d *ThreadDelivery) GetThreadPosts(w http.ResponseWriter, r *http.Request) {
	message := threadLogMessage + "GetThreadPosts:"
	log.Info(message + "started")
	vars := mux.Vars(r)
	slugOrId := vars["slug_od_id"]

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
		log.Error(message+"error = ", err)
		if err == models.ErrThreadNotFound {
			response.SendResponse(w, http.StatusNotFound, models.Error{Message: err.Error()})
		} else if err == models.ErrParentPostNotFound {
			response.SendResponse(w, http.StatusConflict, models.Error{Message: err.Error()})
		} else {
			response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
		}
	}
	response.SendResponse(w, http.StatusOK, posts)
	log.Info(message + "ended")
}

func (d *ThreadDelivery) VoteForThread(w http.ResponseWriter, r *http.Request) {
	message := threadLogMessage + "VoteForThread:"
	log.Info(message + "started")
	vars := mux.Vars(r)
	slugOdId := vars["slugOdId"]
	vote, err := response.GetVoteFromRequest(r.Body)
	if err != nil {
		log.Error(message+"error = ", err)
		response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
	}
	thread, err := d.useCase.VoteForThread(slugOdId, vote)
	if err != nil {
		log.Error(message+"error = ", err)
		if err == models.ErrThreadNotFound {
			response.SendResponse(w, http.StatusNotFound, models.Error{Message: err.Error()})
		} else {
			response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
		}
	}
	response.SendResponse(w, http.StatusOK, thread)
	log.Info(message + "ended")
}
