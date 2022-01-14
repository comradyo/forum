package delivery

import (
	"forum/internal/models"
	"forum/internal/service"

	"forum/pkg/response"
	"net/http"

	routing "github.com/qiangxue/fasthttp-routing"
)

const forumLogMessage = "delivery:forum:"

type ForumDelivery struct {
	useCase service.ForumUseCaseInterface
}

func NewForumDelivery(useCase service.ForumUseCaseInterface) *ForumDelivery {
	return &ForumDelivery{
		useCase: useCase,
	}
}

//ctx *routing.Context
func (d *ForumDelivery) CreateForum(ctx *routing.Context) error {
	forum, err := response.GetForumFromRequest(ctx.Request.Body())
	if err != nil {
		response.SendResponse(ctx, http.StatusInternalServerError, models.Error{Message: err.Error()})
		return nil
	}
	newForum, err := d.useCase.CreateForum(forum)
	if err != nil {
		if err == models.ErrForumExists {
			response.SendResponse(ctx, http.StatusConflict, newForum)
			return nil
		} else if err == models.ErrUserNotFound {
			response.SendResponse(ctx, http.StatusNotFound, models.Error{Message: err.Error()})
			return nil
		} else {
			response.SendResponse(ctx, http.StatusInternalServerError, models.Error{Message: err.Error()})
			return nil
		}
	}
	response.SendResponse(ctx, http.StatusCreated, newForum)
	return nil
}

func (d *ForumDelivery) GetForumDetails(ctx *routing.Context) error {
	slug := ctx.Param("slug")
	forum, err := d.useCase.GetForumDetails(slug)
	if err != nil {
		if err == models.ErrForumNotFound {
			response.SendResponse(ctx, http.StatusNotFound, models.Error{Message: err.Error()})
			return nil
		} else {
			response.SendResponse(ctx, http.StatusInternalServerError, models.Error{Message: err.Error()})
			return nil
		}
	}
	response.SendResponse(ctx, http.StatusOK, forum)
	return nil
}

func (d *ForumDelivery) CreateForumThread(ctx *routing.Context) error {
	slug := ctx.Param("slug")
	thread, err := response.GetThreadFromRequest(ctx.Request.Body())
	if err != nil {
		response.SendResponse(ctx, http.StatusInternalServerError, models.Error{Message: err.Error()})
		return nil
	}
	newThread, err := d.useCase.CreateForumThread(slug, thread)
	if err != nil {
		if err == models.ErrThreadExists {
			response.SendResponse(ctx, http.StatusConflict, newThread)
			return nil
		} else if err == models.ErrUserNotFound || err == models.ErrForumNotFound {
			response.SendResponse(ctx, http.StatusNotFound, models.Error{Message: err.Error()})
			return nil
		} else {
			response.SendResponse(ctx, http.StatusInternalServerError, models.Error{Message: err.Error()})
			return nil
		}
	}
	response.SendResponse(ctx, http.StatusCreated, newThread)
	return nil
}

func (d *ForumDelivery) GetForumUsers(ctx *routing.Context) error {
	slug := ctx.Param("slug")

	limit := string(ctx.QueryArgs().Peek("limit"))
	since := string(ctx.QueryArgs().Peek("since"))
	desc := string(ctx.QueryArgs().Peek("desc"))

	users, err := d.useCase.GetForumUsers(slug, limit, since, desc)
	if err != nil {
		if err == models.ErrForumNotFound {
			response.SendResponse(ctx, http.StatusNotFound, models.Error{Message: err.Error()})
			return nil
		} else {
			response.SendResponse(ctx, http.StatusInternalServerError, models.Error{Message: err.Error()})
			return nil
		}
	}
	if len(users.Users) == 0 {
		response.SendResponse(ctx, http.StatusOK, []models.User{})
		return nil
	}
	response.SendResponse(ctx, http.StatusOK, users.Users)
	return nil
}

func (d *ForumDelivery) GetForumThreads(ctx *routing.Context) error {
	slug := ctx.Param("slug")

	limit := string(ctx.QueryArgs().Peek("limit"))
	since := string(ctx.QueryArgs().Peek("since"))
	desc := string(ctx.QueryArgs().Peek("desc"))

	threads, err := d.useCase.GetForumThreads(slug, limit, since, desc)
	if err != nil {
		if err == models.ErrForumNotFound {
			response.SendResponse(ctx, http.StatusNotFound, models.Error{Message: err.Error()})
			return nil
		} else {
			response.SendResponse(ctx, http.StatusInternalServerError, models.Error{Message: err.Error()})
			return nil
		}
	}
	if len(threads.Threads) == 0 {
		response.SendResponse(ctx, http.StatusOK, []models.Thread{})
		return nil
	}
	response.SendResponse(ctx, http.StatusOK, threads.Threads)
	return nil
}
