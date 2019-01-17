package main

import (
	"awesomeProject/src/middleware/mongo"
	"encoding/json"
	"fmt"
	"github.com/imroc/req"
	"github.com/robfig/cron"
	"io/ioutil"
	"time"
)

type Test struct {
}

func (test *Test) Run() {

}

func main() {
	local, _ := time.LoadLocation("Asia/Shanghai") //服务器上设置的时区

	// 定时任务
	c := cron.NewWithLocation(local)
	_ = c.AddFunc("30 0-59 * * * ?", func() {
		timeStr := time.Now().Format("2006-01-02 15:04:05")
		fmt.Println(timeStr)

		x, _ := getActiveDevice()
		addDAU(x)
	})

	c.Start()

	select {}
}

func addDAU(count int) {
	doc := map[string]interface{}{
		"dau_count": count,
		"date_time": time.Now().Format("2006-01-02 15:04:05"),
	}

	_ = mongo.Insert(mongo.DB_NAME, mongo.DAU, doc)

}

type Info struct {
	Result []map[string]int `json:"result"`
}

func getActiveDevice() (int, error) {
	url := "http://106.12.130.179:18083/api/v2/monitoring/nodes"

	headers := req.Header{
		"Authorization": "Basic YWRtaW46cHVibGlj",
	}
	param := req.Param{
		"curr_page": 1,
		"page_size": 1000,
	}

	r, err := req.Get(url, headers, param)
	if err != nil {
		return 0, err
	}

	body := r.Response().Body
	b, _ := ioutil.ReadAll(body)
	var info Info

	err = json.Unmarshal(b, &info)

	result := info.Result
	return result[0]["clients"], nil
}
