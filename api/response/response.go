package response

import "library/model"

type Response struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
}

type GetSeat struct {
	Data []model.Seat `json:"data"`
}

type Reserve struct {
	Ret int    `json:"ret"`
	Msg string `json:"msg"`
}

type Cancel struct {
	Ret int    `json:"ret"`
	Msg string `json:"msg"`
}

type GetDiscussion struct {
	Ret  int                `json:"ret"`
	Msg  string             `json:"msg"`
	Data []model.Discussion `json:"data"`
}
