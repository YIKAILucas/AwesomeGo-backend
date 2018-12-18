package main

import (
	"awesomeProject/src/mongo"
	"awesomeProject/src/mqttbroker"
	"awesomeProject/src/routers"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

const DB_NAME = "acke"
const DB_COLLECTION = "test"

func MqStart() {
	//mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	mq := &mqttbroker.BrokerInfo{}
	mqttbroker.InfoInit(mq)
	mqttbroker.MqConnect(mq, mqttbroker.HandlerFunc)
	mqttbroker.Sub(mq, "tf/Attendance/v1/notify", 1, nil)
	mqttbroker.Sub(mq, "chat", 1, nil)
	mqttbroker.Sub(mq, "tf/Attendance/v1/devices/+/control", 1, mqttbroker.DeviceControlHandler)
	//mqttbroker.Sub(mq, "tf/Attendance/v1/devices/A6AF-C135-5D68/control", 1, mqttbroker.DeviceControlHandler)
	mqttbroker.Sub(mq, "tf/Attendance/v1/devices/+/info", 1, mqttbroker.DeviceInfoHandler)
	//mqttbroker.Sub(mq, "tf/Attendance/v1/devices/A6AF-C135-5D68/info", 1, mqttbroker.DeviceInfoHandler)
	//c.Disconnect(250)
	select {}
}

func MqSaveDeviceControlLoop() {
	/* 将设备控制指令存入数据库 */
	for data := range mqttbroker.DeviceControlChannel {
		topic := data.Topic()
		payload := data.Payload()

		json_data := make(map[string]interface{})

		err := json.Unmarshal(payload, &json_data)
		if err != nil {
			fmt.Println("话题：%s 的消息，JSON解码失败", topic)
		}
		json_data["topic"] = data.Topic()
		json_data["result"] = make(map[string]interface{})
		_, ok := json_data["arg"]

		if !ok {
			json_data["arg"] = nil
		}
		_, ok = json_data["timeout"]
		if !ok {
			json_data["timeout"] = nil
		}
		json_data["create_time"] = time.Now()
		json_data["last_update_time"] = time.Now()
		mongo.Insert(DB_NAME, DB_COLLECTION, json_data)
		fmt.Printf("控制指令已存入，ID: %d，话题：%s\n", json_data["id"], topic)
	}
}

func MqSaveDeviceInfoLoop() {
	/* 将设备信息存入数据库 */
	for data := range mqttbroker.DeviceInfoChannel {
		topic := data.Topic()
		payload := data.Payload()

		json_data := make(map[string]interface{})
		err := json.Unmarshal(payload, &json_data)
		if err != nil {
			fmt.Printf("话题：%s 的消息，JSON解码失败\n", topic)
		}

		db_data := make(map[string]interface{})
		mongo.FindOne(DB_NAME, DB_COLLECTION, bson.M{"id": json_data["id"]}, nil, db_data)

		if len(db_data) != 0 {
			result := make(map[string]interface{})
			result["success"] = json_data["success"]
			result["message"] = json_data["message"]
			result["device_id"] = json_data["device_id"]
			result["result"] = json_data["result"]
			db_data["result"] = result
			db_data["last_update_time"] = time.Now()
			mongo.Update(DB_NAME, DB_COLLECTION, bson.M{"id": json_data["id"]}, db_data)
			fmt.Printf("设备信息已存入，ID：%d，话题：%s\n", json_data["id"], topic)
		} else {
			fmt.Printf("设备信息未存入，ID：%d，未找到对应的表项\n", json_data["id"])
		}
	}
}

func main() {
	// 初始化MQTT
	go MqStart()
	go MqSaveDeviceControlLoop()
	go MqSaveDeviceInfoLoop()

	// Create the Gin engine.
	gin.SetMode(gin.DebugMode)
	g := gin.New()
	middlewares := []gin.HandlerFunc{}

	// Routes.
	routers.Load(
		// Cores.
		g,

		// Middlwares.
		middlewares...,
	)

	//Ping the server to make sure the router is working.
	go func() {
		if err := pingServer(); err != nil {
			//log.Fatal("The router has no response, or it might took too long to start up.", err)
		}
		log.Print("The router has been deployed successfully.")
	}()

	log.Printf("Start to listening the incoming requests on http address: %s", ":8080")
	log.Printf(http.ListenAndServe(":8080", g).Error())
}

// pingServer pings the http server to make sure the router is working.
func pingServer() error {
	for i := 0; i < 2; i++ {
		// Ping the server by sending a GET request to `/health`.
		resp, err := http.Get("http://127.0.0.1:8080" + "/wechat/push")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		// Sleep for a second to continue the next ping.
		log.Print("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	return errors.New("Cannot connect to the router.")
}
