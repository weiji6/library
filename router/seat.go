package router

import (
	"library/controller"

	"github.com/gin-gonic/gin"
)

type SeatRouter struct {
	sc *controller.SeatController
}

func NewSeatRouter(sc *controller.SeatController) *SeatRouter {
	return &SeatRouter{
		sc: sc,
	}
}

func (sr *SeatRouter) SeatRouter(r *gin.Engine) {
	seatGroup := r.Group("/library")
	{
		// 获取单个房间座位信息
		seatGroup.GET("/seat/fetch", sr.sc.FetchSeat)

		// 获取所有房间座位信息
		seatGroup.GET("/seat/fetch/all", sr.sc.FetchAllSeats)

		seatGroup.POST("/seat/reserve", sr.sc.ReserveSeat)
		seatGroup.GET("/seat/reserve/SSE", sr.sc.SSEReserveSeat)
		seatGroup.GET("/seat/record", sr.sc.GetRecord)
		seatGroup.GET("/seat/cancel/:id", sr.sc.CancelSeat)
	}
}
