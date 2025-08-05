package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"library/api/request"
	"library/api/response"
	"library/model"
	"library/tool"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/redis/go-redis/v9"
)

type SeatService interface {
	FetchSeat(RoomID string) ([]model.Seat, error)
	FetchAllSeats() (map[string][]model.Seat, error)
	StartSeatUpdateService(ctx context.Context, redisClient *redis.Client)
	ReserveSeat(message request.Reserve) (string, error)
	GetRecord() (model.Parsed, error)
	CancelReserve(ID string) (string, error)
}

type SeatServiceImpl struct {
	rc *redis.Client
}

func NewSeatServiceImpl(rc *redis.Client) *SeatServiceImpl {
	return &SeatServiceImpl{
		rc: rc,
	}
}

func (ss *SeatServiceImpl) FetchSeat(RoomID string) ([]model.Seat, error) {
	ls := tool.GetLoginService()

	now := time.Now()
	date := now.Format("2006-01-02")

	minute := now.Minute()
	nextMinute := (minute/5 + 1) * 5
	if nextMinute == 60 {
		nextMinute = 0
		now = now.Add(time.Hour) // 进位到下一个小时
	}
	nextTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), nextMinute, 0, 0, now.Location())

	fullURL := fmt.Sprintf("http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/device.aspx?byType=devcls&classkind=8&display=fp&md=d&room_id=%s&cld_name=default&date=%s&fr_start=%s&fr_end=22:00&act=get_rsv_sta", RoomID, date, nextTime.Format("15:04"))

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := ls.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var seatResp response.GetSeat
	if err = json.Unmarshal(body, &seatResp); err != nil {
		return nil, fmt.Errorf("解析座位信息失败: %v", err)
	}

	return seatResp.Data, nil
}

func (ss *SeatServiceImpl) FetchAllSeats() (map[string][]model.Seat, error) {
	var wg sync.WaitGroup
	results := make(map[string][]model.Seat)
	mutex := &sync.Mutex{}

	for _, roomID := range model.RoomIDs {
		wg.Add(1)
		go func(roomID string) {
			defer wg.Done()
			seats, err := ss.FetchSeat(roomID)
			if err != nil {
				fmt.Printf("获取房间 %s 座位失败: %v", roomID, err)
				mutex.Lock()
				results[roomID] = nil
				mutex.Unlock()
				return
			}
			mutex.Lock()
			results[roomID] = seats
			mutex.Unlock()
		}(roomID)
	}

	wg.Wait()
	return results, nil
}

func (ss *SeatServiceImpl) StartSeatUpdateService(ctx context.Context, redisClient *redis.Client) {
	ticker := time.NewTicker(5 * time.Minute) // 每5分钟更新一次数据
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			seats, err := ss.FetchAllSeats()
			if err != nil {
				fmt.Printf("更新座位数据失败: %v\n", err)
				continue
			}

			// 将数据存入Redis
			seatsJSON, err := json.Marshal(seats)
			if err != nil {
				fmt.Printf("序列化座位数据失败: %v\n", err)
				continue
			}

			err = redisClient.Set(ctx, "all_seats", seatsJSON, 10*time.Minute).Err()
			if err != nil {
				fmt.Printf("存储座位数据到Redis失败: %v\n", err)
			}
		}
	}
}

// todo:实现SSE监听抢座
func (ss *SeatServiceImpl) ReserveSeat(message request.Reserve) (string, error) {
	ls := tool.GetLoginService()

	// 2025-07-08+19:00
	// 2025-07-08+20:00

	fullURL := fmt.Sprintf("http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/reserve.aspx?dev_id=%s&type=dev&start=%s&end=%s&act=set_resv", message.DevID, message.Start, message.End)

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

	var ReserveResp response.Reserve
	if err = json.Unmarshal(body, &ReserveResp); err != nil {
		return "", err
	}

	if ReserveResp.Ret != 1 {
		return "", fmt.Errorf(ReserveResp.Msg)
	}

	return ReserveResp.Msg, nil
}

func (ss *SeatServiceImpl) GetRecord() (model.Parsed, error) {
	var record model.Parsed

	ls := tool.GetLoginService()

	fullRUL := fmt.Sprintf("http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/reserve.aspx?stat_flag=9&act=get_my_resv")

	req, err := http.NewRequest("GET", fullRUL, nil)
	if err != nil {
		return model.Parsed{}, err
	}

	res, err := ls.Client.Do(req)
	if err != nil {
		return model.Parsed{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return model.Parsed{}, err
	}

	fmt.Println(string(body))

	if err = json.Unmarshal(body, &record); err != nil {
		return model.Parsed{}, err
	}

	for i := range record.Data {
		raw := record.Data[i].States
		html := "<div>" + raw + "</div>"

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {
			continue
		}

		var plainStates []string
		doc.Find("span").Each(func(i int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			if text != "" {
				plainStates = append(plainStates, text)
			}
		})

		record.Data[i].States = strings.Join(plainStates, ",")
	}

	return record, nil
}

func (ss *SeatServiceImpl) CancelReserve(ID string) (string, error) {
	ls := tool.GetLoginService()

	fullURL := fmt.Sprintf("http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/reserve.aspx?act=del_resv&id=%s", ID)

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

	jsonRegexp := regexp.MustCompile(`\{[^}]+}`)
	matches := jsonRegexp.FindAll(body, -1)

	var cancelResp response.Cancel

	for _, m := range matches {
		if err = json.Unmarshal(m, &cancelResp); err != nil {
			continue // 忽略无效块
		}
		if cancelResp.Ret == 1 {
			return cancelResp.Msg, nil
		}
		return "", fmt.Errorf(cancelResp.Msg)
	}

	return cancelResp.Msg, nil
}
