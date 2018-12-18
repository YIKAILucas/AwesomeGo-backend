package handler

import (
	"awesomeProject/src/mongo"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

const DB_NAME = "acke"
const DB_COLLECTION = "test"

var HTTPPubChannel = make(chan HTTPPubDadaSheet, 1000)

type http_recv struct {
	Cmd      string  `json:"cmd" binding:"required"`
	DeviceId string  `json:"device_id" binding:"required"`
	Arg      string  `json:"arg"`
	Timeout  float64 `json:"timeout"`
}

type HTTPPubDadaSheet struct {
	Topic   string
	Payload map[string]interface{}
}

// 发送控制指令到设备(即：发布一条消息到MQTT)
func DeviceControl(c *gin.Context) {
	var recv http_recv
	err := c.ShouldBindJSON(&recv)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "数据解析失败", "err": err})
	}
	json_data := make(map[string]interface{})
	json_data["id"] = time.Now().UnixNano() / 1e3 // ID为16位时间戳
	json_data["cmd"] = recv.Cmd
	if recv.Arg != "" {
		json_data["arg"] = recv.Arg
	}
	if recv.Timeout != 0 {
		json_data["timeout"] = recv.Timeout
	}
	var pub_data HTTPPubDadaSheet
	pub_data.Topic = fmt.Sprintf("tf/Attendance/v1/devices/%s/control", recv.DeviceId)
	pub_data.Payload = json_data
	HTTPPubChannel <- pub_data
	fmt.Printf("收到HTTP控制指令，ID: %d，话题：%s\n", json_data["id"], pub_data.Topic)
}

// 获取设备指令结果(即：从数据库中查询结果)
func DeviceInfo(c *gin.Context) {

	id, err := strconv.ParseFloat(c.Param("id"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID格式错误"})
	}
	db_data := make(map[string]interface{})
	mongo.FindOne(DB_NAME, DB_COLLECTION, bson.M{"id": id}, nil, db_data)
	c.JSON(http.StatusOK, db_data)
}

// 由时间戳生成ID： time.Now().UnixNano()/ 1e3
