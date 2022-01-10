package delivery

import (
	"forum/forum/internal/models"
	"forum/forum/internal/service"
	log "forum/forum/pkg/logger"
	"forum/forum/pkg/response"
	"net/http"
)

const serviceLogMessage = "delivery:service:"

type ServiceDelivery struct {
	useCase service.ServiceUseCaseInterface
}

func NewServiceDelivery(useCase service.ServiceUseCaseInterface) *ServiceDelivery {
	return &ServiceDelivery{
		useCase: useCase,
	}
}

func (d *ServiceDelivery) Clear(w http.ResponseWriter, r *http.Request) {
	message := serviceLogMessage + "Clear:"
	log.Info(message + "started")
	err := d.useCase.Clear()
	if err != nil {
		log.Error(message+"error = ", err)
		response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
		return
	}
	response.SendResponse(w, http.StatusOK, nil)
	log.Info(message + "ended")
	return
}

func (d *ServiceDelivery) GetStatus(w http.ResponseWriter, r *http.Request) {
	message := serviceLogMessage + "GetStatus:"
	log.Info(message + "started")
	status, err := d.useCase.GetStatus()
	if err != nil {
		log.Error(message+"error = ", err)
		response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
		return
	}
	response.SendResponse(w, http.StatusOK, status)
	log.Info(message + "ended")
	return
}
