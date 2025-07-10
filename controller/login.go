package controller

import (
	"library/api/request"
	"library/api/response"
	"library/tool"
	"net/http"

	"github.com/gin-gonic/gin"
)

// todo:自动登录
type LoginController struct {
	ls tool.LoginService
}

func NewLoginController(ls tool.LoginService) *LoginController {
	return &LoginController{ls: ls}
}

func (lc *LoginController) Login(c *gin.Context) {
	var login request.Login
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "参数解析错误:" + err.Error(),
		})
		return
	}

	if err := lc.ls.LoginFirst(login); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "登录失败:" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Code:    200,
		Message: "登录成功",
	})
}
