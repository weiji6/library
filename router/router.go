package router

import "github.com/gin-gonic/gin"

type App struct {
	r  *gin.Engine
	sr *SeatRouter
	cr *CreditRouter
	dr *DiscussionRouter
	hr *HistoryRecordRouter
}

func NewApp(sr *SeatRouter, cr *CreditRouter, dr *DiscussionRouter, hr *HistoryRecordRouter) *App {
	r := gin.Default()
	sr.SeatRouter(r)
	cr.CreditRouter(r)
	dr.DiscussionRouter(r)
	hr.HistoryRecordRoute(r)

	return &App{
		r:  r,
		sr: sr,
		cr: cr,
		dr: dr,
		hr: hr,
	}
}

func (a *App) Run() {
	a.r.Run()
}
