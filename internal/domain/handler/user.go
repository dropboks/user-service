package handler

import (
	"net/http"

	_utils "github.com/dropboks/sharedlib/utils"
	"github.com/dropboks/user-service/internal/domain/dto"
	"github.com/dropboks/user-service/internal/domain/service"
	"github.com/dropboks/user-service/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type (
	UserHandler interface {
		GetProfile(ctx *gin.Context)
	}
	userHandler struct {
		userService service.UserService
		logger      zerolog.Logger
	}
)

func NewUserHandler(userService service.UserService, logger zerolog.Logger) UserHandler {
	return &userHandler{
		userService: userService,
		logger:      logger,
	}
}

func (u *userHandler) GetProfile(ctx *gin.Context) {
	userId := utils.GetUserId(ctx)
	if userId == "" {
		u.logger.Error().Msg("unathorized")
		res := _utils.ReturnResponseError(401, dto.Err_UNAUTHORIZED_USER_ID_NOTFOUND.Error())
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
		return
	}
	user, err := u.userService.GetProfile(userId)
	if err != nil {
		if err == dto.Err_INTERNAL_FAILED_BUILD_QUERY || err == dto.Err_INTERNAL_FAILED_SCAN_USER {
			res := _utils.ReturnResponseError(500, err.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		} else if err == dto.Err_NOTFOUND_USER_NOT_FOUND {
			res := _utils.ReturnResponseError(404, err.Error())
			ctx.AbortWithStatusJSON(http.StatusNotFound, res)
			return
		}
	}
	res := _utils.ReturnResponseSuccess(200, "success get profile data", user)
	ctx.JSON(http.StatusOK, res)
}
