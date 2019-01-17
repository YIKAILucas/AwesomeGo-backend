package pushController

import (
	"awesomeProject/src/handler/pushController/constants"
	"encoding/json"
	"errors"
	"github.com/imroc/req"
	"io/ioutil"
)

type DingMessage struct {
	ChatId  string      `json:"chatid"`
	MsgType string      `json:"msgtype"`
	Text    interface{} `json:"text"`
}

func PubMessage(content string) error {
	url := "https://oapi.dingtalk.com/chat/send"
	token := *new(string)
	if getToken(&token) != nil {

	}

	param := req.Param{
		"access_token": token,
	}

	var message DingMessage = DingMessage{
		ChatId:  constants.Chatid,
		MsgType: "text",
		Text: map[string]string{
			"content": content,
		},
	}
	r, err := req.Post(url, req.BodyJSON(&message), param)
	_ = r
	//if r[""] {
	//
	//}
	if err != nil {
		return err
	}

	return nil
}

const tokenUrl = "https://oapi.dingtalk.com/gettoken"

/**
获取ding token
*/
func getToken(token *string) error {
	param := req.Param{
		"appkey":    "dingesahgfr6qayszbxn",
		"appsecret": "uGcD8QfSpyIpaW_GHlgwSmJBnf60KyEmEFSd-IVAru6GsdUJjnKAo7PcRVWU1BJF",
	}

	r, err := req.Get(tokenUrl, param)
	if err != nil {
		return err

	}
	body := r.Response().Body
	b, _ := ioutil.ReadAll(body)
	str := string(b)
	var mapResult map[string]string
	err = json.Unmarshal([]byte(str), &mapResult)
	if err != nil {
		// 用map解析会解析到其他非string类型
		//return err
	}

	*token = mapResult["access_token"]
	return nil
}

func creatChat(chatid *string) error {
	url := "https://oapi.dingtalk.com/chat/create"
	param := req.Param{
		//"access_token": getToken(),
	}
	r, err := req.Post(url, param)
	if err != nil {
		return err
	}
	body := r.Response().Body
	b, _ := ioutil.ReadAll(body)
	str := string(b)
	var mapResult map[string]string
	err = json.Unmarshal([]byte(str), &mapResult)

	if mapResult["chatid"] == "" {
		var err error = errors.New("创建群会话失败")
		return err
	} else {
		*chatid = mapResult["chatid"]
	}

	return nil
}

func getChat() error {
	return nil
}
