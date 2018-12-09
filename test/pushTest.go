package main

import (
	"awesomeProject/src/handler"
)

func main() {
	wechat := new(handler.WeChat)
	info := handler.WeChatInfo{}
	handler.InitWeChatInfo(&info)
	//handler.WechatUploadFile(wechat, info, handler.WechatGetToken(wechat, info), "file", `/Users/acke/www/x.jpeg`)

	mediaID := handler.NewUploadRequest("https://qyapi.weixin.qq.com/cgi-bin/media/upload?access_token=vFH98CkkdT2HyESApuAUVH7zsEmOQcQuq2bL4t6TLzNHbaORByvk0eetq716GxcDlzInIna6bTWdssYJEJSr1c9WlKO32Z4yhXtgvsu9FJb7_1qxyaJ9SwJSiR9Uc3MrkOrwb6mjfUnMcBSTN2d6FsjK9xC4bwirWm2OuS_tz5S6zJi-_b0_7I5bXbC9FjG0vUfffajseTaMiQBR-OqEiA&type=file", "mq", "/Users/acke/www/x.jpeg")

	handler.WechatPushFile(info, handler.WechatGetToken(wechat, info), info.AgentId, mediaID)
}
