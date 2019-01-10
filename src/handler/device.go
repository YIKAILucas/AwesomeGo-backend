package handler

import (
	"awesomeProject/src/middleware/mongo"
	"awesomeProject/src/middleware/mqttbroker"
	"awesomeProject/src/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

// TODO: 将此部分配置放到配置文件中去
var emq_config = map[string]string{
	"host":     "http://106.12.130.179:18083",
	"username": "admin",
	"password": "public",
	"node":     "e6d70b6bfd92@172.17.0.7",
}

var HTTPPubChannel = make(chan HTTPPubDadaSheet, 1000)

type http_recv struct {
	Cmd      string  `json:"cmd" binding:"required"`
	DeviceId string  `json:"device_id" binding:"required"`
	Arg      string  `json:"arg"`
	Timeout  float64 `json:"timeout"`
}

type HTTPPubDadaSheet struct {
	Topic   string
	Payload []byte
}

// 发送控制指令到设备(即：发布一条消息到MQTT)
func DeviceControl(c *gin.Context) {
	var recv http_recv
	err := c.ShouldBindJSON(&recv)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "数据解析失败", "err": err})
		return
	}
	json_data := make(map[string]interface{})
	timestamp := time.Now().UnixNano() / 1e3 // ID为16位时间戳
	json_data["id"] = timestamp
	json_data["cmd"] = recv.Cmd
	if recv.Arg != "" {
		json_data["arg"] = recv.Arg
	}
	if recv.Timeout != 0 {
		json_data["timeout"] = recv.Timeout
	}
	// TODO: 这里写了一次数据库，Pub后又写入了一次数据库？
	mongo.Insert(mqttbroker.DB_NAME, mqttbroker.CMD_COLLECTION_MAP["default"], map[string]interface{}{"id": timestamp, "cmd": recv.Cmd}) // 预先插入一条记录，防止指令返回时异步条目还没有插入的问题

	var pub_data HTTPPubDadaSheet
	pub_data.Topic = fmt.Sprintf("tf/Attendance/v1/devices/%s/control", recv.DeviceId)
	pub_data.Payload, _ = json.Marshal(json_data)
	HTTPPubChannel <- pub_data
	c.JSON(http.StatusOK, gin.H{"id": json_data["id"]})
	fmt.Printf("收到HTTP控制指令，ID: %d，话题：%s\n", json_data["id"], pub_data.Topic)
}

// 获取设备指令结果(即：从数据库中查询结果)
func DeviceInfo(c *gin.Context) {

	id, err := strconv.ParseFloat(c.Param("id"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID格式错误"})
		return
	}
	cmd_data := make(map[string]interface{})
	mongo.FindOne(mqttbroker.DB_NAME, mqttbroker.CMD_COLLECTION_MAP["default"], bson.M{"id": id}, nil, cmd_data)
	if len(cmd_data) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "未找到对应的记录！"})
	} else {
		result_data := make(map[string]interface{})
		result_data["success"] = cmd_data["success"]
		result_data["message"] = cmd_data["message"]

		cmd_name := "default"
		cmd_name, ok := cmd_data["cmd"].(string)
		var result_collection = mqttbroker.CMD_COLLECTION_MAP["default"]
		_, ok = mqttbroker.CMD_COLLECTION_MAP[cmd_name]
		if !ok {
			// 未指定结果表时，从默认表中取
			result_data["last_update_time"] = cmd_data["last_update_time"]
			result_data["result"] = cmd_data["result"]
		} else {
			// 已指定表时，从相应的表中取
			result_collection = mqttbroker.CMD_COLLECTION_MAP[cmd_name]
			db_data := make(map[string]interface{})
			mongo.FindOne(mqttbroker.DB_NAME, result_collection, bson.M{"device_id": cmd_data["device_id"]}, nil, db_data)
			result_data["last_update_time"] = db_data["last_update_time"]
			result_data["result"] = db_data["result"]
		}
		c.JSON(http.StatusOK, result_data)
		return
	}
}

func DeviceList(c *gin.Context) {
	/* 获取设备列表 */
	var devices []model.Device
	var results []map[string]string

	model.DB.Find(&devices)
	for _, device := range devices {
		results = append(results, map[string]string{"device_id": device.DeviceId, "device_name": device.DeviceName})
	}
	c.JSON(http.StatusOK, results)
}

func DeviceOnlineStatus(c *gin.Context) {
	/* 获取某台设备的在线情况 */
	device_id := c.Param("id")
	url := fmt.Sprintf("%s/api/v2/nodes/%s/clients/%s", emq_config["host"], emq_config["node"], device_id)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(emq_config["username"], emq_config["password"])
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("获取设备%s的在线情况失败：%s\n", device_id, err)
		c.JSON(http.StatusOK, gin.H{"success": false, "msg": "请求后端失败", "result": nil})
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		resp_body, _ := ioutil.ReadAll(resp.Body)
		resp_json := make(map[string]interface{})
		err = json.Unmarshal(resp_body, &resp_json)
		if err != nil {
			fmt.Printf("获取设备%s的在线情况失败：%s\n", device_id, err)
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": "json解码失败", "result": nil})
		}
		results, _ := resp_json["result"].(map[string]interface{})
		client_info, _ := results["objects"].([]interface{})
		if len(client_info) > 0 {
			c.JSON(http.StatusOK, gin.H{"success": true, "msg": "", "result": true})
		} else {
			c.JSON(http.StatusOK, gin.H{"success": true, "msg": "", "result": false})
		}

	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Printf("获取设备%s的在线情况失败，后端认证失败\n", device_id)
		c.JSON(http.StatusOK, gin.H{"success": false, "msg": "后端认证失败", "result": nil})
	} else {
		fmt.Printf("获取设备%s的在线情况失败，状态码：%d\n", device_id, resp.StatusCode)
		c.JSON(http.StatusOK, gin.H{"success": false, "msg": "未知错误", "result": nil})
	}
}
