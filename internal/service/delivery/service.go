package delivery

import (
	"forum/forum/internal/models"
	"forum/forum/internal/service"

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
	err := d.useCase.Clear()
	if err != nil {
		response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
		return
	}
	response.SendResponse(w, http.StatusOK, nil)
	return
}

func (d *ServiceDelivery) GetStatus(w http.ResponseWriter, r *http.Request) {
	status, err := d.useCase.GetStatus()
	if err != nil {
		response.SendResponse(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
		return
	}
	response.SendResponse(w, http.StatusOK, status)
	return
}
