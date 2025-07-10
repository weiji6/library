package main

import "library/config"

// todo:数据持久化
func main() {
	config.InitConfig()

	app, err := InitApp("config/config.yaml")
	if err != nil {
		panic(err)
	}
	app.Run()

	//// 创建登录服务实例
	//loginService := tool.NewLoginServiceImpl()
	//// 创建座位服务实例
	//seatService := service.NewSeatServiceImpl(loginService)
	//// 调用FetchSeat，传入阅览室ID
	//response1, err := seatService.FetchSeat("100455822")
	//// 错误处理
	//if err != nil {
	//	fmt.Printf("Error: %v\n", err)
	//	return
	//}
	//// 打印响应结果
	//fmt.Printf("%+v\n", response1)

	//recordsService := tool.NewRecordsServiceImpl(loginService)
	//
	//response2, err := recordsService.GetHistoryRecords()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%+v\n", response2)

	//CreditPointsService := tool.NewCreditPointsServiceImpl(loginService)
	//
	//response3, err := CreditPointsService.GetCreditPoints()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%+v\n", response3)

	//reserveService := tool.NewReserveServiceImpl(loginService)
	//response4, err := reserveService.ReserveSeat("101699817")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%+v\n", response4)
}
