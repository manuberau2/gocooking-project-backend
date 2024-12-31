package utils

import (
	"GoCooking/Backend/clients/responses"

	"github.com/gin-gonic/gin"
)

func SetUserInContext(c *gin.Context, user *responses.UsuarioInfo) {
	c.Set("UsuarioInfo", user)
}

func GetUserInfoFromContext(c *gin.Context) *responses.UsuarioInfo {
	userInfo, _ := c.Get("UsuarioInfo")

	user, _ := userInfo.(*responses.UsuarioInfo)

	return user
}
