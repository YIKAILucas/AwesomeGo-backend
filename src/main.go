package main

import (
	"awesomeProject/src/config"
	"awesomeProject/src/handler"
	"awesomeProject/src/model"
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
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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

	for data := range handler.HTTPPubChannel {
		err := mqttbroker.Pub(mq, data.Topic, 1, data.Payload, false)
		if err != nil {
			fmt.Printf("消息发布失败：%s,Payload：%s", data.Topic, data.Payload)
		}
	}
}

func MqSaveDeviceControlLoop() {
	/* 将设备控制指令存入数据库 */
	for data := range mqttbroker.DeviceControlChannel {
		topic := data.Topic()
		payload := data.Payload()

		json_data := make(map[string]interface{})

		err := json.Unmarshal(payload, &json_data)
		if err != nil {
			fmt.Printf("话题：%s 的消息，JSON解码失败", topic)
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
		if fmt.Sprintf("%T", json_data["id"]) == "string" {
			fmt.Printf("控制指令已存入，ID: %s，话题：%s\n", json_data["id"], topic)
		} else {
			fmt.Printf("控制指令已存入，ID: %.f，话题：%s\n", json_data["id"], topic)
		}
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
			fmt.Println(json_data)

			mongo.Update(DB_NAME, DB_COLLECTION, bson.M{"id": json_data["id"]}, db_data)
			if fmt.Sprintf("%T", json_data["id"]) == "string" {
				fmt.Printf("设备信息已存入，ID：%s，话题：%s\n", json_data["id"], topic)
			} else {
				fmt.Printf("设备信息已存入，ID：%.f，话题：%s\n", json_data["id"], topic)
			}
		} else {
			if fmt.Sprintf("%T", json_data["id"]) == "string" {
				fmt.Printf("设备信息未存入，ID：%s，未找到对应的表项\n", json_data["id"])
			} else {
				fmt.Printf("设备信息未存入，ID：%.f，未找到对应的表项\n", json_data["id"])
			}
		}
	}
}

var (
	cfg = pflag.StringP("config", "c", "", "apiserver config file path.")
)

func main() {
	pflag.Parse()
	// init config
	if err := config.Init(*cfg); err != nil {
		panic(err)
	}
	model.DBInit()

	// 初始化MQTT
	//go MqStart()
	//go MqSaveDeviceControlLoop()
	//go MqSaveDeviceInfoLoop()

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
	log.Printf(http.ListenAndServe(viper.GetString("addr"), g).Error())
}

// pingServer pings the http server to make sure the router is working.
func pingServer() error {
	for i := 0; i < viper.GetInt("max_ping_count"); i++ {
		// Ping the server by sending a GET request to `/health`.
		resp, err := http.Get(viper.GetString("url") + "/sd/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		// Sleep for a second to continue the next ping.
		log.Print("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	return errors.New("Cannot connect to the router.")
}
