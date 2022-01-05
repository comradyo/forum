package delivery

import "forum/forum/internal/service"

const postLogMessage = "delivery:post:"

type PostDelivery struct {
	useCase service.PostUseCaseInterface
}

func NewPostDelivery(useCase service.PostUseCaseInterface) *PostDelivery {
	return &PostDelivery{
		useCase: useCase,
	}
}
