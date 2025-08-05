package controller

import (
	"context"
	"encoding/json"
	"library/api/request"
	"library/api/response"
	"library/model"
	"library/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type SeatController struct {
	ss service.SeatService
	rc *redis.Client
}

func NewSeatController(ss service.SeatService, rc *redis.Client) *SeatController {
	return &SeatController{
		ss: ss,
		rc: rc,
	}
}

func (sc *SeatController) FetchSeat(c *gin.Context) {
	roomID := c.Query("room_id")
	if roomID == "" {
		// 如果没有指定roomID，返回所有座位信息
		sc.FetchAllSeats(c)
		return
	}

	// 从Redis获取数据
	ctx := context.Background()
	seatsJSON, err := sc.rc.Get(ctx, "all_seats").Bytes()
	if err == nil {
		var allSeats map[string][]model.Seat
		if err = json.Unmarshal(seatsJSON, &allSeats); err == nil {
			if seats, ok := allSeats[roomID]; ok {
				c.JSON(http.StatusOK, response.Response{
					Code:    200,
					Message: "座位状态获取成功",
					Data:    seats,
				})
				return
			}
		}
	}

	// 如果Redis中没有数据，直接从服务获取
	seat, err := sc.ss.FetchSeat(roomID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "获取座位信息失败",
		})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Code:    200,
		Message: "座位状态获取成功",
		Data:    seat,
	})
}

func (sc *SeatController) FetchAllSeats(c *gin.Context) {
	ctx := context.Background()
	seatsJSON, err := sc.rc.Get(ctx, "all_seats").Bytes()
	if err == nil {
		var allSeats map[string][]model.Seat
		if err = json.Unmarshal(seatsJSON, &allSeats); err == nil {
			c.JSON(http.StatusOK, response.Response{
				Code:    200,
				Message: "所有座位状态获取成功",
				Data:    allSeats,
			})
			return
		}
	}

	// 如果Redis中没有数据，直接从服务获取
	seats, err := sc.ss.FetchAllSeats()
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "获取所有座位信息失败:" + err.Error(),
		})
		return
	}

	seatsJSON, err = json.Marshal(seats)
	if err == nil {
		sc.rc.Set(ctx, "all_seats", seatsJSON, 10*time.Minute)
	}

	c.JSON(http.StatusOK, response.Response{
		Code:    200,
		Message: "所有座位状态获取成功",
		Data:    seats,
	})
}

func (sc *SeatController) ReserveSeat(c *gin.Context) {
	var reserve request.Reserve

	if err := c.ShouldBindJSON(&reserve); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "参数绑定失败:" + err.Error(),
		})
		return
	}

	result, err := sc.ss.ReserveSeat(reserve)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "座位预约失败:" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Code:    200,
		Message: "座位预约成功",
		Data:    result,
	})
}

func (sc *SeatController) SSEReserveSeat(c *gin.Context) {
	DevId := c.Query("dev_id")
	Start := c.Query("start")
	End := c.Query("end")

	if DevId == "" || Start == "" || End == "" {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "缺少必要的参数",
		})
		return
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Flush()

	now := time.Now()

	next := time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, now.Location())
	if now.After(next) {
		next = next.Add(24 * time.Hour)
	}

	waitDuration := time.Until(next)
	time.Sleep(waitDuration)

	reserve := request.Reserve{
		DevID: DevId,
		Start: Start,
		End:   End,
	}

	result, err := sc.ss.ReserveSeat(reserve)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "座位预约失败" + err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, response.Response{
			Code:    200,
			Message: "预约成功",
			Data:    result,
		})
	}
}

func (sc *SeatController) GetRecord(c *gin.Context) {
	result, err := sc.ss.GetRecord()
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "获取未来预约失败:" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Code:    200,
		Message: "获取未来预约成功",
		Data:    result,
	})
}

func (sc *SeatController) CancelReserve(c *gin.Context) {
	id := c.Param("id")

	result, err := sc.ss.CancelReserve(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "座位取消失败:" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Code:    200,
		Message: "取消预约成功",
		Data:    result,
	})
}
