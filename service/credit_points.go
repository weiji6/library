package service

import (
	"library/model"
	"library/tool"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type CreditPoints interface {
	GetCreditPoints() (*model.CreditPoints, error)
}

type CreditServiceImpl struct{}

func NewCreditServiceImpl() *CreditServiceImpl {
	return &CreditServiceImpl{}
}

func (cs *CreditServiceImpl) GetCreditPoints() (*model.CreditPoints, error) {
	ls := tool.GetLoginService()

	fullURL := "http://kjyy.ccnu.edu.cn/clientweb/m/a/credit.aspx"

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := ls.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var summary model.CreditSummary
	doc.Find("table tbody tr").Each(func(i int, s *goquery.Selection) {
		tds := s.Find("td")
		if tds.Length() >= 3 {
			summary = model.CreditSummary{
				System: strings.TrimSpace(tds.Eq(0).Text()),
				Remain: strings.TrimSpace(tds.Eq(1).Text()),
				Total:  strings.TrimSpace(tds.Eq(2).Text()),
			}
		}
	})

	var records []model.CreditRecord
	doc.Find("#my_resv_list li").Each(func(i int, s *goquery.Selection) {
		record := model.CreditRecord{
			Title:    strings.TrimSpace(s.Find(".item-title").Text()),
			Subtitle: strings.TrimSpace(s.Find(".item-subtitle").Text()),
			Location: strings.TrimSpace(s.Find(".item-text").Text()),
		}
		records = append(records, record)
	})

	return &model.CreditPoints{
		Summary: summary,
		Records: records,
	}, nil
}
