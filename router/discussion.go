package router

import (
	"library/controller"

	"github.com/gin-gonic/gin"
)

type DiscussionRouter struct {
	dc *controller.DiscussionController
}

func NewDiscussionRouter(dc *controller.DiscussionController) *DiscussionRouter {
	return &DiscussionRouter{
		dc: dc,
	}
}

func (dr *DiscussionRouter) DiscussionRouter(r *gin.Engine) {
	discussionGroup := r.Group("/discussion")
	{
		discussionGroup.GET("/get", dr.dc.GetDiscussion)
		discussionGroup.POST("/reserve", dr.dc.ReserveDiscussion)
		discussionGroup.GET("/cancel/:id", dr.dc.CancelDiscussion)

		discussionGroup.GET("/user/search", dr.dc.SearchUser)
	}
}
