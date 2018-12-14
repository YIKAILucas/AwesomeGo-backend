package handler

import (
	"awesomeProject/src/service"
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
type File struct {
	ToUser  string            `json:"touser"`
	ToParty string            `json:"toparty"`
	ToTag   string            `json:"totag"`
	MsgType string            `json:"msgtype"`
	AgentId int               `json:"agentid"`
	File    string            `json:"file"`
	Content map[string]string `json:"content"`
	Safe    int               `json:"safe"`
}

type Image struct {
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
	info := &service.CorpWeChatInfo{}
	xin := new(WeChat)
	//fac := &service.CompanyFactory()
	//service.GetCompany(fac)


	service.InitWeChatInfo(info)
	service.CreateCompany()
	content := c.Query("content")
	token := WechatGetToken(xin, info)

	log.Println("获取到token为:" + token)

	WechatPushString(info, token, info.AgentId, content)
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
	info := &service.CorpWeChatInfo{}
	service.InitWeChatInfo(info)
	wechat := new(WeChat)
	mediaID, err := NewUploadRequest("https://qyapi.weixin.qq.com/cgi-bin/media/upload?access_token="+WechatGetToken(wechat, info)+"&type=file", "mq", out.Name())
	if err != nil {
		log.Fatal(err)
	}
	WechatPushFile(info, WechatGetToken(wechat, info), info.AgentId, mediaID)
	c.String(http.StatusCreated, "upload successful")

}

/**
HTTP 健康检查
 */
func HealthCheck(c *gin.Context) {
	message := "OK"
	c.String(http.StatusOK, "\n"+message)
}
