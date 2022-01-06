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
	posts, err := utils.GetPostsFromRequest(r.Body)
	if err != nil {
		response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
	}
	newPosts, err := d.useCase.CreateThreadPosts(slugOrId, posts)
	if err != nil {
		if err == models.ErrThreadNotFound {
			response.SendResponse(w, http.StatusNotFound, models.Error{Message: err.Error()})
		} else if err == models.ErrPostNotFound {
			response.SendResponse(w, http.StatusConflict, models.Error{Message: err.Error()})
		} else {
			response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
		}
	}
	response.SendResponse(w, http.StatusCreated, newPosts)
}

func (d *ThreadDelivery) GetThreadDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slugOrId := vars["slug_or_id"]
	thread, err := d.useCase.GetThreadDetails(slugOrId)
	if err != nil {
		if err == models.ErrThreadNotFound {
			response.SendResponse(w, http.StatusNotFound, models.Error{Message: err.Error()})
		} else {
			response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
		}
	}
	response.SendResponse(w, http.StatusOK, thread)
}

func (d *ThreadDelivery) UpdateThreadDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slugOrId := vars["slug_or_id"]
	thread, err := utils.GetThreadFromRequest(r.Body)
	if err != nil {
		response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
	}
	updatedThread, err := d.useCase.UpdateThreadDetails(slugOrId, thread)
	if err != nil {
		if err == models.ErrThreadNotFound {
			response.SendResponse(w, http.StatusNotFound, models.Error{Message: err.Error()})
		} else {
			response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
		}
	}
	response.SendResponse(w, http.StatusOK, updatedThread)
}

func (d *ThreadDelivery) GetThreadPosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slugOrId := vars["slug_od_id"]
	q := r.URL.Query()
	var limit int32
	var since string
	var sort string
	var desc bool
	if len(q["limit"]) > 0 {
		limitInt, _ := strconv.Atoi(q["limit"][0])
		limit = int32(limitInt)
	}
	if len(q["since"]) > 0 {
		since = q["since"][0]
	}
	if len(q["sort"]) > 0 {
		sort = q["sort"][0]
	}
	if len(q["desc"]) > 0 {
		descStr := q["desc"][0]
		if descStr == "true" {
			desc = true
		} else {
			desc = false
		}
	}
	posts, err := d.useCase.GetThreadPosts(slugOrId, limit, since, sort, desc)
	if err != nil {
		if err == models.ErrThreadNotFound {
			response.SendResponse(w, http.StatusNotFound, models.Error{Message: err.Error()})
		} else {
			response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
		}
	}
	response.SendResponse(w, http.StatusOK, posts)
}

func (d *ThreadDelivery) VoteForThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slugOdId := vars["slugOdId"]
	vote, err := utils.GetVoteFromRequest(r.Body)
	if err != nil {
		response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
	}
	thread, err := d.useCase.VoteForThread(slugOdId, vote)
	if err != nil {
		if err == models.ErrThreadNotFound {
			response.SendResponse(w, http.StatusNotFound, models.Error{Message: err.Error()})
		} else {
			response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
		}
	}
	response.SendResponse(w, http.StatusOK, thread)
}
