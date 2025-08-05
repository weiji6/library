package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"library/api/request"
	"library/api/response"
	"library/model"
	"library/tool"
	"net/http"
	"strings"

	"github.com/redis/go-redis/v9"
)

type Discussion interface {
	GetDiscussion(ClassID string, Date string) ([]model.Discussion, error)
	SearchUser(StudentId string) (model.Search, error)
	ReserveDiscussion(reserve request.ReserveDiscussion) (string, error)
}

type DiscussionImpl struct {
	rc *redis.Client
}

func NewDiscussionServiceImpl(rc *redis.Client) *DiscussionImpl {
	return &DiscussionImpl{
		rc: rc,
	}
}

func (ds *DiscussionImpl) GetDiscussion(ClassID string, Date string) ([]model.Discussion, error) {
	ls := tool.GetLoginService()

	fullURL := fmt.Sprintf("http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/device.aspx?byType=devcls&classkind=1&display=cld&md=d&class_id=%s&cld_name=default&date=%s&act=get_rsv_sta", ClassID, Date)

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := ls.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var discussionResp response.GetDiscussion
	if err = json.Unmarshal(body, &discussionResp); err != nil {
		return nil, fmt.Errorf("解析座位信息失败: %v", err)
	}

	if discussionResp.Ret != 1 {
		return nil, fmt.Errorf(discussionResp.Msg)
	}

	return discussionResp.Data, nil
}

func (ds *DiscussionImpl) SearchUser(StudentId string) (model.Search, error) {
	ls := tool.GetLoginService()

	fullURL := fmt.Sprintf("http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/data/searchAccount.aspx?type=logonname&ReservaApply=ReservaApply&term=%s", StudentId)

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return model.Search{}, err
	}

	res, err := ls.Client.Do(req)
	if err != nil {
		return model.Search{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return model.Search{}, err
	}

	var searchResp []model.Search
	if err = json.Unmarshal(body, &searchResp); err != nil {
		return model.Search{}, err
	}

	return searchResp[0], nil
}

func (ds *DiscussionImpl) ReserveDiscussion(discussion request.ReserveDiscussion) (string, error) {
	ls := tool.GetLoginService()

	mbList := "$" + strings.Join(discussion.List, ",")

	fullURL := fmt.Sprintf("http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/reserve.aspx?dev_id=%s&lab_id=%s&kind_id=%s&type=dev&test_name=%s&min_user=3&max_user=4&mb_list=%s&start=%s&end=%s&act=set_resv", discussion.DevID, discussion.LabID, discussion.KindID, discussion.Title, mbList, discussion.Start, discussion.End)

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return "", err
	}

	res, err := ls.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var reserveResp response.Reserve
	if err = json.Unmarshal(body, &reserveResp); err != nil {
		return "", err
	}

	if reserveResp.Ret != 1 {
		return "", fmt.Errorf(reserveResp.Msg)
	}

	return reserveResp.Msg, nil
}
