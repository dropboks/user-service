package handler

import "github.com/gin-gonic/gin"

func RegisterUserRoutes(r *gin.Engine, uh UserHandler) *gin.RouterGroup {
	user := r.Group("/user", uh.GetProfile)
	{
		user.GET("/me", uh.GetProfile)
	}
	return user
}
