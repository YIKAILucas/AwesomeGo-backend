package main

import (
	"encoding/json"
	"github.com/imroc/req"
	"io/ioutil"
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

func getToken() string {
	url := "https://qyapi.weixin.qq.com/cgi-bin/gettoken"
	param := req.Param{
		"corpid":     "ww06bd2f666a354c94",
		"corpsecret": "UINmPVLShl4xDGs1kWfX8dzipbSf45SE2GyVDHWf2ZY",
	}

	r, err := req.Get(url, param)
	if err != nil {
		return ""
	}
	body := r.Response().Body
	b, _ := ioutil.ReadAll(body)
	str := string(b)
	var mapResult map[string]string
	err = json.Unmarshal([]byte(str), &mapResult)
	if err != nil {

	}

	x := mapResult["access_token"] //存在
	return x
}
func Push(token string, agentId int, content string) {

	var mes Message = Message{
		ToUser:  "@all",
		ToParty: "",
		ToTag:   "",
		MsgType: "text",
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
	code := r.Response().StatusCode
	if code != 200 {

	}
}
func main() {
	content := "超级测试"

	Push(getToken(), 1000002, content)

}
