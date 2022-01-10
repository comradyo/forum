package delivery

import (
	"forum/forum/internal/models"
	"forum/forum/internal/service"
	log "forum/forum/pkg/logger"
	"forum/forum/pkg/response"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

const postLogMessage = "delivery:post:"

type PostDelivery struct {
	useCase service.PostUseCaseInterface
}

func NewPostDelivery(useCase service.PostUseCaseInterface) *PostDelivery {
	return &PostDelivery{
		useCase: useCase,
	}
}

func (d *PostDelivery) GetPostDetails(w http.ResponseWriter, r *http.Request) {
	message := postLogMessage + "GetPostDetails:"
	log.Info(message + "started")
	vars := mux.Vars(r)
	id := vars["id"]
	idInt, _ := strconv.Atoi(id)
	idInt64 := int64(idInt)

	q := r.URL.Query()
	var related string
	if len(q["related"]) > 0 {
		related = q["related"][0]
	}

	postFull, err := d.useCase.GetPostDetails(idInt64, related)
	if err != nil {
		log.Error(message+"error = ", err)
		if err == models.ErrPostNotFound {
			response.SendResponse(w, http.StatusNotFound, models.Error{Message: err.Error()})
			return
		} else {
			response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
			return
		}
	}
	//TODO: Возможно, здесь будет другая структура
	response.SendResponse(w, http.StatusOK, postFull)
	log.Info(message + "ended")
	return
}

func (d *PostDelivery) UpdatePostDetails(w http.ResponseWriter, r *http.Request) {
	message := postLogMessage + "UpdatePostDetails:"
	log.Info(message + "started")
	vars := mux.Vars(r)
	id := vars["id"]
	post, err := response.GetPostFromRequest(r.Body)
	if err != nil {
		log.Error(message+"error = ", err)
		response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
		return
	}
	idInt, _ := strconv.Atoi(id)
	post.Id = int64(idInt)
	updatedPost, err := d.useCase.UpdatePostDetails(post)
	if err != nil {
		log.Error(message+"error = ", err)
		if err == models.ErrPostNotFound {
			response.SendResponse(w, http.StatusNotFound, models.Error{Message: err.Error()})
			return
		} else {
			response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
			return
		}
	}
	response.SendResponse(w, http.StatusOK, updatedPost)
	log.Info(message + "ended")
	return
}
