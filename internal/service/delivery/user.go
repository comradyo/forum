package delivery

import (
	"forum/internal/models"
	"forum/internal/service"
	"forum/pkg/response"
	"net/http"

	routing "github.com/qiangxue/fasthttp-routing"
)

const userLogMessage = "delivery:user:"

type UserDelivery struct {
	useCase service.UserUseCaseInterface
}

func NewUserDelivery(useCase service.UserUseCaseInterface) *UserDelivery {
	return &UserDelivery{
		useCase: useCase,
	}
}

func (d *UserDelivery) CreateUser(ctx *routing.Context) error {
	nickname := ctx.Param("nickname")

	ctx.Request.Body()
	profile, err := response.GetProfileFromRequest(ctx.PostBody())
	if err != nil {
		response.SendResponse(ctx, http.StatusInternalServerError, models.Error{Message: err.Error()})
		return nil
	}
	profile.Nickname = nickname
	users, err := d.useCase.CreateUser(profile)
	if err != nil {
		if err == models.ErrUserExists {
			response.SendResponse(ctx, http.StatusConflict, users)
			return nil
		} else {
			response.SendResponse(ctx, http.StatusInternalServerError, models.Error{Message: err.Error()})
			return nil
		}
	}
	response.SendResponse(ctx, http.StatusCreated, users[0])
	return nil
}

func (d *UserDelivery) GetUserProfile(ctx *routing.Context) error {
	nickname := ctx.Param("nickname")
	profile, err := d.useCase.GetUserProfile(nickname)
	if err != nil {
		if err == models.ErrUserNotFound {
			response.SendResponse(ctx, http.StatusNotFound, models.Error{Message: err.Error()})
			return nil
		} else {
			response.SendResponse(ctx, http.StatusInternalServerError, models.Error{Message: err.Error()})
			return nil
		}
	}
	response.SendResponse(ctx, http.StatusOK, profile)
	return nil
}

func (d *UserDelivery) UpdateUserProfile(ctx *routing.Context) error {
	nickname := ctx.Param("nickname")
	profile, err := response.GetProfileFromRequest(ctx.Request.Body())
	if err != nil {
		response.SendResponse(ctx, http.StatusInternalServerError, models.Error{Message: err.Error()})
		return nil
	}
	profile.Nickname = nickname
	updatedProfile, err := d.useCase.UpdateUserProfile(profile)
	if err != nil {
		if err == models.ErrUserNotFound {
			response.SendResponse(ctx, http.StatusNotFound, models.Error{Message: err.Error()})
			return nil
		} else if err == models.ErrProfileUpdateConflict {
			response.SendResponse(ctx, http.StatusConflict, models.Error{Message: err.Error()})
			return nil
		} else {
			response.SendResponse(ctx, http.StatusInternalServerError, models.Error{Message: err.Error()})
			return nil
		}
	}
	response.SendResponse(ctx, http.StatusOK, updatedProfile)
	return nil
}
