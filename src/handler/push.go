package handler

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Message struct {
	ToUser  string            `json:"touser"`
	ToParty string            `json:"toparty"`
	ToTag   string            `json:"totag"`
	MsgType string            `json:"msgtype"`
	AgentId int               `json:"agentid"`
	Text    string            `json:"text"`
	Content map[string]string `json:"content"`
	Safe    int               `json:"safe"`
}
type ImageMessage struct {
	ToUser  string            `json:"touser"`
	ToParty string            `json:"toparty"`
	ToTag   string            `json:"totag"`
	MsgType string            `json:"msgtype"`
	AgentId int               `json:"agentid"`
	Image   map[string]string `json:"image"`
	Safe    int               `json:"safe"`
}
type Result struct {
	Errcode      string `json:"errcode"`
	Access_token string `json:"access_token"`
	Errmsg       string `json:"errmsg"`
}

/**
	HTTP 推送文本
 */
func PushController(c *gin.Context) {
	info := WeChatInfo{}
	xin := new(WeChat)

	initWeChat(info)
	content := c.Query("content")
	token := getToken(xin, info)

	log.Println("获取到token为:" + token)

	PushString(token, info.AgentId, content)
	c.String(http.StatusOK, "ok")
}

/**
	HTTP 推送文件
 */
func PushFile(c *gin.Context) {
	name := c.PostForm("name")
	log.Println(name)
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}
	filename := header.Filename

	log.Println(file, err, filename)
	timeLayout := "2006-01-02 15:04:05"
	// 设置时间戳 使用模板格式化为日期字符串
	dataTimeStr := time.Unix(time.Now().Unix(), 0).Format(timeLayout)

	out, err := os.Create("/Users/acke/www/" + dataTimeStr + "-" + filename)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	log.Println(out.Name())
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}
	c.String(http.StatusCreated, "upload successful")

}

/**
	HTTP 健康检查
 */
func HealthCheck(c *gin.Context) {
	message := "OK"
	c.String(http.StatusOK, "\n"+message)
}
