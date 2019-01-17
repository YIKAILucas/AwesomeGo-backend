package main

import (
	"awesomeProject/src/handler"
	"awesomeProject/src/service"
	"fmt"
	"net/url"
)

func main() {
	wechat := new(handler.WeChat)
	info := &service.CorpWeChatInfo{}
	service.InitWeChatInfo(info)
	//handler.WechatUploadFile(wechat, info, handler.WechatGetToken(wechat, info), "file", `/Users/acke/www/x.jpeg`)
	URL := "https://qyapi.weixin.qq.com/cgi-bin/media/upload"
	u, _ := url.Parse(URL)
	u.Query().Set("access_token", handler.WechatGetToken(wechat, info))
	u.Query().Set("type", "file")
	fmt.Println(u.String())
	mediaID, _ := handler.NewUploadRequest(u.String(), "mq", "/Users/acke/www/x.jpeg")

	handler.WechatPushFile(info, handler.WechatGetToken(wechat, info), info.AgentId, mediaID)
}
