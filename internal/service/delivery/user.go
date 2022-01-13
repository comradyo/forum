package delivery

import (
	"forum/forum/internal/models"
	"forum/forum/internal/service"
	log "forum/forum/pkg/logger"
	"forum/forum/pkg/response"
	"net/http"

	"github.com/gorilla/mux"
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
	message := userLogMessage + "CreateUser:"
	log.Info(message + "started")
	vars := mux.Vars(r)
	nickname := vars["nickname"]
	profile, err := response.GetProfileFromRequest(r.Body)
	if err != nil {
		log.Error(message+"error = ", err)
		response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
		return
	}
	profile.Nickname = nickname
	users, err := d.useCase.CreateUser(profile)
	if err != nil {
		log.Error(message+"error = ", err)
		if err == models.ErrUserExists {
			response.SendResponse(w, http.StatusConflict, users)
			return
		} else {
			response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
			return
		}
	}
	response.SendResponse(w, http.StatusCreated, users[0])
	log.Info(message + "ended")
	return
}

func (d *UserDelivery) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	message := userLogMessage + "GetUserProfile:"
	log.Info(message + "started")
	vars := mux.Vars(r)
	nickname := vars["nickname"]
	profile, err := d.useCase.GetUserProfile(nickname)
	if err != nil {
		log.Error(message+"error = ", err)
		if err == models.ErrUserNotFound {
			response.SendResponse(w, http.StatusNotFound, models.Error{Message: err.Error()})
			return
		} else {
			response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
			return
		}
	}
	response.SendResponse(w, http.StatusOK, profile)
	log.Info(message + "ended")
	return
}

func (d *UserDelivery) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	message := userLogMessage + "UpdateUserProfile:"
	log.Info(message + "started")
	vars := mux.Vars(r)
	nickname := vars["nickname"]
	profile, err := response.GetProfileFromRequest(r.Body)
	if err != nil {
		log.Error(message+"error = ", err)
		response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
		return
	}
	profile.Nickname = nickname
	updatedProfile, err := d.useCase.UpdateUserProfile(profile)
	if err != nil {
		log.Error(message+"error = ", err)
		if err == models.ErrUserNotFound {
			response.SendResponse(w, http.StatusNotFound, models.Error{Message: err.Error()})
			return
		} else if err == models.ErrProfileUpdateConflict {
			response.SendResponse(w, http.StatusConflict, models.Error{Message: err.Error()})
			return
		} else {
			response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
			return
		}
	}
	response.SendResponse(w, http.StatusOK, updatedProfile)
	log.Info(message + "ended")
	return
}
