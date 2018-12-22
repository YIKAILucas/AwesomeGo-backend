package mqttbroker

import (
	"awesomeProject/src/middleware/mongo"
	"encoding/json"
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gopkg.in/mgo.v2/bson"
)

const DB_NAME = "acke"
const DB_COLLECTION = "test"

var ChannelString chan string = make(chan string, 5)
var DeviceControlChannel = make(chan mqtt.Message, 2000)
var DeviceInfoChannel = make(chan mqtt.Message, 2000)

type MQCallback interface {
	handlerFunc(client mqtt.Client, msg mqtt.Message)
}

type Content struct {
	Content string `json:"content"`
}

type BrokerInfo struct {
	client    mqtt.Client
	brokerURL string
	clientId  string
	userName  string
}

func InfoInit(mq *BrokerInfo) {
	//mq.brokerURL = "tcp://106.12.130.179:1883"
	mq.brokerURL = "tcp://127.0.0.1:1883"
	//mq.clientId = string(rand.Int())
	mq.clientId = "E470-B8A3-1"
	mq.userName = "golang-server"
}

func MqConnect(mq *BrokerInfo, handler mqtt.MessageHandler) bool {
	// 连接broker
	opts := mqtt.NewClientOptions().AddBroker(mq.brokerURL).SetClientID(mq.clientId)
	opts.SetUsername(mq.userName)
	opts.SetKeepAlive(30 * time.Second)
	opts.SetAutoReconnect(true)
	// 设置回调
	opts.SetDefaultPublishHandler(handler)
	//opts.SetPingTimeout(1 * time.Second)

	//create client
	mq.client = mqtt.NewClient(opts)
	if token := mq.client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return mq.client.IsConnected()
}

func MqSaveDeviceControlLoop() {
	/* 将设备控制指令存入数据库 */
	for data := range DeviceControlChannel {
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
	for data := range DeviceInfoChannel {
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

func Sub(mq *BrokerInfo, topic string, qos byte, callback mqtt.MessageHandler) {
	if token := mq.client.Subscribe(topic, qos, callback); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}

func Pub(mq *BrokerInfo, topic string, qos byte, payload interface{}, retain bool) error {
	token := mq.client.Publish(topic, qos, retain, payload)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func UnSub(mq *BrokerInfo, topic string) {
	if token := mq.client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}
