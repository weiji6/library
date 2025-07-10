package request

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Reserve struct {
	DevID string `json:"dev_id"`
	Start string `json:"start"`
	End   string `json:"end"`
}

type ReserveDiscussion struct {
	DevID  string   `json:"dev_id"`
	LabID  string   `json:"lab_id"`
	KindID string   `json:"kind_id"`
	Title  string   `json:"title"`
	List   []string `json:"list"`
	Start  string   `json:"start"`
	End    string   `json:"end"`
}
