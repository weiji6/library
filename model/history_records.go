package model

type HistoryRecord struct {
	Place      string `json:"place"`
	Floor      string `json:"floor"`
	Status     string `json:"status"`
	Date       string `json:"date"`
	SubmitTime string `json:"submitTime"`
}
