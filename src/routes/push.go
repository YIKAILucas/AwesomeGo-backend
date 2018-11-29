package routes

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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

func PushController() *gin.Engine {
	router := gin.Default()
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	content := "http测试"
	agentId := 1000002
	router.POST("/push", func(c *gin.Context) {
		Push(getToken(), agentId, content)

		c.JSON(http.StatusOK, gin.H{
			"status": gin.H{
				"status_code": http.StatusOK,
				"status":      "ok",
			},
		})
	})
	/**
	分组路由
 	*/
	v1 := router.Group("/v1")

	v1.GET("/login", func(c *gin.Context) {
		c.String(http.StatusOK, "v1 login")
	})

	v2 := router.Group("/v2")

	v2.GET("/login", func(c *gin.Context) {
		c.String(http.StatusOK, "v2 login")
	})

	return router
}
