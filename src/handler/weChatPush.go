package handler

import (
	"bytes"
	"encoding/json"
	"github.com/imroc/req"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type WeChatInfo struct {
	TokenURL   string
	StringURL  string
	FileURL    string
	CorpId     string
	CorpSecret string
	AgentId    int
}

type CompanyFactory interface {
	Work(task *string)
}

type WeChat interface {
	getToken()
	PushString(token string, agentId int, content string)
}

func initWeChat(chat WeChatInfo) {
	//chat.TokenURL = "https://qyapi.weixin.qq.com/cgi-bin/gettoken"
	//chat.FileURL = "https://qyapi.weixin.qq.com/cgi-bin/media/upload"

	//chat.CorpId = "ww4a3407dd4c27e725"
	//chat.CorpSecret = "cOO2IqQXZGwSLYHdaRoMuwo0Bhk4bvrqBH4httj_Vv8"
}

/**
	企业微信获取token
 */
func getToken(weChat *WeChat, info WeChatInfo) string {
	param := req.Param{
		"corpid":     info.CorpId,
		"corpsecret": info.CorpSecret,
	}

	r, err := req.Get(info.TokenURL, param)
	if err != nil {
		return ""
	}
	body := r.Response().Body
	b, _ := ioutil.ReadAll(body)
	str := string(b)
	var mapResult map[string]string
	err = json.Unmarshal([]byte(str), &mapResult)
	if err != nil {
		log.Println("json解析错误")
		log.Println(err)
	}
	if mapResult["errmsg"] != "ok" {
		log.Println(mapResult)
		return "请求失败,错误码为:" + mapResult["errmsg"]
	}

	if _, ok := mapResult["access_token"]; !ok {
		return "token不存在"
	}
	return mapResult["access_token"]
}

func UploadFile(weChat *WeChat, info WeChatInfo, token string, fileType string) {
	header := req.Header{
		"Accept": "multipart/form-data",
	}
	param := req.Param{
		"access_token": token,
		"type":         fileType,
	}
	r, err := req.Post(info.FileURL, header, param)
	if err != nil {
		log.Fatal(err)
	}
	var foo map[string]string
	err = r.ToJSON(&foo)
	panic(err)

}

func NewUploadRequest(link string, params map[string]string, name, path string) (*http.Request, error) {
	fp, err := os.Open(path) // 打开文件句柄
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	body := &bytes.Buffer{}                                       // 初始化body参数
	writer := multipart.NewWriter(body)                           // 实例化multipart
	part, err := writer.CreateFormFile(name, filepath.Base(path)) // 创建multipart 文件字段
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, fp) // 写入文件数据到multipart
	for key, val := range params {
		_ = writer.WriteField(key, val) // 写入body中额外参数，比如七牛上传时需要提供token
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", link, body) // 新建请求
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "multipart/form-data") // 设置请求头,!!!非常重要，否则远端无法识别请求
	return req, nil
}

/**
	微信推送文本
 */
func PushString(token string, agentId int, content string) {
	msgType := "text"
	log.Println("token:" + token)
	var mes Message = Message{
		ToUser:  "@all",
		ToParty: "",
		ToTag:   "",
		MsgType: msgType,
		AgentId: agentId,
		Text:    "Text",
		Content: map[string]string{"content": content},
		Safe:    0,
	}

	param := req.Param{
		"access_token": token,
	}
	sendUrl := "https://qyapi.weixin.qq.com/cgi-bin/message/send"

	r, err := req.Post(sendUrl, req.BodyJSON(&mes), param)
	if err != nil {
	}
	x := r.Response().Header
	_ = x

}
