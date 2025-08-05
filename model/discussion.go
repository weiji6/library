package model

var DiscussionID = []string{
	"100455304", // 主馆图书馆研讨间
	"103915682", // 南湖图书馆研讨间
}

type Discussion struct {
	LabID    string `json:"labId"`
	LabName  string `json:"labName"`
	KindID   string `json:"kindId"`
	KindName string `json:"kindName"`
	DevID    string `json:"devId"`
	DevName  string `json:"devName"`
	TS       []DiscussionTS
}

type DiscussionTS struct {
	Start  string `json:"start"`
	End    string `json:"end"`
	State  string `json:"state"`
	Title  string `json:"title"`
	Owner  string `json:"owner"`
	Occupy bool   `json:"occupy"`
}

type Search struct {
	ID    string `json:"id"`
	Pid   string `json:"Pid"`
	Name  string `json:"name"`
	Label string `json:"label"`
}
