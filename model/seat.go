package model

var RoomIDs = []string{
	"100455820", // 主馆图书馆一楼-一楼综合学习室
	"100455822", // 主馆图书馆二楼-二楼借阅室（一）
	"100671994", // 主馆图书馆二楼-二楼借阅室（二）
	"100455824", // 主馆图书馆三楼-三楼借阅室（三）
	"100455826", // 主馆图书馆四楼-四楼自主学习中心
	"100455828", // 主馆图书馆五楼-五楼借阅室（四）
	"100746476", // 主馆图书馆五楼-五楼借阅室（五）
	"100746204", // 主馆图书馆六楼-六楼阅览室（一）
	"100455830", // 主馆图书馆六楼-六楼外文借阅室
	"100455832", // 主馆图书馆七楼-七楼阅览室（二）
	"100746480", // 主馆图书馆七楼-七楼阅览室（三）
	"100455834", // 主馆图书馆九楼-九楼阅览室
	"101699179", // 南湖分馆一楼-南湖分馆一楼开敞座位区
	"101699187", // 南湖分馆一楼-南湖分馆一楼中庭开敞座位区
	"101699189", // 南湖分馆二楼-南湖分馆二楼开敞座位区
	"101699191", // 南湖分馆二楼-南湖分馆二楼卡座区
}

type Seat struct {
	LabName  string   `json:"labName"`
	KindName string   `json:"kindName"`
	DevID    string   `json:"devId"`
	DevName  string   `json:"devName"`
	TS       []SeatTS `json:"ts"`
}

type SeatTS struct {
	Start  string `json:"start"`
	End    string `json:"end"`
	State  string `json:"state"`
	Owner  string `json:"owner"`
	Occupy bool   `json:"occupy"`
}

type Parsed struct {
	Data []Record `json:"data"`
}

type Record struct {
	ID       string `json:"id"`
	Owner    string `json:"owner"`
	Start    string `json:"start"`
	End      string `json:"end"`
	TimeDesc string `json:"timeDesc"`
	Occur    string `json:"occur"`
	States   string `json:"states"`
	DevName  string `json:"devName"`
	RoomID   string `json:"roomId"`
	RoomName string `json:"roomName"`
	LabName  string `json:"labName"`
}
