package controller

import (
	"library/api/response"
	"library/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreditController struct {
	CreditPointsService service.CreditPoints
}

func NewCreditController(creditPointsService service.CreditPoints) *CreditController {
	return &CreditController{
		CreditPointsService: creditPointsService,
	}
}

func (cc *CreditController) GetCreditPoint(c *gin.Context) {
	CreditPoints, err := cc.CreditPointsService.GetCreditPoints()
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "个人信用分获取失败:" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Code:    200,
		Message: "个人信用分获取成功",
		Data:    CreditPoints,
	})
}
