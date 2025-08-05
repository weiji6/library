package router

import (
	"library/controller"

	"github.com/gin-gonic/gin"
)

type HistoryRecordRouter struct {
	hc *controller.HistoryRecordController
}

func NewHistoryRecordRouter(hc *controller.HistoryRecordController) *HistoryRecordRouter {
	return &HistoryRecordRouter{
		hc: hc,
	}
}

func (hr *HistoryRecordRouter) HistoryRecordRoute(r *gin.Engine) {
	historyRecordGroup := r.Group("/history_record")
	{
		historyRecordGroup.GET("/get", hr.hc.GetHistoryRecord)
	}
}
