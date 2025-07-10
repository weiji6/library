package service

import (
	"fmt"
	"io/ioutil"
	"library/tool"
	"net/http"
)

type Records interface {
	GetHistoryRecords() (string, error)
}

type RecordsServiceImpl struct{}

func NewRecordsServiceImpl() *RecordsServiceImpl {
	return &RecordsServiceImpl{}
}

// todo:完成解析
func (rs *RecordsServiceImpl) GetHistoryRecords() (string, error) {
	ls := tool.GetLoginService()

	fullURL := "http://kjyy.ccnu.edu.cn/clientweb/m/a/resvlist.aspx"

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}

	res, err := ls.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}

	return string(body), nil
}
