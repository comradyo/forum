package delivery

import "forum/forum/internal/service"

const forumLogMessage = "delivery:forum:"

type ForumDelivery struct {
	useCase service.ForumUseCaseInterface
}

func NewForumDelivery(useCase service.ForumUseCaseInterface) *ForumDelivery {
	return &ForumDelivery{
		useCase: useCase,
	}
}
