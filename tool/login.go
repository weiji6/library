package tool

import (
	"fmt"
	"io/ioutil"
	"library/api/request"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
)

var (
	globalLoginService *LoginServiceImpl
	once               sync.Once
)

type LoginService interface {
	LoginFirst(loginRequest request.Login) error
	LoginAuto() error
	LoginSecond() error
}

type LoginServiceImpl struct {
	LoginService
	Client *http.Client
}

func NewLoginServiceImpl() *LoginServiceImpl {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic("创建 cookie jar 失败")
	}

	return &LoginServiceImpl{
		Client: &http.Client{
			Jar: jar,
		},
	}
}

func (ls *LoginServiceImpl) LoginFirst(loginRequest request.Login) error {
	var Lt, Execution, Cookies string

	URL1 := "https://account.ccnu.edu.cn/cas/login?service=http://kjyy.ccnu.edu.cn/loginall.aspx?page="
	res, err := ls.Client.Get(URL1)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	lt := regexp.MustCompile(`name="lt"\s+value="([^"]+)"`)
	execution := regexp.MustCompile(`name="execution"\s+value="([^"]+)"`)
	cookies := regexp.MustCompile(`JSESSIONID=([^;]+)`)

	match1 := lt.FindStringSubmatch(string(body))
	if len(match1) > 1 {
		// fmt.Println("lt value:", match1[1])
		Lt = match1[1]
	} else {
		fmt.Println("No match found")
	}

	match2 := execution.FindStringSubmatch(string(body))
	if len(match2) > 1 {
		// fmt.Println("execution value:", match2[1])
		Execution = match2[1]
	} else {
		fmt.Println("No match found")
	}

	match3 := cookies.FindStringSubmatch(fmt.Sprintf("%v", res.Header))
	if len(match3) > 1 {
		// fmt.Println("cookies value:", match3[1])
		Cookies = "JSESSIONID" + match3[1]
	} else {
		fmt.Println("No match found")
	}

	data := url.Values{}
	data.Set("username", loginRequest.Username)
	data.Set("password", loginRequest.Password)
	data.Set("lt", Lt)
	data.Set("execution", Execution)
	data.Set("_eventId", "submit")
	data.Set("submit", "登录")

	playload := strings.NewReader(data.Encode())

	URL2 := "https://account.ccnu.edu.cn/cas/login;jsessionid=" + match3[1] + "?service=http://kjyy.ccnu.edu.cn/loginall.aspx?page="
	req, err := http.NewRequest("POST", URL2, playload)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", Cookies)

	resp, err := ls.Client.Do(req)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	defer resp.Body.Close()

	return nil
}

func (ls *LoginServiceImpl) LoginSecond() error {
	var Lt, Execution, Cookies string

	URL1 := "https://account.ccnu.edu.cn/cas/login?service=http://kjyy.ccnu.edu.cn/loginall.aspx?page="
	res, err := ls.Client.Get(URL1)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	lt := regexp.MustCompile(`name="lt"\s+value="([^"]+)"`)
	execution := regexp.MustCompile(`name="execution"\s+value="([^"]+)"`)
	cookies := regexp.MustCompile(`JSESSIONID=([^;]+)`)

	match1 := lt.FindStringSubmatch(string(body))
	if len(match1) > 1 {
		// fmt.Println("lt value:", match1[1])
		Lt = match1[1]
	} else {
		return fmt.Errorf("lt 获取失败")
	}

	match2 := execution.FindStringSubmatch(string(body))
	if len(match2) > 1 {
		// fmt.Println("execution value:", match2[1])
		Execution = match2[1]
	} else {
		return fmt.Errorf("execution 获取失败")
	}

	match3 := cookies.FindStringSubmatch(fmt.Sprintf("%v", res.Header))
	if len(match3) > 1 {
		// fmt.Println("cookies value:", match3[1])
		Cookies = "JSESSIONID" + match3[1]
	} else {
		return fmt.Errorf("cookies 获取失败")
	}

	username := viper.GetString("user.username")
	password := viper.GetString("user.password")
	// fmt.Println("username:", username, "password:", password)

	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)
	data.Set("lt", Lt)
	data.Set("execution", Execution)
	data.Set("_eventId", "submit")
	data.Set("submit", "登录")

	playload := strings.NewReader(data.Encode())

	URL2 := "https://account.ccnu.edu.cn/cas/login;jsessionid=" + match3[1] + "?service=http://kjyy.ccnu.edu.cn/loginall.aspx?page="
	req, err := http.NewRequest("POST", URL2, playload)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", Cookies)

	resp, err := ls.Client.Do(req)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	defer resp.Body.Close()

	return nil
}

func (ls *LoginServiceImpl) LoginAuto() error {
	ticker := time.NewTicker(30 * time.Minute) // 设置30分钟登录间隙
	defer ticker.Stop()

	for {
		var Lt, Execution, Cookies string

		URL1 := "https://account.ccnu.edu.cn/cas/login?service=http://kjyy.ccnu.edu.cn/loginall.aspx?page="
		res, err := ls.Client.Get(URL1)
		if err != nil {
			return fmt.Errorf("client.Get error:" + err.Error())
		}
		res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("ioutil.ReadAll error:" + err.Error())
		}

		lt := regexp.MustCompile(`name="lt"\s+value="([^"]+)"`)
		execution := regexp.MustCompile(`name="execution"\s+value="([^"]+)"`)
		cookies := regexp.MustCompile(`JSESSIONID=([^;]+)`)

		match1 := lt.FindStringSubmatch(string(body))
		if len(match1) > 1 {
			// fmt.Println("lt value:", match1[1])
			Lt = match1[1]
		} else {
			return fmt.Errorf("lt 获取失败")
		}

		match2 := execution.FindStringSubmatch(string(body))
		if len(match2) > 1 {
			// fmt.Println("execution value:", match2[1])
			Execution = match2[1]
		} else {
			return fmt.Errorf("execution 获取失败")
		}

		match3 := cookies.FindStringSubmatch(fmt.Sprintf("%v", res.Header))
		if len(match3) > 1 {
			// fmt.Println("cookies value:", match3[1])
			Cookies = "JSESSIONID" + match3[1]
		} else {
			return fmt.Errorf("cookies 获取失败")
		}

		username := viper.GetString("user.username")
		password := viper.GetString("user.password")

		data := url.Values{}
		data.Set("username", username)
		data.Set("password", password)
		data.Set("lt", Lt)
		data.Set("execution", Execution)
		data.Set("_eventId", "submit")
		data.Set("submit", "登录")

		playload := strings.NewReader(data.Encode())

		URL2 := "https://account.ccnu.edu.cn/cas/login;jsessionid=" + match3[1] + "?service=http://kjyy.ccnu.edu.cn/loginall.aspx?page="
		req, err := http.NewRequest("POST", URL2, playload)
		if err != nil {
			return fmt.Errorf("登陆失败:%v", err.Error())
		}

		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Cookie", Cookies)

		resp, err := ls.Client.Do(req)
		if err != nil {
			return fmt.Errorf("登陆失败:%v", err.Error())
		}
		resp.Body.Close()

		<-ticker.C // 等待下一次定时触发
	}
}

func GetLoginService() *LoginServiceImpl {
	once.Do(func() {
		globalLoginService = NewLoginServiceImpl()
		if err := globalLoginService.LoginSecond(); err != nil {
			fmt.Println("登陆失败")
			panic(err)
		}
	})

	return globalLoginService
}
