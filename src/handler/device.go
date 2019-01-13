package handler

import (
	"awesomeProject/src/middleware/mongo"
	"awesomeProject/src/middleware/mqttbroker"
	"awesomeProject/src/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	http_url "net/url"
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

func requestEMQBackend(url string, params http_url.Values) (data map[string]interface{}, err error) {
	/* 请求EMQ后端，请求结果 */
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.URL.RawQuery = params.Encode()
	req.SetBasicAuth(emq_config["username"], emq_config["password"])
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求EMQ后端失败：请求失败")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		resp_body, _ := ioutil.ReadAll(resp.Body)
		resp_json := make(map[string]interface{})
		err = json.Unmarshal(resp_body, &resp_json)
		if err != nil {
			return nil, fmt.Errorf("请求EMQ后端失败：JSON解码失败")
		}
		return resp_json, nil
	} else if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("请求EMQ后端失败：认证失败")
	} else {
		return nil, fmt.Errorf("请求EMQ后端失败：状态码：%d", resp.StatusCode)
	}
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
	var db_results []map[string]string
	//var emq_results []map[string]interface{}
	var url = fmt.Sprintf("%s/api/v2/nodes/%s/clients", emq_config["host"], emq_config["node"])

	model.DB.Find(&devices)
	for _, device := range devices {
		db_results = append(db_results, map[string]string{"device_id": device.DeviceId, "device_name": device.DeviceName, "connected_at": ""})
	}

	params := http_url.Values{}
	params.Set("page_size", "1")
	req1, err := requestEMQBackend(url, params)
	if err != nil {
		fmt.Printf("设备列表：请求EMQ后端失败：%s\n", err)
	} else {
		r, _ := req1["result"].(map[string]interface{})
		device_total, _ := r["total_num"].(float64)
		if device_total > 0 {
			params.Set("page_size", strconv.FormatFloat(device_total, 'f', 0, 64))
			req2, err := requestEMQBackend(url, params)
			if err != nil {
				fmt.Printf("设备列表：请求EMQ后端失败：%s\n", err)
			} else {
				r, _ := req2["result"].(map[string]interface{})
				infos, _ := r["objects"].([]interface{})
				for _, v := range infos {
					v, _ := v.(map[string]interface{})
					d_id, _ := v["client_id"].(string)
					d_name, _ := v["username"].(string)
					d_connected_at, _ := v["connected_at"].(string)

					for index, db_value := range db_results {
						if db_value["device_id"] == d_id {
							db_results = append(db_results[:index], db_results[index+1:]...) // 删除元素
							db_results = append(db_results, map[string]string{"device_id": d_id, "device_name": d_name, "connected_at": d_connected_at})
						}
					}
				}
				c.JSON(http.StatusOK, db_results)
			}

		}
	}
}

func DeviceOnlineStatus(c *gin.Context) {
	/* 获取某台设备的在线情况 */
	device_id := c.Param("id")
	url := fmt.Sprintf("%s/api/v2/nodes/%s/clients/%s", emq_config["host"], emq_config["node"], device_id)
	req, err := requestEMQBackend(url, http_url.Values{})
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{"success": false, "msg": "请求数据后端错误", "result": nil})
	}
	r, _ := req["result"].(map[string]interface{})
	info, _ := r["objects"].([]interface{})
	if len(info) > 0 {
		c.JSON(http.StatusOK, gin.H{"success": true, "msg": "", "result": true})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": true, "msg": "", "result": false})
	}
}
