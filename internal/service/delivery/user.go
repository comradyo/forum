package delivery

import "forum/forum/internal/service"

const userLogMessage = "delivery:user:"

type UserDelivery struct {
	useCase service.UserUseCaseInterface
}

func NewUserDelivery(useCase service.UserUseCaseInterface) *UserDelivery {
	return &UserDelivery{
		useCase: useCase,
	}
}
