package router

import "github.com/gin-gonic/gin"

type App struct {
	r  *gin.Engine
	sr *SeatRouter
	cr *CreditRouter
	dr *DiscussionRouter
}

func NewApp(sr *SeatRouter, cr *CreditRouter, dr *DiscussionRouter) *App {
	r := gin.Default()
	sr.SeatRouter(r)
	cr.CreditRouter(r)
	dr.DiscussionRouter(r)

	return &App{
		r:  r,
		sr: sr,
		cr: cr,
		dr: dr,
	}
}

func (a *App) Run() {
	a.r.Run()
}
