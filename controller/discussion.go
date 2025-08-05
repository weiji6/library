package controller

import (
	"library/api/request"
	"library/api/response"
	"library/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type DiscussionController struct {
	ds service.Discussion
	rc *redis.Client
}

func NewDiscussionController(ds service.Discussion, rc *redis.Client) *DiscussionController {
	return &DiscussionController{
		ds: ds,
		rc: rc,
	}
}

func (dc *DiscussionController) GetDiscussion(c *gin.Context) {
	classID := c.Query("classID")
	date := c.Query("date")

	//classID := "103915682"
	//date := "20250708"

	discussion, err := dc.ds.GetDiscussion(classID, date)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "获取研讨间信息失败:" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Code:    200,
		Message: "研讨间信息获取成功",
		Data:    discussion,
	})
}

func (dc *DiscussionController) SearchUser(c *gin.Context) {
	StudentId := c.Query("studentID")

	Student, err := dc.ds.SearchUser(StudentId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "查找失败:" + err.Error(),
		})
		return
	}

	if Student.Name == "" {
		c.JSON(http.StatusNotFound, response.Response{
			Code:    404,
			Message: "查找失败:无该学生",
		})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Code:    200,
		Message: "查找成功",
		Data:    Student,
	})
}

func (dc *DiscussionController) ReserveDiscussion(c *gin.Context) {
	var Reserve request.ReserveDiscussion

	if err := c.ShouldBind(&Reserve); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "参数绑定失败:" + err.Error(),
		})
		return
	}

	result, err := dc.ds.ReserveDiscussion(Reserve)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "研讨间预约失败:" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Code:    200,
		Message: "研讨间预约成功",
		Data:    result,
	})
}
