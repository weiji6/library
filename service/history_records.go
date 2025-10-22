package service

import (
	"library/model"
	"library/tool"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type HistoryRecords interface {
	GetHistoryRecords() ([]model.HistoryRecord, error)
}

type RecordsServiceImpl struct{}

func NewHistoryRecordsServiceImpl() *RecordsServiceImpl {
	return &RecordsServiceImpl{}
}

func (rs *RecordsServiceImpl) GetHistoryRecords() ([]model.HistoryRecord, error) {
	ls := tool.GetLoginService()

	fullURL := "http://kjyy.ccnu.edu.cn/clientweb/m/a/resvlist.aspx"

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("origin", "https://account.ccnu.edu.cn")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")

	res, err := ls.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var records []model.HistoryRecord

	doc.Find("li.item-content").Each(func(i int, item *goquery.Selection) {
		place := item.Find(".item-title").Text()
		status := item.Find(".item-after").Text()
		date := item.Find(".item-subtitle").Text()
		submitText := item.Find(".item-text").Text()
		submitParts := strings.Split(submitText, ",")
		if len(submitParts) >= 2 {
			floor := submitParts[0]
			floor = strings.TrimSpace(floor)
			submitTime := submitParts[2]
			submitTime = strings.TrimSpace(submitTime)

			records = append(records, model.HistoryRecord{
				Place:      place,
				Floor:      floor,
				Status:     status,
				Date:       date,
				SubmitTime: submitTime,
			})
		}
	})

	return records, nil
}
