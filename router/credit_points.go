package router

import (
	"library/controller"

	"github.com/gin-gonic/gin"
)

type CreditRouter struct {
	CreditController *controller.CreditController
}

func NewCreditRouter(cc *controller.CreditController) *CreditRouter {
	return &CreditRouter{
		CreditController: cc,
	}
}

func (cr *CreditRouter) CreditRouter(r *gin.Engine) {
	creditGroup := r.Group("/credit")
	{
		creditGroup.GET("/get", cr.CreditController.GetCreditPoint)
	}

}
