package handler

import (
	"awesomeProject/src/service"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/imroc/req"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

type WeChat interface {
	getToken()
	PushString(token string, agentId int, content string)
}

/*
企业微信获取token
 */
func WechatGetToken(weChat *WeChat, info *service.CorpWeChatInfo) string {
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

// TODO 测试使用req库
func WechatUploadFile(weChat *WeChat, info *service.CorpWeChatInfo, token string, fileType string, file string) {
	param := req.Param{
		"access_token": token,
		"type":         fileType,
	}
	l, _ := os.Create(file)
	fileUpload := req.FileUpload{
		FileName: "123",
		File:     l,
	}

	r, err := req.Post(info.FileURL, param, fileUpload)
	if err != nil {
		log.Fatal(err)
	}
	var foo map[string]string
	err = r.ToJSON(&foo)
	if err != nil {
		log.Println(err)
	}

}

/*
企业微信上传文件
 */
func NewUploadRequest(URL string, name, path string) (string, error) {
	client := &http.Client{}

	body := &bytes.Buffer{}                 // 初始化body参数
	bodyWriter := multipart.NewWriter(body) // 实例化multipart

	// 创建multipart 文件字段
	fileWriter, err := bodyWriter.CreateFormFile(name, path)
	if err != nil {
		return "", err
	}
	// 打开文件句柄
	fp, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer fp.Close()

	// 写入文件数据到multipart
	_, err = io.Copy(fileWriter, fp)
	if err != nil {
		log.Fatal(err)
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	req, err := http.Post(URL, contentType, body) // 新建请求
	if err != nil {
		return "", err
	}
	defer req.Body.Close()

	reqBody, err := ioutil.ReadAll(req.Body)
	var result map[string]string
	_ = json.Unmarshal(reqBody, &result)

	mediaId := result["media_id"]

	_ = client

	//req.Header.Set("Content-Type", "multipart/form-data") // 设置请求头,!!!非常重要，否则远端无法识别请求
	//response, err := client.Do(req)
	if mediaId == "" {
		var val error = errors.New("没有获取到文件")
		return mediaId, val
	}
	return mediaId, err
}

/*
企业微信推送文件
 */
func WechatPushFile(info *service.CorpWeChatInfo, token string, agentId int, content string) {
	log.Println("token:" + token)
	var fi File = File{
		ToUser:  "@all",
		ToParty: "",
		ToTag:   "",
		MsgType: "file",
		AgentId: agentId,
		Content: map[string]string{"media_id": content},
		Safe:    0,
	}

	param := req.Param{
		"access_token": token,
	}

	r, err := req.Post(info.StringURL, req.BodyJSON(&fi), param)
	if err != nil {
	}
	x := r.Response().Header
	_ = x
}

/*
微信推送文本
 */
func WechatPushString(info *service.CorpWeChatInfo, token string, agentId int, content string) {
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

	r, err := req.Post(info.StringURL, req.BodyJSON(&mes), param)
	if err != nil {
	}
	x := r.Response().Header
	_ = x

}
