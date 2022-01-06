package delivery

import (
	"forum/forum/internal/models"
	"forum/forum/internal/service"
	"forum/forum/internal/utils"
	"forum/forum/pkg/response"
	"github.com/gorilla/mux"
	"net/http"
)

const userLogMessage = "delivery:user:"

type UserDelivery struct {
	useCase service.UserUseCaseInterface
}

func NewUserDelivery(useCase service.UserUseCaseInterface) *UserDelivery {
	return &UserDelivery{
		useCase: useCase,
	}
}

func (d *UserDelivery) CreateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nickname := vars["nickname"]
	profile, err := utils.GetProfileFromRequest(r.Body)
	if err != nil {
		response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
	}
	profile.Nickname = nickname
	newUser, err := d.useCase.CreateUser(profile)
	if err != nil {
		if err == models.ErrUserExists {
			response.SendResponse(w, http.StatusConflict, newUser)
		} else {
			response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
		}
	}
	response.SendResponse(w, http.StatusCreated, newUser)
}

func (d *UserDelivery) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nickname := vars["nickname"]
	profile, err := d.useCase.GetUserProfile(nickname)
	if err != nil {
		if err == models.ErrUserNotFound {
			response.SendResponse(w, http.StatusNotFound, models.Error{Message: err.Error()})
		} else {
			response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
		}
	}
	response.SendResponse(w, http.StatusOK, profile)
}

func (d *UserDelivery) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nickname := vars["nickname"]
	profile, err := utils.GetProfileFromRequest(r.Body)
	if err != nil {
		response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
	}
	profile.Nickname = nickname
	updatedProfile, err := d.useCase.UpdateUserProfile(profile)
	if err != nil {
		if err == models.ErrUserNotFound {
			response.SendResponse(w, http.StatusNotFound, models.Error{Message: err.Error()})
		} else if err == models.ErrProfileUpdateConflict {
			response.SendResponse(w, http.StatusConflict, models.Error{Message: err.Error()})
		} else {
			response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
		}
	}
	response.SendResponse(w, http.StatusOK, updatedProfile)
}
