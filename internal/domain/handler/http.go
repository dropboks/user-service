package handler

import "github.com/gin-gonic/gin"

func RegisterUserRoutes(r *gin.Engine, uh UserHandler) *gin.Engine {
	{
		r.GET("/me", uh.GetProfile)
	}
	return r
}
