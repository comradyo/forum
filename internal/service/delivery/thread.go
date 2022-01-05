package delivery

import "forum/forum/internal/service"

const threadLogMessage = "delivery:thread:"

type ThreadDelivery struct {
	useCase service.ThreadUseCaseInterface
}

func NewThreadDelivery(useCase service.ThreadUseCaseInterface) *ThreadDelivery {
	return &ThreadDelivery{
		useCase: useCase,
	}
}
