package delivery

import (
	"forum/forum/internal/models"
	"forum/forum/internal/service"
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
	vars := mux.Vars(r)
	id := vars["id"]
	q := r.URL.Query()
	var related string
	if len(q["related"]) > 0 {
		related = q["related"][0]
	}
	postFull, err := d.useCase.GetPostDetails(id, related)
	if err != nil {
		if err == models.ErrPostNotFound {
			response.SendResponse(w, http.StatusNotFound, models.Error{Message: err.Error()})
		} else {
			response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
		}
	}
	//TODO: Возможно, здесь будет другая структура
	response.SendResponse(w, http.StatusOK, postFull)
}

func (d *PostDelivery) UpdatePostDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	post, err := response.GetPostFromRequest(r.Body)
	if err != nil {
		response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
	}
	idInt, _ := strconv.Atoi(id)
	post.Id = int64(idInt)
	updatedPost, err := d.useCase.UpdatePostDetails(post)
	if err != nil {
		if err == models.ErrPostNotFound {
			response.SendResponse(w, http.StatusNotFound, models.Error{Message: err.Error()})
		} else {
			response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
		}
	}
	response.SendResponse(w, http.StatusOK, updatedPost)
}
