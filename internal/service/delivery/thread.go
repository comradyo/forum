package delivery

import (
	"forum/internal/models"
	"forum/internal/service"

	"forum/pkg/response"
	"net/http"

	routing "github.com/qiangxue/fasthttp-routing"
)

const threadLogMessage = "delivery:thread:"

type ThreadDelivery struct {
	useCase service.ThreadUseCaseInterface
}

func NewThreadDelivery(useCase service.ThreadUseCaseInterface) *ThreadDelivery {
	return &ThreadDelivery{
		useCase: useCase,
	}
}

func (d *ThreadDelivery) CreateThreadPosts(ctx *routing.Context) error {
	slugOrId := ctx.Param("slug_or_id")
	posts, err := response.GetPostsFromRequest(ctx.Request.Body())
	if err != nil {
		response.SendResponse(ctx, http.StatusInternalServerError, models.Error{Message: err.Error()})
		return nil
	}
	newPosts, err := d.useCase.CreateThreadPosts(slugOrId, &models.Posts{Posts: posts})
	if err != nil {
		if err == models.ErrThreadNotFound || err == models.ErrUserNotFound {
			response.SendResponse(ctx, http.StatusNotFound, models.Error{Message: err.Error()})
			return nil
		} else if err == models.ErrPostNotFound {
			response.SendResponse(ctx, http.StatusConflict, models.Error{Message: err.Error()})
			return nil
		} else {
			response.SendResponse(ctx, http.StatusInternalServerError, models.Error{Message: err.Error()})
			return nil
		}
	}
	if len(newPosts.Posts) == 0 {
		response.SendResponse(ctx, http.StatusCreated, []models.Post{})
		return nil
	}
	response.SendResponse(ctx, http.StatusCreated, newPosts.Posts)
	return nil
}

func (d *ThreadDelivery) GetThreadDetails(ctx *routing.Context) error {
	slugOrId := ctx.Param("slug_or_id")
	thread, err := d.useCase.GetThreadDetails(slugOrId)
	if err != nil {
		if err == models.ErrThreadNotFound {
			response.SendResponse(ctx, http.StatusNotFound, models.Error{Message: err.Error()})
			return nil
		} else {
			response.SendResponse(ctx, http.StatusInternalServerError, models.Error{Message: err.Error()})
			return nil
		}
	}
	response.SendResponse(ctx, http.StatusOK, thread)
	return nil
}

func (d *ThreadDelivery) UpdateThreadDetails(ctx *routing.Context) error {
	slugOrId := ctx.Param("slug_or_id")
	thread, err := response.GetThreadFromRequest(ctx.Request.Body())
	if err != nil {
		response.SendResponse(ctx, http.StatusInternalServerError, models.Error{Message: err.Error()})
		return nil
	}
	updatedThread, err := d.useCase.UpdateThreadDetails(slugOrId, thread)
	if err != nil {
		if err == models.ErrThreadNotFound {
			response.SendResponse(ctx, http.StatusNotFound, models.Error{Message: err.Error()})
			return nil
		} else {
			response.SendResponse(ctx, http.StatusInternalServerError, models.Error{Message: err.Error()})
			return nil
		}
	}
	response.SendResponse(ctx, http.StatusOK, updatedThread)
	return nil
}

func (d *ThreadDelivery) GetThreadPosts(ctx *routing.Context) error {
	slugOrId := ctx.Param("slug_or_id")

	limit := string(ctx.QueryArgs().Peek("limit"))
	since := string(ctx.QueryArgs().Peek("since"))
	sort := string(ctx.QueryArgs().Peek("sort"))
	desc := string(ctx.QueryArgs().Peek("desc"))

	posts, err := d.useCase.GetThreadPosts(slugOrId, limit, since, sort, desc)
	if err != nil {
		if err == models.ErrThreadNotFound {
			response.SendResponse(ctx, http.StatusNotFound, models.Error{Message: err.Error()})
			return nil
		} else {
			response.SendResponse(ctx, http.StatusInternalServerError, models.Error{Message: err.Error()})
			return nil
		}
	}

	if len(posts.Posts) == 0 {
		response.SendResponse(ctx, http.StatusOK, []models.Post{})
		return nil
	}
	response.SendResponse(ctx, http.StatusOK, posts.Posts)
	return nil
}

func (d *ThreadDelivery) VoteForThread(ctx *routing.Context) error {
	slugOrId := ctx.Param("slug_or_id")
	vote, err := response.GetVoteFromRequest(ctx.Request.Body())
	if err != nil {
		response.SendResponse(ctx, http.StatusInternalServerError, models.Error{Message: err.Error()})
		return nil
	}
	thread, err := d.useCase.VoteForThread(slugOrId, vote)
	if err != nil {
		if err == models.ErrThreadNotFound || err == models.ErrUserNotFound {
			response.SendResponse(ctx, http.StatusNotFound, models.Error{Message: err.Error()})
			return nil
		} else {
			response.SendResponse(ctx, http.StatusInternalServerError, models.Error{Message: err.Error()})
			return nil
		}
	}
	response.SendResponse(ctx, http.StatusOK, thread)
	return nil
}
