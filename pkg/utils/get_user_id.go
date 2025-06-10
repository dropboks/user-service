package utils

import (
	"encoding/json"

	"github.com/dropboks/user-service/internal/domain/dto"
	"github.com/gin-gonic/gin"
)

func GetUserId(c *gin.Context) string {
	userDataHeader := c.Request.Header.Get("User-Data")
	if userDataHeader != "" {
		var ud dto.UserData
		err := json.Unmarshal([]byte(userDataHeader), &ud)
		if err == nil {
			return ud.UserId
		}
	}
	return ""
}
