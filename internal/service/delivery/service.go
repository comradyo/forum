package delivery

import (
	"forum/internal/models"
	"forum/internal/service"

	"forum/pkg/response"
	"net/http"

	routing "github.com/qiangxue/fasthttp-routing"
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

func (d *ServiceDelivery) Clear(ctx *routing.Context) error {
	err := d.useCase.Clear()
	if err != nil {
		response.SendResponse(ctx, http.StatusInternalServerError, models.Error{Message: err.Error()})
		return nil
	}
	response.SendResponse(ctx, http.StatusOK, nil)
	return nil
}

func (d *ServiceDelivery) GetStatus(ctx *routing.Context) error {
	status, err := d.useCase.GetStatus()
	if err != nil {
		response.SendResponse(ctx, http.StatusInternalServerError, models.Error{Message: err.Error()})
		return nil
	}
	response.SendResponse(ctx, http.StatusOK, status)
	return nil
}
