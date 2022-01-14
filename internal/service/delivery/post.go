package delivery

import (
	"forum/internal/models"
	"forum/internal/service"

	"forum/pkg/response"
	routing "github.com/qiangxue/fasthttp-routing"
	"net/http"
	"strconv"
)

const postLogMessage = "delivery:post:"

type PostDelivery struct {
	useCase service.PostUseCaseInterface
}

func NewPostDelivery(useCase service.PostUseCaseInterface) *PostDelivery {
	return &PostDelivery{
		useCase: useCase,
	}
}

func (d *PostDelivery) GetPostDetails(ctx *routing.Context) error {
	id := ctx.Param("id")
	idInt, _ := strconv.Atoi(id)
	idInt64 := int64(idInt)

	related := string(ctx.QueryArgs().Peek("related"))

	postFull, err := d.useCase.GetPostDetails(idInt64, related)
	if err != nil {
		if err == models.ErrPostNotFound {
			response.SendResponse(ctx, http.StatusNotFound, models.Error{Message: err.Error()})
			return nil
		} else {
			response.SendResponse(ctx, http.StatusInternalServerError, models.Error{Message: err.Error()})
			return nil
		}
	}
	//TODO: Возможно, здесь будет другая структура
	response.SendResponse(ctx, http.StatusOK, postFull)
	return nil
}

func (d *PostDelivery) UpdatePostDetails(ctx *routing.Context) error {
	id := ctx.Param("id")
	post, err := response.GetPostFromRequest(ctx.Request.Body())
	if err != nil {
		response.SendResponse(ctx, http.StatusInternalServerError, models.Error{Message: err.Error()})
		return nil
	}
	idInt, _ := strconv.Atoi(id)
	post.Id = int64(idInt)
	updatedPost, err := d.useCase.UpdatePostDetails(post)
	if err != nil {
		if err == models.ErrPostNotFound {
			response.SendResponse(ctx, http.StatusNotFound, models.Error{Message: err.Error()})
			return nil
		} else {
			response.SendResponse(ctx, http.StatusInternalServerError, models.Error{Message: err.Error()})
			return nil
		}
	}
	response.SendResponse(ctx, http.StatusOK, updatedPost)
	return nil
}
