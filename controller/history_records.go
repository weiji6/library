package controller

import (
	"library/api/response"
	"library/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HistoryRecordController struct {
	hs service.HistoryRecords
}

func NewHistoryRecordController(hs service.HistoryRecords) *HistoryRecordController {
	return &HistoryRecordController{
		hs: hs,
	}
}

func (hc *HistoryRecordController) GetHistoryRecord(c *gin.Context) {
	record, err := hc.hs.GetHistoryRecords()
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "获取预约记录失败:" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Code:    200,
		Message: "预约记录获取成功",
		Data:    record,
	})
}
