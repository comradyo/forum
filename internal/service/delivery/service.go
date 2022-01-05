package delivery

import "forum/forum/internal/service"

const serviceLogMessage = "delivery:service:"

type ServiceDelivery struct {
	useCase service.ServiceUseCaseInterface
}

func NewServiceDelivery(useCase service.ServiceUseCaseInterface) *ServiceDelivery {
	return &ServiceDelivery{
		useCase: useCase,
	}
}
